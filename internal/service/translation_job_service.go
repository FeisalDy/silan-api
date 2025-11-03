package service

import (
	"context"
	"errors"
	"fmt"

	"simple-go/internal/domain/chapter"
	"simple-go/internal/domain/job"
	"simple-go/internal/domain/volume"
	"simple-go/internal/repository"
	"simple-go/pkg/logger"
	"simple-go/pkg/queue"

	"gorm.io/gorm"
)

type TranslationJobService struct {
	uow         repository.UnitOfWork
	jobRepo     repository.TranslationJobRepository
	novelRepo   repository.NovelRepository
	volumeRepo  repository.VolumeRepository
	chapterRepo repository.ChapterRepository
	redisQueue  *queue.RedisQueue
}

func NewTranslationJobService(
	uow repository.UnitOfWork,
	jobRepo repository.TranslationJobRepository,
	novelRepo repository.NovelRepository,
	volumeRepo repository.VolumeRepository,
	chapterRepo repository.ChapterRepository,
	redisQueue *queue.RedisQueue,
) *TranslationJobService {
	return &TranslationJobService{
		uow:         uow,
		jobRepo:     jobRepo,
		novelRepo:   novelRepo,
		volumeRepo:  volumeRepo,
		chapterRepo: chapterRepo,
		redisQueue:  redisQueue,
	}
}

// CreateTranslationJob creates a new translation job for a novel with all subtasks
func (s *TranslationJobService) CreateTranslationJob(
	ctx context.Context,
	userID string,
	dto job.CreateTranslationJobDTO,
) (*job.TranslationJobResponseDTO, error) {
	var createdJob *job.TranslationJob

	err := s.uow.Do(ctx, func(provider repository.RepositoryProvider) error {
		novel, err := provider.Novel().GetByID(ctx, dto.NovelID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("novel not found")
			}
			logger.Error(err, "failed to get novel")
			return errors.New("unable to verify novel")
		}

		existingJob, err := provider.TranslationJob().GetByNovelAndLang(ctx, dto.NovelID, dto.TargetLang)
		if err == nil && existingJob != nil {
			// Check if it's an active job
			if existingJob.Status == job.TranslationJobStatusPending || existingJob.Status == job.TranslationJobStatusInProgress {
				return fmt.Errorf("active translation job already exists for this novel and language (job_id: %s)", existingJob.ID)
			}
		}

		// 3. Get all volumes with chapters for the novel
		volumes, err := provider.Volume().GetAllWithChaptersByNovelID(ctx, dto.NovelID)
		if err != nil {
			logger.Error(err, "failed to get volumes")
			return errors.New("unable to get volumes for novel")
		}

		// 4. Collect all chapters from volumes
		chapters := make([]chapter.Chapter, 0)
		for _, vol := range volumes {
			chapters = append(chapters, vol.Chapters...)
		}

		if len(chapters) == 0 {
			return errors.New("novel has no chapters to translate")
		}

		// 5. Create the parent job
		newJob := &job.TranslationJob{
			NovelID:    dto.NovelID,
			FromLang:   novel.OriginalLanguage,
			TargetLang: dto.TargetLang,
			Status:     job.TranslationJobStatusPending,
			CreatedBy:  &userID,
		}

		createdJob, err = provider.TranslationJob().Create(ctx, newJob)
		if err != nil {
			logger.Error(err, "failed to create translation job")
			return errors.New("unable to create translation job")
		}

		// 6. Create subtasks for chapters, volumes, and novel
		subtasks := make([]job.TranslationSubtask, 0)

		// Chapter subtasks (priority 200)
		for i, ch := range chapters {
			subtasks = append(subtasks, job.TranslationSubtask{
				JobID:          createdJob.ID,
				EntityType:     job.EntityTypeChapter,
				EntityID:       ch.ID,
				ParentVolumeID: &ch.VolumeID,
				Seq:            i + 1,
				Priority:       job.PriorityChapter,
				Status:         job.TranslationSubtaskStatusPending,
			})
		}

		// Volume subtasks (priority 150)
		for i, vol := range volumes {
			subtasks = append(subtasks, job.TranslationSubtask{
				JobID:      createdJob.ID,
				EntityType: job.EntityTypeVolume,
				EntityID:   vol.ID,
				Seq:        i + 1,
				Priority:   job.PriorityVolume,
				Status:     job.TranslationSubtaskStatusPending,
			})
		}

		// Novel subtask (priority 100 - highest, processed last)
		subtasks = append(subtasks, job.TranslationSubtask{
			JobID:      createdJob.ID,
			EntityType: job.EntityTypeNovel,
			EntityID:   dto.NovelID,
			Seq:        1,
			Priority:   job.PriorityNovel,
			Status:     job.TranslationSubtaskStatusPending,
		})

		// 7. Batch insert subtasks
		if err := provider.TranslationJob().CreateSubtasksBatch(ctx, subtasks); err != nil {
			logger.Error(err, "failed to create subtasks")
			return errors.New("unable to create translation subtasks")
		}

		// 8. Update job with total subtasks count
		createdJob.TotalSubtasks = len(subtasks)
		if _, err := provider.TranslationJob().Update(ctx, createdJob); err != nil {
			logger.Error(err, "failed to update job with total subtasks")
			return errors.New("unable to update job")
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 9. Push job to Redis queue (after successful DB transaction)
	if s.redisQueue != nil {
		queueMsg := queue.TranslationJobMessage{
			JobID:        createdJob.ID,
			TargetLang:   createdJob.TargetLang,
			SourceLang:   createdJob.FromLang,
			TargetFields: []string{"title", "description", "content"},
		}

		if err := s.redisQueue.Publish(ctx, queueMsg); err != nil {
			logger.Error(err, "failed to push job to Redis queue")
			// Don't fail the request if queue push fails - job is already created
			// Worker can poll from DB if needed
		} else {
			logger.Info(fmt.Sprintf("Successfully pushed job %s to Redis queue", createdJob.ID))
		}
	}

	response := job.MapTranslationJobToDTO(*createdJob)
	return &response, nil
}

// GetJobByID retrieves a translation job by ID with all subtasks
func (s *TranslationJobService) GetJobByID(ctx context.Context, id string) (*job.TranslationJobDetailResponseDTO, error) {
	j, err := s.jobRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("translation job not found")
		}
		logger.Error(err, "failed to get translation job")
		return nil, errors.New("unable to retrieve translation job")
	}

	response := job.MapTranslationJobToDetailDTO(*j)
	return &response, nil
}

