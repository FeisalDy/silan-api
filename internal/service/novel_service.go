package service

import (
	"context"
	"errors"
	dommedia "simple-go/internal/domain/media"
	"simple-go/internal/domain/novel"
	"simple-go/internal/domain/volume"
	"simple-go/internal/repository"
	"simple-go/pkg/logger"

	"gorm.io/gorm"
)

type NovelService struct {
	uow        repository.UnitOfWork
	novelRepo  repository.NovelRepository
	mediaSrvc  *MediaService
	volumeSrvc *VolumeService
}

func NewNovelService(uow repository.UnitOfWork, novelRepo repository.NovelRepository, mediaSrvc *MediaService, volumeSrvc *VolumeService) *NovelService {
	return &NovelService{
		uow:        uow,
		novelRepo:  novelRepo,
		mediaSrvc:  mediaSrvc,
		volumeSrvc: volumeSrvc,
	}
}

func (s *NovelService) CreateNovelWithTranslation(
	ctx context.Context,
	creatorID string,
	dto novel.CreateNovelDTO,
) (*novel.Novel, *novel.NovelTranslation, error) {
	var newNovel *novel.Novel
	var newTranslation *novel.NovelTranslation

	err := s.uow.Do(ctx, func(provider repository.RepositoryProvider) error {
		newNovel = &novel.Novel{
			CreatedBy:        creatorID,
			OriginalLanguage: dto.OriginalLanguage,
			OriginalAuthor:   dto.OriginalAuthor,
			Source:           dto.Source,
			Status:           dto.Status,
		}

		if err := provider.Novel().Create(ctx, newNovel); err != nil {
			logger.Error(err, "failed to create novel in database")
			return errors.New("unable to create novel")
		}

		newTranslation = &novel.NovelTranslation{
			NovelID:      newNovel.ID,
			Lang:         dto.OriginalLanguage,
			Title:        dto.Title,
			Description:  dto.Description,
			TranslatorID: creatorID,
		}

		if err := provider.Novel().CreateTranslation(ctx, newTranslation); err != nil {
			logger.Error(err, "failed to create novel translation")
			return errors.New("unable to create novel translation")
		}

		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	return newNovel, newTranslation, nil
}

func (s *NovelService) GetByID(ctx context.Context, id, lang string) (*novel.NovelResponseDTO, error) {
	n, err := s.novelRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("novel not found")
		}
		logger.Error(err, "failed to get novel by ID")
		return nil, errors.New("unable to retrieve novel")
	}

	var selected *novel.NovelTranslation

	if lang != "" {
		for i := range n.Translations {
			if n.Translations[i].Lang == lang {
				selected = &n.Translations[i]
				break
			}
		}
	}

	if selected == nil {
		for i := range n.Translations {
			if n.Translations[i].Lang == n.OriginalLanguage {
				selected = &n.Translations[i]
				break
			}
		}
	}

	if selected == nil && len(n.Translations) > 0 {
		selected = &n.Translations[0]
	}

	if selected == nil {
		return nil, errors.New("no translations available for this novel")
	}

	var coverURL *string
	if n.Media != nil && n.Media.URL != nil {
		coverURL = n.Media.URL
	}

	return &novel.NovelResponseDTO{
		ID:               n.ID,
		OriginalLanguage: n.OriginalLanguage,
		OriginalAuthor:   n.OriginalAuthor,
		Source:           n.Source,
		Status:           n.Status,
		WordCount:        n.WordCount,
		CoverURL:         coverURL,
		Lang:             selected.Lang,
		Title:            selected.Title,
		Description:      selected.Description,
		CreatedAt:        n.CreatedAt,
		UpdatedAt:        n.UpdatedAt,
	}, nil
}

