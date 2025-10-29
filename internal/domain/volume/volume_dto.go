package volume

import "simple-go/internal/domain/chapter"

type CreateVolumeDTO struct {
	OriginalLanguage string `json:"original_language" binding:"required"`
	Number           int    `json:"number" binding:"required"`
	NovelID          string `json:"novel_id" binding:"required"`
	IsVirtual        bool   `json:"is_virtual"`
	//this go to translation table
	Title       string  `json:"title" binding:"required"`
	Description *string `json:"description"`
}

type CreateVolumeTranslationDTO struct {
	VolumeID    string  `json:"volume_id" binding:"required"`
	Lang        string  `json:"lang" binding:"required"`
	Title       string  `json:"title" binding:"required"`
	Description *string `json:"description"`
}

type VolumeResponseDTO struct {
	ID               string  `json:"id"`
	NovelID          string  `json:"novel_id,omitempty"`
	OriginalLanguage string  `json:"original_language"`
	Number           int     `json:"number"`
	CoverURL         *string `json:"cover_url"`
	Lang             string  `json:"lang"`
	Title            string  `json:"title"`
	Description      *string `json:"description"`
	IsVirtual        bool    `json:"is_virtual"`

	Chapters []chapter.ChapterResponseDTO `json:"chapters,omitempty"`
}
