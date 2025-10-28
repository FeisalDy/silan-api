package media

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Media struct {
	ID          string    `gorm:"type:uuid;primaryKey"`
	URL         *string   `gorm:"type:text"`
	Type        *string   `gorm:"type:varchar(50)"`
	UploadedBy  *string   `gorm:"type:uuid"`
	UploadedAt  time.Time `gorm:"autoCreateTime"`
	Description *string   `gorm:"type:text"`
	FileSize    *int64    `gorm:"type:int"`
	MimeType    *string   `gorm:"type:varchar(100)"`
}

func (m *Media) BeforeCreate(tx *gorm.DB) error {
	if m.ID == "" {
		m.ID = uuid.New().String()
	}
	return nil
}

func (Media) TableName() string {
	return "medias"
}
