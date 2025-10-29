package chapter

func SelectTranslation(translations []ChapterTranslation, lang string) *ChapterTranslation {
	if len(translations) == 0 {
		return nil
	}

	if lang != "" {
		for i := range translations {
			if translations[i].Lang == lang {
				return &translations[i]
			}
		}
	}

	for i := range translations {
		if translations[i].Lang == translations[0].Lang {
			return &translations[i]
		}
	}

	return &translations[0]
}
