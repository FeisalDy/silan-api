package genre

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Genre represents a genre category for novels
type Genre struct {
	ID          string    `gorm:"type:uuid;primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(100);unique;not null" json:"name"`
	Slug        string    `gorm:"type:varchar(100);unique;not null;index" json:"slug"`
	Description *string   `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// BeforeCreate will set a UUID
func (g *Genre) BeforeCreate(tx *gorm.DB) error {
	if g.ID == "" {
		g.ID = uuid.New().String()
	}
	return nil
}

// TableName specifies the table name
func (Genre) TableName() string {
	return "genres"
}

// NovelGenre represents the many-to-many relationship between novels and genres
type NovelGenre struct {
	ID        string    `gorm:"type:uuid;primaryKey" json:"id"`
	NovelID   string    `gorm:"type:uuid;not null;index" json:"novel_id"`
	GenreID   string    `gorm:"type:uuid;not null;index" json:"genre_id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// BeforeCreate will set a UUID
func (ng *NovelGenre) BeforeCreate(tx *gorm.DB) error {
	if ng.ID == "" {
		ng.ID = uuid.New().String()
	}
	return nil
}

// TableName specifies the table name
func (NovelGenre) TableName() string {
	return "novel_genres"
}

// CreateGenreDTO is used for genre creation
type CreateGenreDTO struct {
	Name        string  `json:"name" binding:"required"`
	Slug        string  `json:"slug" binding:"required"`
	Description *string `json:"description"`
}

// UpdateGenreDTO is used for genre updates
type UpdateGenreDTO struct {
	Name        *string `json:"name"`
	Slug        *string `json:"slug"`
	Description *string `json:"description"`
}
