# Translation Job System - Implementation Summary

## Files Created

### Domain Layer (`internal/domain/job/`)
1. **translation_job_model.go** - Enhanced TranslationJob model with full lifecycle tracking
2. **translation_subtask_model.go** - TranslationSubtask model for granular task tracking
3. **translation_job_dto.go** - Request/Response DTOs for API
4. **translation_job_mapper.go** - Mapper functions to convert models to DTOs

### Repository Layer
1. **internal/repository/translation_job_repository.go** - Repository interface
2. **internal/repository/gormrepo/translation_job_repository.go** - GORM implementation
3. **Updated**: `internal/repository/unit_of_work.go` - Added TranslationJob() method
4. **Updated**: `internal/repository/gormrepo/unit_of_work.go` - Added TranslationJob() provider

### Service Layer
1. **internal/service/translation_job_service.go** - Business logic for translation jobs
   - CreateTranslationJob: Creates job with all subtasks (chapters, volumes, novel)
   - GetJobByID: Retrieves job with subtasks
   - GetAllJobs: List all jobs with pagination and status filter
   - GetJobsByNovelID: List jobs for specific novel
   - CancelJob: Cancel a job

### Handler Layer
1. **internal/handler/translation_job_handler.go** - HTTP handlers for translation jobs
   - POST /translation-jobs - Create new job
   - GET /translation-jobs/:id - Get job detail
   - GET /translation-jobs - List all jobs
   - GET /novels/:novel_id/translation-jobs - List jobs for novel
   - PUT /translation-jobs/:id/cancel - Cancel job

## Database Schema

The models will create these tables when you run migration:

### translation_jobs
```sql
- id (UUID, PK)
- novel_id (UUID, FK to novels, indexed)
- from_lang (VARCHAR(10))
- target_lang (VARCHAR(10))
- status (VARCHAR(20), default 'PENDING')
- progress (INT, default 0)
- total_subtasks (INT, default 0)
- completed_subtasks (INT, default 0)
- staging_prefix (TEXT, nullable)
- error_message (TEXT, nullable)
- created_by (UUID, nullable, indexed)
- started_at (TIMESTAMP, nullable)
- finished_at (TIMESTAMP, nullable)
- created_at (TIMESTAMP)
- updated_at (TIMESTAMP)
- UNIQUE INDEX (novel_id, target_lang)
```

### translation_subtasks
```sql
- id (UUID, PK)
- job_id (UUID, FK to translation_jobs, cascade delete, indexed)
- entity_type (VARCHAR(20)) -- 'chapter', 'volume', 'novel'
- entity_id (UUID)
- parent_volume_id (UUID, nullable)
- seq (INT)
- priority (INT, default 100)
- status (VARCHAR(20), default 'PENDING')
- result_path (TEXT, nullable)
- result_text (TEXT, nullable)
- error_message (TEXT, nullable)
- started_at (TIMESTAMP, nullable)
- finished_at (TIMESTAMP, nullable)
- created_at (TIMESTAMP)
- updated_at (TIMESTAMP)
- UNIQUE INDEX (job_id, entity_type, entity_id)
```

## How to Wire Up in Your Application

### 1. Run Database Migration

Add to your migration or AutoMigrate:

```go
db.AutoMigrate(
    &job.TranslationJob{},
    &job.TranslationSubtask{},
    // ... other models
)
```

### 2. Initialize Repository, Service, and Handler

In your `internal/app/app.go` or dependency injection setup:

```go
// Repositories
translationJobRepo := gormrepo.NewTranslationJobRepository(db)

// Services
translationJobService := service.NewTranslationJobService(
    uow,
    translationJobRepo,
    novelRepo,
    volumeRepo,
    chapterRepo,
)

// Handlers
translationJobHandler := handler.NewTranslationJobHandler(translationJobService)
```

### 3. Register Routes

In your router setup (e.g., `internal/server/gin/routes.go`):

