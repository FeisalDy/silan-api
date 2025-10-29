package volume

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type VolumeTranslation struct {
	ID           string    `gorm:"type:uuid;primaryKey"`
	VolumeID     string    `gorm:"type:uuid;not null;index:idx_volume_lang,unique"`
	Lang         string    `gorm:"type:varchar(10);not null;index:idx_volume_lang,unique"`
	Title        string    `gorm:"type:varchar(500);not null"`
	Description  *string   `gorm:"type:text"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
	TranslatorID string    `gorm:"type:uuid;not null;index"`
}

func (vt *VolumeTranslation) BeforeCreate(tx *gorm.DB) error {
	if vt.ID == "" {
		vt.ID = uuid.New().String()
	}
	return nil
}

func (VolumeTranslation) TableName() string {
	return "volume_translations"
}
