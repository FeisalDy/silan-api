package novel

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Novel represents a novel in the system
type Novel struct {
	ID               string    `gorm:"type:uuid;primaryKey" json:"id"`
	CreatedBy        string    `gorm:"type:uuid;not null;index" json:"created_by"`
	OriginalLanguage string    `gorm:"type:varchar(10);not null" json:"original_language"`
	OriginalAuthor   *string   `gorm:"type:varchar(255)" json:"original_author"`
	Source           *string   `gorm:"type:varchar(500)" json:"source"`
	Status           *string   `gorm:"type:varchar(50)" json:"status"`
	WordCount        *int      `gorm:"type:int" json:"word_count"`
	CoverMediaID     *string   `gorm:"type:uuid;index" json:"cover_media_id"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// BeforeCreate will set a UUID
func (n *Novel) BeforeCreate(tx *gorm.DB) error {
	if n.ID == "" {
		n.ID = uuid.New().String()
	}
	return nil
}

// TableName specifies the table name
func (Novel) TableName() string {
	return "novels"
}

// NovelTranslation represents a translation of a novel
type NovelTranslation struct {
	ID           string    `gorm:"type:uuid;primaryKey" json:"id"`
	NovelID      string    `gorm:"type:uuid;not null;index" json:"novel_id"`
	Lang         string    `gorm:"type:varchar(10);not null;index" json:"lang"`
	Title        string    `gorm:"type:varchar(500);not null" json:"title"`
	Description  *string   `gorm:"type:text" json:"description"`
	Summary      *string   `gorm:"type:text" json:"summary"`
	TranslatorID string    `gorm:"type:uuid;not null;index" json:"translator_id"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// BeforeCreate will set a UUID
func (nt *NovelTranslation) BeforeCreate(tx *gorm.DB) error {
	if nt.ID == "" {
		nt.ID = uuid.New().String()
	}
	return nil
}

// TableName specifies the table name
func (NovelTranslation) TableName() string {
	return "novel_translations"
}

// CreateNovelDTO is used for novel creation
type CreateNovelDTO struct {
	OriginalLanguage string  `json:"original_language" binding:"required"`
	OriginalAuthor   *string `json:"original_author"`
	Source           *string `json:"source"`
	Status           *string `json:"status"`
	CoverMediaID     *string `json:"cover_media_id"`
}

// UpdateNovelDTO is used for novel updates
type UpdateNovelDTO struct {
	OriginalAuthor *string `json:"original_author"`
	Source         *string `json:"source"`
	Status         *string `json:"status"`
	WordCount      *int    `json:"word_count"`
	CoverMediaID   *string `json:"cover_media_id"`
}

// CreateNovelTranslationDTO is used for creating novel translations
type CreateNovelTranslationDTO struct {
	NovelID     string  `json:"novel_id" binding:"required"`
	Lang        string  `json:"lang" binding:"required"`
	Title       string  `json:"title" binding:"required"`
	Description *string `json:"description"`
	Summary     *string `json:"summary"`
}

// UpdateNovelTranslationDTO is used for updating novel translations
type UpdateNovelTranslationDTO struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Summary     *string `json:"summary"`
}

// NovelWithTranslation combines novel and its translation
type NovelWithTranslation struct {
	Novel       Novel            `json:"novel"`
	Translation NovelTranslation `json:"translation"`
}

// NovelResponse is used for API responses
type NovelResponse struct {
	ID               string             `json:"id"`
	CreatedBy        string             `json:"created_by"`
	OriginalLanguage string             `json:"original_language"`
	OriginalAuthor   *string            `json:"original_author"`
	Source           *string            `json:"source"`
	Status           *string            `json:"status"`
	WordCount        *int               `json:"word_count"`
	CoverMediaID     *string            `json:"cover_media_id"`
	Translations     []NovelTranslation `json:"translations,omitempty"`
	CreatedAt        time.Time          `json:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at"`
}
