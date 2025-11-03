# Redis Queue Integration for Translation Jobs

This document describes the Redis queue publisher integration for the translation job system.

## Overview

When a translation job is created via the API, it is:
1. Saved to the PostgreSQL database with all subtasks
2. Published to a Redis queue for async processing
3. Picked up by Python workers for actual translation

## Architecture

```
┌─────────────┐      ┌──────────┐      ┌───────────┐      ┌──────────────┐
│   Client    │─────▶│ Go API   │─────▶│  Redis    │─────▶│ Python Worker│
│  (HTTP API) │      │ Service  │      │  Queue    │      │  (Translator)│
└─────────────┘      └──────────┘      └───────────┘      └──────────────┘
                           │                                       │
                           │                                       │
                           ▼                                       ▼
                     ┌──────────┐                          ┌──────────┐
                     │PostgreSQL│◀─────────────────────────│PostgreSQL│
                     │   DB     │  Updates job status      │   DB     │
                     └──────────┘  & translations          └──────────┘
```

## Components

### 1. Redis Queue Publisher (Go)

**File**: `pkg/queue/redis_queue.go`

**Key Functions**:
- `NewRedisQueue(redisURL, queueName)` - Creates connection to Redis
- `Publish(ctx, TranslationJobMessage)` - Publishes a single job
- `PublishBatch(ctx, []TranslationJobMessage)` - Publishes multiple jobs
- `GetQueueLength(ctx)` - Returns current queue size
- `HealthCheck(ctx)` - Checks Redis connection

**Message Structure**:
```go
type TranslationJobMessage struct {
    JobID            string   `json:"job_id"`
    TargetLang       string   `json:"target_lang,omitempty"`
    SourceLang       string   `json:"source_lang,omitempty"`
    OutputDir        string   `json:"output_dir,omitempty"`
    TargetFields     []string `json:"target_fields,omitempty"`
    EnableCodeFilter bool     `json:"enable_code_filter,omitempty"`
}
```

### 2. Service Integration

**File**: `internal/service/translation_job_service.go`

The `TranslationJobService` now accepts an optional `*queue.RedisQueue` parameter. After successfully creating a job in the database, it publishes the job to Redis:

```go
queueMsg := queue.TranslationJobMessage{
    JobID:        createdJob.ID,
    TargetLang:   createdJob.TargetLang,
    SourceLang:   createdJob.FromLang,
    OutputDir:    "outputs/novel_jobs",
    TargetFields: []string{"title", "description", "content"},
}

if err := s.redisQueue.Publish(ctx, queueMsg); err != nil {
    logger.Error(err, "failed to push job to Redis queue")
    // Job is still created, worker can poll from DB if needed
}
```

### 3. Configuration

**File**: `pkg/config/config.go`

Added `RedisConfig` to application configuration:

```go
type RedisConfig struct {
    URL       string  // Redis connection URL
    QueueName string  // Queue name for translation jobs
}
```

**Environment Variables**:
- `REDIS_URL` - Redis connection URL (default: `redis://localhost:6379/0`)
- `TRANSLATION_QUEUE` - Queue name (default: `translation_jobs`)

### 4. Application Initialization

**File**: `internal/app/app.go`

Redis queue is initialized during application startup. If Redis is unavailable, the service continues to work (jobs are created in DB but not queued):

```go
redisQueue, err = queue.NewRedisQueue(cfg.Redis.URL, cfg.Redis.QueueName)
if err != nil {
    logger.Warn("Failed to connect to Redis, translation jobs will not be queued")
    redisQueue = nil
}
```

## Python Worker

The Python worker (provided in your code) continuously polls the Redis queue:

```python
# Main worker loop
while True:
    # Blocking pop from Redis queue (30 second timeout)
    result = redis_client.blpop(queue_name, timeout=30)
    
    if result is None:
        continue  # No jobs available
    
    # Parse job data
    _, job_json = result
    job_data = json.loads(job_json)
    
    # Process the job
    process_job(job_data)
```

## Setup Instructions

### 1. Install Redis

**Using Docker**:
```bash
docker run -d --name redis -p 6379:6379 redis:7-alpine
```

**Using Docker Compose** (add to your `docker-compose.yml`):
```yaml
services:
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    restart: unless-stopped

volumes:
  redis_data:
```

### 2. Configure Go Service

Add to your `.env` file:
```bash
REDIS_URL=redis://localhost:6379/0
TRANSLATION_QUEUE=translation_jobs
```

