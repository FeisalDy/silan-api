package tag

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Tag represents a tag that can be applied to novels
type Tag struct {
	ID          string    `gorm:"type:uuid;primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar(100);unique;not null" json:"name"`
	Slug        string    `gorm:"type:varchar(100);unique;not null;index" json:"slug"`
	Description *string   `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// BeforeCreate will set a UUID
func (t *Tag) BeforeCreate(tx *gorm.DB) error {
	if t.ID == "" {
		t.ID = uuid.New().String()
	}
	return nil
}

// TableName specifies the table name
func (Tag) TableName() string {
	return "tags"
}

// NovelTag represents the many-to-many relationship between novels and tags
type NovelTag struct {
	ID        string    `gorm:"type:uuid;primaryKey" json:"id"`
	NovelID   string    `gorm:"type:uuid;not null;index" json:"novel_id"`
	TagID     string    `gorm:"type:uuid;not null;index" json:"tag_id"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// BeforeCreate will set a UUID
func (nt *NovelTag) BeforeCreate(tx *gorm.DB) error {
	if nt.ID == "" {
		nt.ID = uuid.New().String()
	}
	return nil
}

// TableName specifies the table name
func (NovelTag) TableName() string {
	return "novel_tags"
}

// CreateTagDTO is used for tag creation
type CreateTagDTO struct {
	Name        string  `json:"name" binding:"required"`
	Slug        string  `json:"slug" binding:"required"`
	Description *string `json:"description"`
}

// UpdateTagDTO is used for tag updates
type UpdateTagDTO struct {
	Name        *string `json:"name"`
	Slug        *string `json:"slug"`
	Description *string `json:"description"`
}
