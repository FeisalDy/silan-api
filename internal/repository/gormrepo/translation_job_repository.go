package gormrepo

import (
	"context"
	"simple-go/internal/domain/job"
	"simple-go/internal/repository"

	"gorm.io/gorm"
)

type translationJobRepository struct {
	db *gorm.DB
}

func NewTranslationJobRepository(db *gorm.DB) repository.TranslationJobRepository {
	return &translationJobRepository{db: db}
}

func (r *translationJobRepository) Create(ctx context.Context, j *job.TranslationJob) (*job.TranslationJob, error) {
	if err := r.db.WithContext(ctx).Create(j).Error; err != nil {
		return nil, err
	}
	return j, nil
}

func (r *translationJobRepository) GetByID(ctx context.Context, id string) (*job.TranslationJob, error) {
	var j job.TranslationJob
	err := r.db.WithContext(ctx).
		Preload("Subtasks", func(db *gorm.DB) *gorm.DB {
			return db.Order("translation_subtasks.priority ASC, translation_subtasks.seq ASC")
		}).
		First(&j, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &j, nil
}

func (r *translationJobRepository) GetByNovelAndLang(ctx context.Context, novelID, targetLang string) (*job.TranslationJob, error) {
	var j job.TranslationJob
	err := r.db.WithContext(ctx).
		Where("novel_id = ? AND target_lang = ?", novelID, targetLang).
		First(&j).Error
	if err != nil {
		return nil, err
	}
	return &j, nil
}

func (r *translationJobRepository) GetAll(ctx context.Context, limit, offset int, status string) ([]job.TranslationJob, error) {
	var jobs []job.TranslationJob

	query := r.db.WithContext(ctx).
		Order("created_at DESC")

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&jobs).Error; err != nil {
		return nil, err
	}

	return jobs, nil
}

func (r *translationJobRepository) GetByNovelID(ctx context.Context, novelID string, limit, offset int) ([]job.TranslationJob, error) {
	var jobs []job.TranslationJob

	query := r.db.WithContext(ctx).
		Where("novel_id = ?", novelID).
		Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&jobs).Error; err != nil {
		return nil, err
	}

	return jobs, nil
}

func (r *translationJobRepository) Update(ctx context.Context, j *job.TranslationJob) (*job.TranslationJob, error) {
	if err := r.db.WithContext(ctx).Save(j).Error; err != nil {
		return nil, err
	}
	return j, nil
}

func (r *translationJobRepository) UpdateStatus(ctx context.Context, id, status string) error {
	return r.db.WithContext(ctx).
		Model(&job.TranslationJob{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *translationJobRepository) UpdateProgress(ctx context.Context, id string, progress, completedSubtasks int) error {
	return r.db.WithContext(ctx).
		Model(&job.TranslationJob{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"progress":           progress,
			"completed_subtasks": completedSubtasks,
		}).Error
}

func (r *translationJobRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&job.TranslationJob{}).Count(&count).Error
	return count, err
}

// Subtask operations

func (r *translationJobRepository) CreateSubtask(ctx context.Context, subtask *job.TranslationSubtask) (*job.TranslationSubtask, error) {
	if err := r.db.WithContext(ctx).Create(subtask).Error; err != nil {
		return nil, err
	}
	return subtask, nil
}

func (r *translationJobRepository) CreateSubtasksBatch(ctx context.Context, subtasks []job.TranslationSubtask) error {
	return r.db.WithContext(ctx).CreateInBatches(subtasks, 100).Error
}

func (r *translationJobRepository) GetSubtaskByID(ctx context.Context, id string) (*job.TranslationSubtask, error) {
	var subtask job.TranslationSubtask
	err := r.db.WithContext(ctx).First(&subtask, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &subtask, nil
}

func (r *translationJobRepository) GetSubtasksByJobID(ctx context.Context, jobID string) ([]job.TranslationSubtask, error) {
	var subtasks []job.TranslationSubtask
	err := r.db.WithContext(ctx).
		Where("job_id = ?", jobID).
		Order("priority ASC, seq ASC").
		Find(&subtasks).Error
	if err != nil {
		return nil, err
	}
	return subtasks, nil
}

func (r *translationJobRepository) UpdateSubtask(ctx context.Context, subtask *job.TranslationSubtask) (*job.TranslationSubtask, error) {
	if err := r.db.WithContext(ctx).Save(subtask).Error; err != nil {
		return nil, err
	}
	return subtask, nil
}

func (r *translationJobRepository) UpdateSubtaskStatus(ctx context.Context, id, status string) error {
	return r.db.WithContext(ctx).
		Model(&job.TranslationSubtask{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *translationJobRepository) CountSubtasksByStatus(ctx context.Context, jobID, status string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&job.TranslationSubtask{}).
		Where("job_id = ? AND status = ?", jobID, status).
		Count(&count).Error
	return count, err
}
