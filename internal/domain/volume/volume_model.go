package volume

import (
	"simple-go/internal/domain/chapter"
	"simple-go/internal/domain/media"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Volume struct {
	ID               string       `gorm:"type:uuid;primaryKey"`
	Number           int          `gorm:"type:int;not null"`
	OriginalLanguage string       `gorm:"type:varchar(10);not null"`
	CoverMediaID     *string      `gorm:"type:uuid;index"`
	Media            *media.Media `gorm:"foreignKey:CoverMediaID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	NovelID          string       `gorm:"type:uuid;not null;index"`
	IsVirtual        bool         `gorm:"type:boolean;not null"`

	Translations []VolumeTranslation `gorm:"foreignKey:VolumeID;constraint:OnDelete:CASCADE;"`
	Chapters     []chapter.Chapter   `gorm:"foreignKey:VolumeID;constraint:OnDelete:CASCADE;"`
}

func (v *Volume) BeforeCreate(tx *gorm.DB) error {
	if v.ID == "" {
		v.ID = uuid.New().String()
	}
	return nil
}

func (Volume) TableName() string {
	return "volumes"
}
