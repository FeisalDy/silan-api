package chapter

func MapChapterToDTO(c Chapter, lang string) ChapterResponseDTO {
	selected := SelectTranslation(c.Translations, lang)

	return ChapterResponseDTO{
		ID:        c.ID,
		VolumeID:  c.VolumeID,
		Number:    c.Number,
		WordCount: c.WordCount,
		Title:     selected.Title,
		Content:   selected.Content,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

func MapChapterAndTranslationToDTO(c Chapter, t ChapterTranslation) ChapterResponseDTO {
	return ChapterResponseDTO{
		ID:        c.ID,
		VolumeID:  c.VolumeID,
		Number:    c.Number,
		WordCount: c.WordCount,
		Title:     t.Title,
		Content:   t.Content,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}
