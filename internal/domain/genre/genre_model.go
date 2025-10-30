package genre

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Genre struct {
	ID          string    `gorm:"type:uuid;primaryKey"`
	Name        string    `gorm:"type:varchar(100);unique;not null"`
	Slug        string    `gorm:"type:varchar(100);unique;not null;index"`
	Description *string   `gorm:"type:text"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

func (g *Genre) BeforeCreate(tx *gorm.DB) error {
	if g.ID == "" {
		g.ID = uuid.New().String()
	}
	return nil
}

func (Genre) TableName() string {
	return "genres"
}
