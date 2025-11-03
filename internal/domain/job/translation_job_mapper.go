package job

import "github.com/google/uuid"

func MapTranslationJobToDTO(job TranslationJob) TranslationJobResponseDTO {
	// Defensive zero-value check
	if job.ID == "" || job.ID == uuid.Nil.String() {
		return TranslationJobResponseDTO{}
	}

	return TranslationJobResponseDTO{
		ID:                job.ID,
		NovelID:           job.NovelID,
		FromLang:          job.FromLang,
		TargetLang:        job.TargetLang,
		Status:            job.Status,
		Progress:          job.Progress,
		TotalSubtasks:     job.TotalSubtasks,
		CompletedSubtasks: job.CompletedSubtasks,
		ErrorMessage:      job.ErrorMessage,
		StartedAt:         job.StartedAt,
		FinishedAt:        job.FinishedAt,
		CreatedAt:         job.CreatedAt,
		UpdatedAt:         job.UpdatedAt,
	}
}

func MapTranslationSubtaskToDTO(subtask TranslationSubtask) TranslationSubtaskResponseDTO {
	// Defensive zero-value check
	if subtask.ID == "" || subtask.ID == uuid.Nil.String() {
		return TranslationSubtaskResponseDTO{}
	}

	return TranslationSubtaskResponseDTO{
		ID:             subtask.ID,
		JobID:          subtask.JobID,
		EntityType:     subtask.EntityType,
		EntityID:       subtask.EntityID,
		ParentVolumeID: subtask.ParentVolumeID,
		Seq:            subtask.Seq,
		Priority:       subtask.Priority,
		Status:         subtask.Status,
		ErrorMessage:   subtask.ErrorMessage,
		StartedAt:      subtask.StartedAt,
		FinishedAt:     subtask.FinishedAt,
		CreatedAt:      subtask.CreatedAt,
		UpdatedAt:      subtask.UpdatedAt,
	}
}

func MapTranslationJobToDetailDTO(job TranslationJob) TranslationJobDetailResponseDTO {
	// Handle nil or empty subtasks gracefully
	var subtasks []TranslationSubtaskResponseDTO
	if len(job.Subtasks) > 0 {
		subtasks = make([]TranslationSubtaskResponseDTO, len(job.Subtasks))
		for i, subtask := range job.Subtasks {
			subtasks[i] = MapTranslationSubtaskToDTO(subtask)
		}
	} else {
		subtasks = []TranslationSubtaskResponseDTO{} // return empty slice, not nil
	}

	return TranslationJobDetailResponseDTO{
		TranslationJobResponseDTO: MapTranslationJobToDTO(job),
		Subtasks:                  subtasks,
	}
}