```go
// Translation Jobs routes
translationJobRoutes := api.Group("/translation-jobs")
{
    translationJobRoutes.POST("", authMiddleware, translationJobHandler.CreateTranslationJob)
    translationJobRoutes.GET("", authMiddleware, translationJobHandler.GetAllJobs)
    translationJobRoutes.GET("/:id", authMiddleware, translationJobHandler.GetJobByID)
    translationJobRoutes.PUT("/:id/cancel", authMiddleware, translationJobHandler.CancelJob)
}

// Novel-specific translation jobs
novelRoutes.GET("/:novel_id/translation-jobs", authMiddleware, translationJobHandler.GetJobsByNovelID)
```

## API Endpoints

### Create Translation Job
```http
POST /api/translation-jobs
Content-Type: application/json
Authorization: Bearer <token>

{
    "novel_id": "uuid-here",
    "target_lang": "en"
}
```

Response:
```json
{
    "status": "success",
    "message": "Translation job created successfully",
    "data": {
        "id": "job-uuid",
        "novel_id": "novel-uuid",
        "from_lang": "ja",
        "target_lang": "en",
        "status": "PENDING",
        "progress": 0,
        "total_subtasks": 153,
        "completed_subtasks": 0,
        "created_at": "2025-11-03T10:00:00Z",
        "updated_at": "2025-11-03T10:00:00Z"
    }
}
```

### Get Job Detail
```http
GET /api/translation-jobs/:id
```

Response includes job info + all subtasks ordered by priority and sequence.

### List All Jobs
```http
GET /api/translation-jobs?page=1&limit=10&status=PENDING
```

### Get Jobs for Novel
```http
GET /api/novels/:novel_id/translation-jobs?page=1&limit=10
```

### Cancel Job
```http
PUT /api/translation-jobs/:id/cancel
```

## Job Creation Flow

When a user creates a translation job:

1. **Validates** novel exists and gets original language
2. **Checks** for existing active jobs (novel + target_lang is unique)
3. **Fetches** all volumes and chapters for the novel
4. **Creates** parent job with status PENDING
5. **Creates subtasks**:
   - Chapter subtasks (priority 200) - one per chapter
   - Volume subtasks (priority 150) - one per volume  
   - Novel subtask (priority 100) - one for novel metadata
6. **Updates** job with total_subtasks count
7. **Returns** job details

## Priority System

- **Priority 100** (highest): Novel metadata - processed last
- **Priority 150**: Volume metadata - processed after chapters
- **Priority 200**: Chapter content - processed first

Workers should query:
```sql
ORDER BY priority ASC, seq ASC
FOR UPDATE SKIP LOCKED
```

## Status Values

### Job Statuses:
- `PENDING` - Job created, waiting to be picked up
- `IN_PROGRESS` - Worker is processing
- `COMPLETED` - All subtasks done, translations promoted
- `FAILED` - Job failed or was cancelled

### Subtask Statuses:
- `PENDING` - Not started
- `IN_PROGRESS` - Worker processing
- `DONE` - Successfully translated
- `FAILED` - Translation failed

## Next Steps for Python Microservice

1. **Worker subscribes to Redis queue** for new job IDs
2. **Worker queries subtasks** for the job (ORDER BY priority, seq)
3. **For each subtask**:
   - Fetch entity content (chapter/volume/novel)
   - Translate using your ML model
   - Store result in staging (S3 or result_path/result_text)
   - Update subtask status to DONE
   - Update job progress
4. **When all subtasks DONE**:
   - Run promotion transaction
   - Upsert all translations atomically
   - Mark job COMPLETED

## Files Summary

- **4 domain files** (models + DTOs + mapper)
- **2 repository files** (interface + implementation)
- **1 service file** (business logic)
- **1 handler file** (HTTP endpoints)
- **2 updated files** (unit of work)

Total: 10 files created/updated
