package service

import (
	"context"
	"errors"
	"simple-go/internal/domain/novel"
	"simple-go/internal/repository"
	"simple-go/pkg/logger"

	"gorm.io/gorm"
)

type NovelService struct {
	uow       repository.UnitOfWork
	novelRepo repository.NovelRepository
}

func NewNovelService(uow repository.UnitOfWork, novelRepo repository.NovelRepository) *NovelService {
	return &NovelService{
		uow:       uow,
		novelRepo: novelRepo,
	}
}

func (s *NovelService) CreateNovelWithTranslation(
	ctx context.Context,
	creatorID string,
	dto novel.CreateNovelDTO,
	translationDTO novel.CreateNovelTranslationDTO,
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
			CoverMediaID:     dto.CoverMediaID,
		}

		if err := provider.Novel().Create(ctx, newNovel); err != nil {
			logger.Error(err, "failed to create novel in database")
			return errors.New("unable to create novel")
		}

		newTranslation = &novel.NovelTranslation{
			NovelID:      newNovel.ID,
			Lang:         translationDTO.Lang,
			Title:        translationDTO.Title,
			Description:  translationDTO.Description,
			Summary:      translationDTO.Summary,
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

func (s *NovelService) GetByIDWithTranslations(ctx context.Context, id, lang string) (*novel.NovelResponse, error) {
	n, err := s.novelRepo.GetByIDWithTranslations(ctx, id, lang)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("novel not found")
		}
		logger.Error(err, "failed to get novel by ID")
		return nil, errors.New("unable to retrieve novel")
	}

	var novelTranslationResponse []novel.NovelTranslationResponse

	for _, t := range n.Translations {
		novelTranslationResponse = append(novelTranslationResponse, novel.NovelTranslationResponse(t))
	}

	return &novel.NovelResponse{
		ID:               n.ID,
		CreatedBy:        n.CreatedBy,
		OriginalLanguage: n.OriginalLanguage,
		OriginalAuthor:   n.OriginalAuthor,
		Source:           n.Source,
		Status:           n.Status,
		WordCount:        n.WordCount,
		CoverMediaID:     n.CoverMediaID,
		Translations:     novelTranslationResponse,
		CreatedAt:        n.CreatedAt,
		UpdatedAt:        n.UpdatedAt,
	}, nil
}

func (s *NovelService) GetAll(ctx context.Context, limit, offset int, title, lang string) ([]novel.NovelResponse, int64, error) {
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

	var novelResponses []novel.NovelResponse

	for _, n := range novels {
		var translations []novel.NovelTranslationResponse
		for _, t := range n.Translations {
			translations = append(translations, novel.NovelTranslationResponse(t))
		}

		novelResponse := novel.NovelResponse{
			ID:               n.ID,
			CreatedBy:        n.CreatedBy,
			OriginalLanguage: n.OriginalLanguage,
			OriginalAuthor:   n.OriginalAuthor,
			Source:           n.Source,
			Status:           n.Status,
			WordCount:        n.WordCount,
			CoverMediaID:     n.CoverMediaID,
			Translations:     translations,
			CreatedAt:        n.CreatedAt,
			UpdatedAt:        n.UpdatedAt,
		}

		novelResponses = append(novelResponses, novelResponse)
	}

	return novelResponses, count, nil
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
		Summary:      dto.Summary,
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
