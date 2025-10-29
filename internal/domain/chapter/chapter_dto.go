package chapter

import "time"

type CreateChapterDTO struct {
	VolumeID  string `json:"volume_id" binding:"required"`
	Number    int    `json:"number" binding:"required"`
	WordCount *int   `json:"word_count"`
}

type UpdateChapterDTO struct {
	Number    *int `json:"number"`
	WordCount *int `json:"word_count"`
}

type CreateChapterTranslationDTO struct {
	ChapterID string `json:"chapter_id" binding:"required"`
	Lang      string `json:"lang" binding:"required"`
	Title     string `json:"title" binding:"required"`
	Content   string `json:"content" binding:"required"`
}

type UpdateChapterTranslationDTO struct {
	Title   *string `json:"title"`
	Content *string `json:"content"`
}

type ChapterWithTranslation struct {
	Chapter     Chapter            `json:"chapter"`
	Translation ChapterTranslation `json:"translation"`
}

type ChapterResponseDTO struct {
	ID        string     `json:"id"`
	VolumeID  string     `json:"volume_id,omitempty"`
	Number    int        `json:"number"`
	WordCount *int       `json:"word_count,omitempty"`
	Title     string     `json:"title"`
	Content   string     `json:"content,omitempty"`
	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`
}
