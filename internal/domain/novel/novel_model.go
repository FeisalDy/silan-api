package novel

import (
	"time"

	"simple-go/internal/domain/media"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Novel struct {
	ID               string       `gorm:"type:uuid;primaryKey"`
	CreatedBy        string       `gorm:"type:uuid;not null;index"`
	OriginalLanguage string       `gorm:"type:varchar(10);not null"`
	OriginalAuthor   *string      `gorm:"type:varchar(255)"`
	Source           *string      `gorm:"type:varchar(500)"`
	Status           *string      `gorm:"type:varchar(50)"`
	WordCount        *int         `gorm:"type:int"`
	CreatedAt        time.Time    `gorm:"autoCreateTime"`
	UpdatedAt        time.Time    `gorm:"autoUpdateTime"`
	CoverMediaID     *string      `gorm:"type:uuid;index"`
	Media            *media.Media `gorm:"foreignKey:CoverMediaID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	Translations []NovelTranslation `gorm:"foreignKey:NovelID;constraint:OnDelete:CASCADE;"`
}

func (n *Novel) BeforeCreate(tx *gorm.DB) error {
	if n.ID == "" {
		n.ID = uuid.New().String()
	}
	return nil
}

func (Novel) TableName() string {
	return "novels"
}
