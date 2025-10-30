package novel

import (
	"time"

	"simple-go/internal/domain/genre"
	"simple-go/internal/domain/media"
	"simple-go/internal/domain/tag"
	"simple-go/internal/domain/volume"

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
	Volumes      []volume.Volume    `gorm:"foreignKey:NovelID;constraint:OnDelete:CASCADE;"`

	Tags   []tag.Tag     `gorm:"many2many:novel_tags;joinForeignKey:NovelID;joinReferences:TagID"`
	Genres []genre.Genre `gorm:"many2many:novel_genres;joinForeignKey:NovelID;joinReferences:GenreID"`
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
