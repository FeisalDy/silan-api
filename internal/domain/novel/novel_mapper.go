package novel

import (
	"simple-go/internal/domain/genre"
	"simple-go/internal/domain/tag"
)

func MapNovelToDTO(n Novel, lang string) NovelResponseDTO {
	selected := SelectTranslation(n.Translations, lang, n.OriginalLanguage)

	var coverURL *string
	if n.Media != nil && n.Media.URL != nil {
		coverURL = n.Media.URL
	}

	return NovelResponseDTO{
		ID:               n.ID,
		OriginalLanguage: n.OriginalLanguage,
		OriginalAuthor:   n.OriginalAuthor,
		Source:           n.Source,
		Status:           n.Status,
		WordCount:        n.WordCount,
		CoverURL:         coverURL,
		Lang:             selected.Lang,
		Title:            selected.Title,
		Description:      selected.Description,
		Tags:             tag.MapTagsToUpdateDTOs(n.Tags),
		Genres:           genre.MapGenresToUpdateDTOs(n.Genres),
		CreatedAt:        n.CreatedAt,
		UpdatedAt:        n.UpdatedAt,
	}
}