### 3. Start Python Worker

```bash
# Install dependencies
pip install redis psycopg2-binary

# Set environment variables
export DATABASE_DSN=postgresql://user:pass@localhost:5432/dbname
export REDIS_URL=redis://localhost:6379/0
export TRANSLATION_QUEUE=translation_jobs

# Run worker
python redis_queue_worker.py
```

### 4. Test the Integration

**Create a translation job**:
```bash
curl -X POST http://localhost:8000/api/v1/translation-jobs \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "novel_id": "some-uuid",
    "target_lang": "en"
  }'
```

**Monitor the worker logs** - you should see:
```
Processing translation job: <job-id>
✅ Job <job-id> completed: X entities translated
```

## Monitoring

### Check Queue Length

You can monitor the queue using Redis CLI:
```bash
redis-cli LLEN translation_jobs
```

Or programmatically:
```go
length, err := redisQueue.GetQueueLength(ctx)
```

### Health Check

Check Redis connection health:
```go
if err := redisQueue.HealthCheck(ctx); err != nil {
    // Redis is down
}
```

## Failure Handling

### Redis Unavailable at Startup
- Service starts successfully
- Jobs are created in database
- Queue publishing is skipped
- Workers can poll database directly if needed

### Redis Fails During Operation
- Job creation succeeds
- Queue publish fails (logged as error)
- Job remains in database with `PENDING` status
- Workers can pick up from database

### Worker Fails During Processing
- Job status remains `IN_PROGRESS`
- Implement timeout recovery (reset stuck jobs back to `PENDING`)
- Subtasks track individual progress for retry

## Scaling

### Multiple Workers
The Redis queue pattern naturally supports multiple workers:
- Each worker uses `BLPOP` (blocking pop) with `SKIP LOCKED` semantics
- Only one worker receives each job
- Workers can run on different machines

### Rate Limiting
Implement rate limiting on the Go API side to prevent queue overflow:
```go
// Check queue length before accepting job
if length, _ := redisQueue.GetQueueLength(ctx); length > 1000 {
    return errors.New("queue is full, please try again later")
}
```

## Future Enhancements

1. **Priority Queues**: Use different queue names for different priorities
2. **Dead Letter Queue**: Move failed jobs to a separate queue
3. **Job Retry Logic**: Automatic retry with exponential backoff
4. **Progress Updates**: Workers publish progress updates back to Redis
5. **WebSocket Notifications**: Push real-time updates to clients

## Troubleshooting

### Jobs Created But Not Processed
- Check Redis is running: `redis-cli ping`
- Check worker is running: `ps aux | grep redis_queue_worker`
- Check queue has jobs: `redis-cli LLEN translation_jobs`

### Worker Can't Connect to Redis
- Verify `REDIS_URL` environment variable
- Check network connectivity
- Verify Redis is listening on correct port

### Jobs Stuck in `IN_PROGRESS`
- Worker crashed during processing
- Implement watchdog to reset old jobs:
```sql
UPDATE translation_jobs 
SET status = 'PENDING' 
WHERE status = 'IN_PROGRESS' 
AND updated_at < NOW() - INTERVAL '1 hour';
```

## Dependencies

### Go Dependencies
```bash
go get github.com/redis/go-redis/v9
```

### Python Dependencies
```bash
pip install redis psycopg2-binary
```

## API Examples

### Create Job and Monitor Queue

```go
// Create job
job, err := jobService.CreateTranslationJob(ctx, userID, dto)

// Check if it was queued
length, err := redisQueue.GetQueueLength(ctx)
fmt.Printf("Queue length: %d\n", length)
```

### Batch Job Creation

```go
messages := []queue.TranslationJobMessage{
    {JobID: "job-1", TargetLang: "en"},
    {JobID: "job-2", TargetLang: "fr"},
    {JobID: "job-3", TargetLang: "es"},
}

err := redisQueue.PublishBatch(ctx, messages)
```

## Summary

The Redis queue integration provides:
- ✅ **Asynchronous processing** - API responds immediately
- ✅ **Decoupled architecture** - Go and Python services independent
- ✅ **Scalability** - Multiple workers can process jobs in parallel
- ✅ **Reliability** - Jobs persisted in database, queue is optional
- ✅ **Monitoring** - Track queue length and worker status
- ✅ **Graceful degradation** - System works even if Redis is down
