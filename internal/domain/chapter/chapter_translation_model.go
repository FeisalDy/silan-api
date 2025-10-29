package chapter

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

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

func (ct *ChapterTranslation) BeforeCreate(tx *gorm.DB) error {
	if ct.ID == "" {
		ct.ID = uuid.New().String()
	}
	return nil
}

func (ChapterTranslation) TableName() string {
	return "chapter_translations"
}
