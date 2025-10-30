package genre

// MapGenreToDTO converts a Genre model to GenreResponseDTO
func MapGenreToDTO(g Genre) GenreResponseDTO {
	return GenreResponseDTO{
		ID:          g.ID,
		Name:        g.Name,
		Slug:        g.Slug,
		Description: g.Description,
	}
}

func MapGenresToResponseDTOs(genres []Genre) []GenreResponseDTO {
	if genres == nil {
		return []GenreResponseDTO{}
	}

	dtos := make([]GenreResponseDTO, len(genres))
	for i, g := range genres {
		dtos[i] = MapGenreToDTO(g)
	}
	return dtos
}

func MapGenreToUpdateDTO(g Genre) UpdateGenreDTO {
	return UpdateGenreDTO{
		Name:        &g.Name,
		Slug:        &g.Slug,
		Description: g.Description,
	}
}

func MapGenresToUpdateDTOs(genres []Genre) []UpdateGenreDTO {
	if genres == nil {
		return []UpdateGenreDTO{}
	}

	dtos := make([]UpdateGenreDTO, len(genres))
	for i, g := range genres {
		dtos[i] = MapGenreToUpdateDTO(g)
	}
	return dtos
}
