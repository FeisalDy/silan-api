package volume

func SelectTranslation(translations []VolumeTranslation, lang, originalLang string) *VolumeTranslation {
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
		if translations[i].Lang == originalLang {
			return &translations[i]
		}
	}

	return &translations[0]
}
