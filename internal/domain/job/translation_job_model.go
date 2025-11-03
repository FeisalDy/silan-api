package job

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	TranslationJobStatusPending    = "PENDING"
	TranslationJobStatusInProgress = "IN_PROGRESS"
	TranslationJobStatusCompleted  = "COMPLETED"
	TranslationJobStatusFailed     = "FAILED"
)

type TranslationJob struct {
	ID                string               `gorm:"type:uuid;primaryKey"`
	NovelID           string               `gorm:"type:uuid;not null;index;uniqueIndex:idx_translation_job_novel_lang,priority:1"`
	FromLang          string               `gorm:"type:varchar(10);not null"`
	TargetLang        string               `gorm:"type:varchar(10);not null;uniqueIndex:idx_translation_job_novel_lang,priority:2"`
	Status            string               `gorm:"type:varchar(20);not null;default:'PENDING'"`
	Progress          int                  `gorm:"type:int;not null;default:0"`
	TotalSubtasks     int                  `gorm:"type:int;not null;default:0"`
	CompletedSubtasks int                  `gorm:"type:int;not null;default:0"`
	StagingPrefix     *string              `gorm:"type:text"`
	ErrorMessage      *string              `gorm:"type:text"`
	CreatedBy         *string              `gorm:"type:uuid;index"`
	StartedAt         *time.Time           `gorm:"type:timestamp"`
	FinishedAt        *time.Time           `gorm:"type:timestamp"`
	CreatedAt         time.Time            `gorm:"autoCreateTime"`
	UpdatedAt         time.Time            `gorm:"autoUpdateTime"`
	Subtasks          []TranslationSubtask `gorm:"foreignKey:JobID;constraint:OnDelete:CASCADE"`
}

func (j *TranslationJob) BeforeCreate(tx *gorm.DB) error {
	if j.ID == "" {
		j.ID = uuid.New().String()
	}
	return nil
}

func (TranslationJob) TableName() string {
	return "translation_jobs"
}
