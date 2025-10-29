package volume

import "simple-go/internal/domain/chapter"

func MapVolumeToDTO(v Volume, lang string) VolumeResponseDTO {
	selected := SelectTranslation(v.Translations, lang, v.OriginalLanguage)

	var coverURL *string
	if v.Media != nil && v.Media.URL != nil {
		coverURL = v.Media.URL
	}

	return VolumeResponseDTO{
		ID:               v.ID,
		NovelID:          v.NovelID,
		OriginalLanguage: v.OriginalLanguage,
		Number:           v.Number,
		CoverURL:         coverURL,
		Lang:             selected.Lang,
		Title:            selected.Title,
		Description:      selected.Description,
		IsVirtual:        v.IsVirtual,
		Chapters:         mapChaptersToDTO(v.Chapters, lang),
	}
}

func mapChaptersToDTO(chapters []chapter.Chapter, lang string) []chapter.ChapterResponseDTO {
	chapterDTOs := make([]chapter.ChapterResponseDTO, 0, len(chapters))
	for _, ch := range chapters {
		selected := chapter.SelectTranslation(ch.Translations, lang)
		title := ""
		if selected != nil {
			title = selected.Title
		}

		chapterDTOs = append(chapterDTOs, chapter.ChapterResponseDTO{
			ID:     ch.ID,
			Number: ch.Number,
			Title:  title,
		})
	}
	return chapterDTOs
}
