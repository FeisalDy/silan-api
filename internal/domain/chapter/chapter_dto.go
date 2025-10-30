package chapter

import "time"

type CreateChapterDTO struct {
	VolumeID  string `json:"volume_id" binding:"required"`
	Number    int    `json:"number" binding:"required"`
	WordCount *int   `json:"word_count"`
	Title     string `json:"title" binding:"required"`
	Content   string `json:"content" binding:"required"`
	Lang      string `json:"lang" binding:"required"`
}

type CreateChapterTranslationDTO struct {
	ChapterID string `json:"chapter_id" binding:"required"`
	Lang      string `json:"lang" binding:"required"`
	Title     string `json:"title" binding:"required"`
	Content   string `json:"content" binding:"required"`
}

type ChapterResponseDTO struct {
	ID        string    `json:"id"`
	VolumeID  string    `json:"volume_id,omitempty"`
	Number    int       `json:"number"`
	WordCount *int      `json:"word_count"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
