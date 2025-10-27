package chapter

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Chapter represents a chapter of a novel
type Chapter struct {
	ID        string    `gorm:"type:uuid;primaryKey" json:"id"`
	NovelID   string    `gorm:"type:uuid;not null;index" json:"novel_id"`
	Number    int       `gorm:"type:int;not null" json:"number"`
	WordCount *int      `gorm:"type:int" json:"word_count"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// BeforeCreate will set a UUID
func (c *Chapter) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return nil
}

// TableName specifies the table name
func (Chapter) TableName() string {
	return "chapters"
}

// ChapterTranslation represents a translation of a chapter
type ChapterTranslation struct {
	ID           string    `gorm:"type:uuid;primaryKey" json:"id"`
	ChapterID    string    `gorm:"type:uuid;not null;index" json:"chapter_id"`
	Lang         string    `gorm:"type:varchar(10);not null;index" json:"lang"`
	Title        string    `gorm:"type:varchar(500);not null" json:"title"`
	Content      string    `gorm:"type:text;not null" json:"content"`
	TranslatorID string    `gorm:"type:uuid;not null;index" json:"translator_id"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// BeforeCreate will set a UUID
func (ct *ChapterTranslation) BeforeCreate(tx *gorm.DB) error {
	if ct.ID == "" {
		ct.ID = uuid.New().String()
	}
	return nil
}

// TableName specifies the table name
func (ChapterTranslation) TableName() string {
	return "chapter_translations"
}

// CreateChapterDTO is used for chapter creation
type CreateChapterDTO struct {
	NovelID   string `json:"novel_id" binding:"required"`
	Number    int    `json:"number" binding:"required"`
	WordCount *int   `json:"word_count"`
}

// UpdateChapterDTO is used for chapter updates
type UpdateChapterDTO struct {
	Number    *int `json:"number"`
	WordCount *int `json:"word_count"`
}

// CreateChapterTranslationDTO is used for creating chapter translations
type CreateChapterTranslationDTO struct {
	ChapterID string `json:"chapter_id" binding:"required"`
	Lang      string `json:"lang" binding:"required"`
	Title     string `json:"title" binding:"required"`
	Content   string `json:"content" binding:"required"`
}

// UpdateChapterTranslationDTO is used for updating chapter translations
type UpdateChapterTranslationDTO struct {
	Title   *string `json:"title"`
	Content *string `json:"content"`
}

// ChapterWithTranslation combines chapter and its translation
type ChapterWithTranslation struct {
	Chapter     Chapter            `json:"chapter"`
	Translation ChapterTranslation `json:"translation"`
}

// ChapterResponse is used for API responses
type ChapterResponse struct {
	ID           string               `json:"id"`
	NovelID      string               `json:"novel_id"`
	Number       int                  `json:"number"`
	WordCount    *int                 `json:"word_count"`
	Translations []ChapterTranslation `json:"translations,omitempty"`
	CreatedAt    time.Time            `json:"created_at"`
	UpdatedAt    time.Time            `json:"updated_at"`
}
