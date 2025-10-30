package role

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role struct {
	ID          string    `gorm:"type:uuid;primaryKey"`
	Name        string    `gorm:"type:varchar(50);unique;not null"`
	Description *string   `gorm:"type:text"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

func (r *Role) BeforeCreate(tx *gorm.DB) error {
	if r.ID == "" {
		r.ID = uuid.New().String()
	}
	return nil
}

func (Role) TableName() string {
	return "roles"
}

const (
	RoleAdmin      = "admin"
	RoleTranslator = "translator"
	RoleUser       = "user"
)
