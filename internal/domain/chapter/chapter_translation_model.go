package chapter

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ChapterTranslation struct {
	ID        string    `gorm:"type:uuid;primaryKey"`
	ChapterID string    `gorm:"type:uuid;not null;index"`
	Lang      string    `gorm:"type:varchar(10);not null;index"`
	Title     string    `gorm:"type:varchar(500);not null"`
	Content   string    `gorm:"type:text;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
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
