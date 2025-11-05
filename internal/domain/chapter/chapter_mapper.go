package chapter

func MapChapterToDTO(c Chapter, lang string, nextID, prevID *string) ChapterResponseDTO {
	selected := SelectTranslation(c.Translations, lang)

	return ChapterResponseDTO{
		ID:                c.ID,
		VolumeID:          c.VolumeID,
		Number:            c.Number,
		WordCount:         c.WordCount,
		Title:             selected.Title,
		Content:           selected.Content,
		NextChapterID:     nextID,
		PreviousChapterID: prevID,
		CreatedAt:         c.CreatedAt,
		UpdatedAt:         c.UpdatedAt,
	}
}

func MapChapterAndTranslationToDTO(c Chapter, t ChapterTranslation) ChapterResponseDTO {
	return ChapterResponseDTO{
		ID:                c.ID,
		VolumeID:          c.VolumeID,
		Number:            c.Number,
		WordCount:         c.WordCount,
		Title:             t.Title,
		Content:           t.Content,
		NextChapterID:     nil,
		PreviousChapterID: nil,
		CreatedAt:         c.CreatedAt,
		UpdatedAt:         c.UpdatedAt,
	}
}
