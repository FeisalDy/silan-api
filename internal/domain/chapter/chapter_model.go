package chapter

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Chapter struct {
	ID        string    `gorm:"type:uuid;primaryKey"`
	VolumeID  string    `gorm:"type:uuid;not null;index"`
	Number    int       `gorm:"type:int;not null"`
	WordCount *int      `gorm:"type:int"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`

	Translations []ChapterTranslation `gorm:"foreignKey:ChapterID;constraint:OnDelete:CASCADE"`
}

func (c *Chapter) BeforeCreate(tx *gorm.DB) error {
	if c.ID == "" {
		c.ID = uuid.New().String()
	}
	return nil
}

func (Chapter) TableName() string {
	return "chapters"
}
