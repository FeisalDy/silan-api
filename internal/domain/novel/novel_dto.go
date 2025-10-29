package novel

import "time"

type CreateNovelDTO struct {
	OriginalLanguage string  `json:"original_language" binding:"required"`
	OriginalAuthor   *string `json:"original_author"`
	Source           *string `json:"source"`
	Status           *string `json:"status"`
	//this go to translation table
	Title       string  `json:"title" binding:"required"`
	Description *string `json:"description"`
}

type CreateNovelTranslationDTO struct {
	NovelID     string  `json:"novel_id" binding:"required"`
	Lang        string  `json:"lang" binding:"required"`
	Title       string  `json:"title" binding:"required"`
	Description *string `json:"description"`
}

type NovelResponseDTO struct {
	ID               string    `json:"id"`
	OriginalLanguage string    `json:"original_language"`
	OriginalAuthor   *string   `json:"original_author"`
	Source           *string   `json:"source"`
	Status           *string   `json:"status"`
	WordCount        *int      `json:"word_count"`
	CoverURL         *string   `json:"cover_url"`
	Lang             string    `json:"lang"`
	Title            string    `json:"title"`
	Description      *string   `json:"description"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type NovelTranslationResponseDTO struct {
	ID          string  `json:"id"`
	NovelID     string  `json:"novel_id"`
	Lang        string  `json:"lang"`
	Title       string  `json:"title"`
	Description *string `json:"description"`
	// TranslatorID string    `json:"translator_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UpdateCoverMediaDTO struct {
	FileName   string `json:"-"`
	FileBytes  []byte `json:"-"`
	UploaderID string `json:"-"`
}
