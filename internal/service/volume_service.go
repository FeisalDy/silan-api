package service

import (
	"context"
	"errors"
	"simple-go/internal/domain/volume"
	"simple-go/internal/repository"
	"simple-go/pkg/logger"
)

type VolumeService struct {
	uow        repository.UnitOfWork
	volumeRepo repository.VolumeRepository
	mediaSrvc  *MediaService
}

func NewVolumeService(uow repository.UnitOfWork, volumeRepo repository.VolumeRepository, mediaSrvc *MediaService) *VolumeService {
	return &VolumeService{
		uow:        uow,
		volumeRepo: volumeRepo,
		mediaSrvc:  mediaSrvc,
	}
}

func (s *VolumeService) GetNovelVolumes(ctx context.Context, novelID, lang string) ([]volume.VolumeResponseDTO, error) {
	var volumes []volume.Volume
	var err error

	if lang != "" {
		volumes, err = s.volumeRepo.GetAllWithChaptersByNovelIDAndLang(ctx, novelID, lang)
	} else {
		volumes, err = s.volumeRepo.GetAllWithChaptersByNovelID(ctx, novelID)
	}

	if err != nil {
		logger.Error(err, "failed to get novel volumes")
		return nil, errors.New("unable to get novel volumes")
	}

	response := make([]volume.VolumeResponseDTO, len(volumes))
	for i, v := range volumes {
		response[i] = volume.MapVolumeToDTO(v, lang)
	}

	return response, nil
}