func (s *NovelService) GetAll(ctx context.Context, limit, offset int, title, lang string) ([]novel.NovelResponseDTO, int64, error) {
	novels, err := s.novelRepo.GetAll(ctx, limit, offset, title, lang)
	if err != nil {
		logger.Error(err, "failed to get all novels")
		return nil, 0, errors.New("unable to retrieve novels")
	}

	count, err := s.novelRepo.Count(ctx)
	if err != nil {
		logger.Error(err, "failed to count novels")
		return nil, 0, errors.New("unable to retrieve novels")
	}

	var responses []novel.NovelResponseDTO

	for _, n := range novels {
		var selected *novel.NovelTranslation

		if lang != "" {
			for i := range n.Translations {
				if n.Translations[i].Lang == lang {
					selected = &n.Translations[i]
					break
				}
			}
		}

		if selected == nil {
			for i := range n.Translations {
				if n.Translations[i].Lang == n.OriginalLanguage {
					selected = &n.Translations[i]
					break
				}
			}
		}

		if selected == nil && len(n.Translations) > 0 {
			selected = &n.Translations[0]
		}

		if selected == nil {
			continue
		}

		var coverURL *string
		if n.Media != nil && n.Media.URL != nil {
			coverURL = n.Media.URL
		}

		responses = append(responses, novel.NovelResponseDTO{
			ID:               n.ID,
			OriginalLanguage: n.OriginalLanguage,
			OriginalAuthor:   n.OriginalAuthor,
			Source:           n.Source,
			Status:           n.Status,
			WordCount:        n.WordCount,
			CoverURL:         coverURL,
			Lang:             selected.Lang,
			Title:            selected.Title,
			Description:      selected.Description,
			CreatedAt:        n.CreatedAt,
			UpdatedAt:        n.UpdatedAt,
		})
	}

	return responses, count, nil
}

func (s *NovelService) Delete(ctx context.Context, id string) error {
	if err := s.novelRepo.Delete(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("novel not found")
		}
		logger.Error(err, "failed to delete novel")
		return errors.New("unable to delete novel")
	}
	return nil
}

func (s *NovelService) UpdateCoverMedia(ctx context.Context, id string, dto novel.UpdateCoverMediaDTO) error {
	uploadParams := dommedia.UploadAndSaveDTO{
		Name:       dto.FileName,
		FileBytes:  dto.FileBytes,
		UploaderID: dto.UploaderID,
	}

	return s.uow.Do(ctx, func(provider repository.RepositoryProvider) error {
		savedMedia, _, err := s.mediaSrvc.UploadAndSaveWithRepo(ctx, provider.Media(), uploadParams)
		if err != nil {
			logger.Error(err, "failed to upload and save media for novel cover")
			return errors.New("unable to upload cover media")
		}

		if err := provider.Novel().UpdateCoverMedia(ctx, id, savedMedia.ID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("novel not found")
			}
			logger.Error(err, "failed to update novel cover media")
			return errors.New("unable to update novel cover")
		}

		return nil
	})
}

func (s *NovelService) CreateTranslation(ctx context.Context, translatorID string, dto novel.CreateNovelTranslationDTO) (*novel.NovelTranslation, error) {
	_, err := s.novelRepo.GetByID(ctx, dto.NovelID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("novel not found")
		}
		logger.Error(err, "failed to verify novel exists")
		return nil, errors.New("unable to create translation")
	}

	existing, err := s.novelRepo.GetTranslation(ctx, dto.NovelID, dto.Lang)
	if err == nil && existing != nil {
		return nil, errors.New("translation for this language already exists")
	}

	nt := &novel.NovelTranslation{
		NovelID:      dto.NovelID,
		Lang:         dto.Lang,
		Title:        dto.Title,
		Description:  dto.Description,
		TranslatorID: translatorID,
	}

	if err := s.novelRepo.CreateTranslation(ctx, nt); err != nil {
		logger.Error(err, "failed to create novel translation")
		return nil, errors.New("unable to create translation")
	}

	return nt, nil
}

func (s *NovelService) DeleteTranslation(ctx context.Context, id string) error {
	if err := s.novelRepo.DeleteTranslation(ctx, id); err != nil {
		logger.Error(err, "failed to delete novel translation")
		return errors.New("unable to delete translation")
	}
	return nil
}

func (s *NovelService) GetNovelVolumes(ctx context.Context, novelID, lang string) ([]volume.VolumeResponseDTO, error) {
	return s.volumeSrvc.GetNovelVolumes(ctx, novelID, lang)
}