func (s *TranslationJobService) GetAllJobs(ctx context.Context, limit, offset int, status string) ([]job.TranslationJobResponseDTO, int64, error) {
	jobs, err := s.jobRepo.GetAll(ctx, limit, offset, status)
	if err != nil {
		logger.Error(err, "failed to get translation jobs")
		return nil, 0, errors.New("unable to retrieve translation jobs")
	}

	count, err := s.jobRepo.Count(ctx)
	if err != nil {
		logger.Error(err, "failed to count translation jobs")
		return nil, 0, errors.New("unable to count translation jobs")
	}

	// Return empty array if no jobs found
	if len(jobs) == 0 {
		return []job.TranslationJobResponseDTO{}, 0, nil
	}

	response := make([]job.TranslationJobResponseDTO, len(jobs))
	for i, j := range jobs {
		response[i] = job.MapTranslationJobToDTO(j)
	}

	return response, count, nil
}

// GetJobsByNovelID retrieves all translation jobs for a specific novel
func (s *TranslationJobService) GetJobsByNovelID(ctx context.Context, novelID string, limit, offset int) ([]job.TranslationJobResponseDTO, error) {
	jobs, err := s.jobRepo.GetByNovelID(ctx, novelID, limit, offset)
	if err != nil {
		logger.Error(err, "failed to get translation jobs for novel")
		return nil, errors.New("unable to retrieve translation jobs")
	}

	response := make([]job.TranslationJobResponseDTO, len(jobs))
	for i, j := range jobs {
		response[i] = job.MapTranslationJobToDTO(j)
	}

	return response, nil
}

// CancelJob cancels a pending or in-progress translation job
func (s *TranslationJobService) CancelJob(ctx context.Context, id string) error {
	j, err := s.jobRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("translation job not found")
		}
		logger.Error(err, "failed to get translation job")
		return errors.New("unable to retrieve translation job")
	}

	if j.Status == job.TranslationJobStatusCompleted {
		return errors.New("cannot cancel a completed job")
	}

	if j.Status == job.TranslationJobStatusFailed {
		return errors.New("job already failed")
	}

	if err := s.jobRepo.UpdateStatus(ctx, id, job.TranslationJobStatusFailed); err != nil {
		logger.Error(err, "failed to cancel translation job")
		return errors.New("unable to cancel translation job")
	}

	return nil
}

// Helper method to get chapters by novel ID (if not already available in repository)
func (s *TranslationJobService) getChaptersByNovelID(ctx context.Context, novelID string) ([]chapter.Chapter, error) {
	volumes, err := s.volumeRepo.GetAllWithChaptersByNovelID(ctx, novelID)
	if err != nil {
		return nil, err
	}

	var allChapters []chapter.Chapter
	for _, vol := range volumes {
		allChapters = append(allChapters, vol.Chapters...)
	}

	return allChapters, nil
}

// Helper method to get volumes by novel ID
func (s *TranslationJobService) getVolumesByNovelID(ctx context.Context, novelID string) ([]volume.Volume, error) {
	return s.volumeRepo.GetAllWithChaptersByNovelID(ctx, novelID)
}
