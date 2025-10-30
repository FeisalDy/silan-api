package novel

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NovelTranslation struct {
	ID          string    `gorm:"type:uuid;primaryKey"`
	NovelID     string    `gorm:"type:uuid;not null;index:idx_novel_lang,unique"`
	Lang        string    `gorm:"type:varchar(10);not null;index:idx_novel_lang,unique"`
	Title       string    `gorm:"type:varchar(500);not null"`
	Description *string   `gorm:"type:text"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

func (nt *NovelTranslation) BeforeCreate(tx *gorm.DB) error {
	if nt.ID == "" {
		nt.ID = uuid.New().String()
	}
	return nil
}

func (NovelTranslation) TableName() string {
	return "novel_translations"
}
