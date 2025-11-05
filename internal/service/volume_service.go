package service

import (
	"context"
	"errors"
	"simple-go/internal/domain/chapter"
	"simple-go/internal/domain/volume"
	"simple-go/internal/repository"
	"simple-go/pkg/logger"
)

type VolumeService struct {
	uow         repository.UnitOfWork
	volumeRepo  repository.VolumeRepository
	chapterRepo repository.ChapterRepository
	mediaSrvc   *MediaService
}

func NewVolumeService(uow repository.UnitOfWork, volumeRepo repository.VolumeRepository, chapterRepo repository.ChapterRepository, mediaSrvc *MediaService) *VolumeService {
	return &VolumeService{
		uow:         uow,
		volumeRepo:  volumeRepo,
		chapterRepo: chapterRepo,
		mediaSrvc:   mediaSrvc,
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

// GetChapterWithCrossVolumeNavigation returns a chapter with next/prev IDs including cross-volume navigation
func (s *VolumeService) GetChapterWithCrossVolumeNavigation(ctx context.Context, chapterID, lang string) (*chapter.ChapterResponseDTO, error) {
	// Get the chapter first
	c, err := s.chapterRepo.GetByID(ctx, chapterID)
	if err != nil {
		logger.Error(err, "failed to get chapter by ID")
		return nil, errors.New("unable to retrieve chapter")
	}

	// Get the volume to access novel_id and volume number
	vol, err := s.volumeRepo.GetByID(ctx, c.VolumeID)
	if err != nil {
		logger.Error(err, "failed to get volume for chapter navigation")
		return nil, errors.New("unable to retrieve volume")
	}

	// Fetch next chapter ID (same volume first, then next volume)
	nextID, err := s.chapterRepo.GetNextChapterID(ctx, c.VolumeID, c.Number)
	if err != nil {
		logger.Error(err, "failed to get next chapter ID")
		nextID = nil
	}

	// If no next chapter in current volume, check next volume
	if nextID == nil {
		nextVolumeID, err := s.volumeRepo.GetNextVolumeID(ctx, vol.NovelID, vol.Number)
		if err != nil {
			logger.Error(err, "failed to get next volume ID")
		} else if nextVolumeID != nil {
			// Get first chapter of next volume
			nextID, err = s.chapterRepo.GetFirstChapterIDOfVolume(ctx, *nextVolumeID)
			if err != nil {
				logger.Error(err, "failed to get first chapter of next volume")
			}
		}
	}

	// Fetch previous chapter ID (same volume first, then previous volume)
	prevID, err := s.chapterRepo.GetPreviousChapterID(ctx, c.VolumeID, c.Number)
	if err != nil {
		logger.Error(err, "failed to get previous chapter ID")
		prevID = nil
	}

	// If no previous chapter in current volume, check previous volume
	if prevID == nil {
		prevVolumeID, err := s.volumeRepo.GetPreviousVolumeID(ctx, vol.NovelID, vol.Number)
		if err != nil {
			logger.Error(err, "failed to get previous volume ID")
		} else if prevVolumeID != nil {
			// Get last chapter of previous volume
			prevID, err = s.chapterRepo.GetLastChapterIDOfVolume(ctx, *prevVolumeID)
			if err != nil {
				logger.Error(err, "failed to get last chapter of previous volume")
			}
		}
	}

	// Map to DTO with navigation
	res := chapter.MapChapterToDTO(*c, lang, nextID, prevID)
	return &res, nil
}
