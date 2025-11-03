package job

import "time"

type CreateTranslationJobDTO struct {
	NovelID    string `json:"novel_id" binding:"required"`
	TargetLang string `json:"target_lang" binding:"required"`
}

type TranslationJobResponseDTO struct {
	ID                string     `json:"id"`
	NovelID           string     `json:"novel_id"`
	FromLang          string     `json:"from_lang"`
	TargetLang        string     `json:"target_lang"`
	Status            string     `json:"status"`
	Progress          int        `json:"progress"`
	TotalSubtasks     int        `json:"total_subtasks"`
	CompletedSubtasks int        `json:"completed_subtasks"`
	ErrorMessage      *string    `json:"error_message,omitempty"`
	StartedAt         *time.Time `json:"started_at,omitempty"`
	FinishedAt        *time.Time `json:"finished_at,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

type TranslationSubtaskResponseDTO struct {
	ID             string     `json:"id"`
	JobID          string     `json:"job_id"`
	EntityType     string     `json:"entity_type"`
	EntityID       string     `json:"entity_id"`
	ParentVolumeID *string    `json:"parent_volume_id,omitempty"`
	Seq            int        `json:"seq"`
	Priority       int        `json:"priority"`
	Status         string     `json:"status"`
	ErrorMessage   *string    `json:"error_message,omitempty"`
	StartedAt      *time.Time `json:"started_at,omitempty"`
	FinishedAt     *time.Time `json:"finished_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

type TranslationJobDetailResponseDTO struct {
	TranslationJobResponseDTO
	Subtasks []TranslationSubtaskResponseDTO `json:"subtasks"`
}
