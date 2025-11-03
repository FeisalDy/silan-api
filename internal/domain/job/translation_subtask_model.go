package job

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	TranslationSubtaskStatusPending    = "PENDING"
	TranslationSubtaskStatusInProgress = "IN_PROGRESS"
	TranslationSubtaskStatusDone       = "DONE"
	TranslationSubtaskStatusFailed     = "FAILED"

	// Entity types for subtasks
	EntityTypeChapter = "chapter"
	EntityTypeVolume  = "volume"
	EntityTypeNovel   = "novel"

	// Priority values (lower = higher priority in processing)
	PriorityNovel   = 100 // Processed last (highest priority number)
	PriorityVolume  = 150 // Processed after chapters
	PriorityChapter = 200 // Processed first (lowest priority number)
)

type TranslationSubtask struct {
	ID             string     `gorm:"type:uuid;primaryKey"`
	JobID          string     `gorm:"type:uuid;not null;index;uniqueIndex:idx_translation_subtask_entity,priority:1"`
	EntityType     string     `gorm:"type:varchar(20);not null;uniqueIndex:idx_translation_subtask_entity,priority:2"`
	EntityID       string     `gorm:"type:uuid;not null;uniqueIndex:idx_translation_subtask_entity,priority:3"`
	ParentVolumeID *string    `gorm:"type:uuid"`
	Seq            int        `gorm:"type:int"`
	Priority       int        `gorm:"type:int;not null;default:100"`
	Status         string     `gorm:"type:varchar(20);not null;default:'PENDING'"`
	ResultPath     *string    `gorm:"type:text"`
	ResultText     *string    `gorm:"type:text"`
	ErrorMessage   *string    `gorm:"type:text"`
	StartedAt      *time.Time `gorm:"type:timestamp"`
	FinishedAt     *time.Time `gorm:"type:timestamp"`
	CreatedAt      time.Time  `gorm:"autoCreateTime"`
	UpdatedAt      time.Time  `gorm:"autoUpdateTime"`
}

func (s *TranslationSubtask) BeforeCreate(tx *gorm.DB) error {
	if s.ID == "" {
		s.ID = uuid.New().String()
	}
	return nil
}

func (TranslationSubtask) TableName() string {
	return "translation_subtasks"
}
