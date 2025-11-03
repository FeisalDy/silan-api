package repository

import (
	"context"
	"simple-go/internal/domain/job"
)

type TranslationJobRepository interface {
	Create(ctx context.Context, j *job.TranslationJob) (*job.TranslationJob, error)
	GetByID(ctx context.Context, id string) (*job.TranslationJob, error)
	GetByNovelAndLang(ctx context.Context, novelID, targetLang string) (*job.TranslationJob, error)
	GetAll(ctx context.Context, limit, offset int, status string) ([]job.TranslationJob, error)
	GetByNovelID(ctx context.Context, novelID string, limit, offset int) ([]job.TranslationJob, error)
	Update(ctx context.Context, j *job.TranslationJob) (*job.TranslationJob, error)
	UpdateStatus(ctx context.Context, id, status string) error
	UpdateProgress(ctx context.Context, id string, progress, completedSubtasks int) error
	Count(ctx context.Context) (int64, error)

	// Subtask operations
	CreateSubtask(ctx context.Context, subtask *job.TranslationSubtask) (*job.TranslationSubtask, error)
	CreateSubtasksBatch(ctx context.Context, subtasks []job.TranslationSubtask) error
	GetSubtaskByID(ctx context.Context, id string) (*job.TranslationSubtask, error)
	GetSubtasksByJobID(ctx context.Context, jobID string) ([]job.TranslationSubtask, error)
	UpdateSubtask(ctx context.Context, subtask *job.TranslationSubtask) (*job.TranslationSubtask, error)
	UpdateSubtaskStatus(ctx context.Context, id, status string) error
	CountSubtasksByStatus(ctx context.Context, jobID, status string) (int64, error)
}
