package novel

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Novel struct {
	ID               string    `gorm:"type:uuid;primaryKey"`
	CreatedBy        string    `gorm:"type:uuid;not null;index"`
	OriginalLanguage string    `gorm:"type:varchar(10);not null"`
	OriginalAuthor   *string   `gorm:"type:varchar(255)"`
	Source           *string   `gorm:"type:varchar(500)"`
	Status           *string   `gorm:"type:varchar(50)"`
	WordCount        *int      `gorm:"type:int"`
	CoverMediaID     *string   `gorm:"type:uuid;index"`
	CreatedAt        time.Time `gorm:"autoCreateTime"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime"`

	Translations []NovelTranslation `gorm:"foreignKey:NovelID;constraint:OnDelete:CASCADE;"`
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
	ID           string    `gorm:"type:uuid;primaryKey"`
	NovelID      string    `gorm:"type:uuid;not null;index:idx_novel_lang,unique"`
	Lang         string    `gorm:"type:varchar(10);not null;index:idx_novel_lang,unique"`
	Title        string    `gorm:"type:varchar(500);not null"`
	Description  *string   `gorm:"type:text"`
	Summary      *string   `gorm:"type:text"`
	TranslatorID string    `gorm:"type:uuid;not null;index"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
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
	ID               string                     `json:"id"`
	CreatedBy        string                     `json:"created_by"`
	OriginalLanguage string                     `json:"original_language"`
	OriginalAuthor   *string                    `json:"original_author"`
	Source           *string                    `json:"source"`
	Status           *string                    `json:"status"`
	WordCount        *int                       `json:"word_count"`
	CoverMediaID     *string                    `json:"cover_media_id"`
	Translations     []NovelTranslationResponse `json:"translations,omitempty"`
	CreatedAt        time.Time                  `json:"created_at"`
	UpdatedAt        time.Time                  `json:"updated_at"`
}

type NovelTranslationResponse struct {
	ID           string    `json:"id"`
	NovelID      string    `json:"novel_id"`
	Lang         string    `json:"lang"`
	Title        string    `json:"title"`
	Description  *string   `json:"description"`
	Summary      *string   `json:"summary"`
	TranslatorID string    `json:"translator_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
