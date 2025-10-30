package tag

// MapTagToDTO converts a Tag model to TagResponseDTO
func MapTagToDTO(t Tag) TagResponseDTO {
	return TagResponseDTO{
		ID:          t.ID,
		Name:        t.Name,
		Slug:        t.Slug,
		Description: t.Description,
	}
}

// MapTagsToResponseDTOs converts a slice of Tag models to TagResponseDTO
func MapTagsToResponseDTOs(tags []Tag) []TagResponseDTO {
	if tags == nil {
		return []TagResponseDTO{}
	}

	dtos := make([]TagResponseDTO, len(tags))
	for i, t := range tags {
		dtos[i] = MapTagToDTO(t)
	}
	return dtos
}

// MapTagToUpdateDTO converts a Tag model to UpdateTagDTO
func MapTagToUpdateDTO(t Tag) UpdateTagDTO {
	return UpdateTagDTO{
		Name:        &t.Name,
		Slug:        &t.Slug,
		Description: t.Description,
	}
}

// MapTagsToUpdateDTOs converts a slice of Tag models to UpdateTagDTO
func MapTagsToUpdateDTOs(tags []Tag) []UpdateTagDTO {
	if tags == nil {
		return []UpdateTagDTO{}
	}

	dtos := make([]UpdateTagDTO, len(tags))
	for i, t := range tags {
		dtos[i] = MapTagToUpdateDTO(t)
	}
	return dtos
}
