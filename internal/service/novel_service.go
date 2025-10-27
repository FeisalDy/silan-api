package service

import (
	"context"
	"errors"
	"fmt"
	"simple-go/internal/domain/novel"
	"simple-go/internal/repository"

	"gorm.io/gorm"
)

// NovelService handles novel-related business logic
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

// CreateNovelWithTranslation creates a novel and its initial translation in one transaction
func (s *NovelService) CreateNovelWithTranslation(
	ctx context.Context,
	creatorID string,
	dto novel.CreateNovelDTO,
	translationDTO novel.CreateNovelTranslationDTO,
) (*novel.Novel, *novel.NovelTranslation, error) {
	var newNovel *novel.Novel
	var newTranslation *novel.NovelTranslation

	err := s.uow.Do(ctx, func(provider repository.RepositoryProvider) error {
		// Create novel
		newNovel = &novel.Novel{
			CreatedBy:        creatorID,
			OriginalLanguage: dto.OriginalLanguage,
			OriginalAuthor:   dto.OriginalAuthor,
			Source:           dto.Source,
			Status:           dto.Status,
			CoverMediaID:     dto.CoverMediaID,
		}

		if err := provider.Novel().Create(ctx, newNovel); err != nil {
			return fmt.Errorf("failed to create novel: %w", err)
		}

		// Create translation
		newTranslation = &novel.NovelTranslation{
			NovelID:      newNovel.ID,
			Lang:         translationDTO.Lang,
			Title:        translationDTO.Title,
			Description:  translationDTO.Description,
			Summary:      translationDTO.Summary,
			TranslatorID: creatorID,
		}

		if err := provider.Novel().CreateTranslation(ctx, newTranslation); err != nil {
			return fmt.Errorf("failed to create translation: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, nil, err
	}

	return newNovel, newTranslation, nil
}

// GetByID retrieves a novel by ID
func (s *NovelService) GetByID(ctx context.Context, id string) (*novel.Novel, error) {
	n, err := s.novelRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("novel not found")
		}
		return nil, fmt.Errorf("failed to get novel: %w", err)
	}
	return n, nil
}

// GetByIDWithTranslations retrieves a novel with all its translations
func (s *NovelService) GetByIDWithTranslations(ctx context.Context, id string) (*novel.NovelResponse, error) {
	n, err := s.novelRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("novel not found")
		}
		return nil, fmt.Errorf("failed to get novel: %w", err)
	}

	translations, err := s.novelRepo.GetTranslations(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get translations: %w", err)
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
		Translations:     translations,
		CreatedAt:        n.CreatedAt,
		UpdatedAt:        n.UpdatedAt,
	}, nil
}

// GetAll retrieves all novels with pagination
func (s *NovelService) GetAll(ctx context.Context, limit, offset int) ([]novel.Novel, int64, error) {
	novels, err := s.novelRepo.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get novels: %w", err)
	}

	count, err := s.novelRepo.Count(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count novels: %w", err)
	}

	return novels, count, nil
}

// GetByCreator retrieves novels created by a specific user
func (s *NovelService) GetByCreator(ctx context.Context, creatorID string, limit, offset int) ([]novel.Novel, error) {
	novels, err := s.novelRepo.GetByCreator(ctx, creatorID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get novels: %w", err)
	}
	return novels, nil
}

// Update updates a novel
func (s *NovelService) Update(ctx context.Context, id string, dto novel.UpdateNovelDTO) (*novel.Novel, error) {
	n, err := s.novelRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("novel not found")
		}
		return nil, fmt.Errorf("failed to get novel: %w", err)
	}

	// Update fields
	if dto.OriginalAuthor != nil {
		n.OriginalAuthor = dto.OriginalAuthor
	}
	if dto.Source != nil {
		n.Source = dto.Source
	}
	if dto.Status != nil {
		n.Status = dto.Status
	}
	if dto.WordCount != nil {
		n.WordCount = dto.WordCount
	}
	if dto.CoverMediaID != nil {
		n.CoverMediaID = dto.CoverMediaID
	}

	if err := s.novelRepo.Update(ctx, n); err != nil {
		return nil, fmt.Errorf("failed to update novel: %w", err)
	}

	return n, nil
}

// Delete deletes a novel
func (s *NovelService) Delete(ctx context.Context, id string) error {
	if err := s.novelRepo.Delete(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("novel not found")
		}
		return fmt.Errorf("failed to delete novel: %w", err)
	}
	return nil
}

// Translation operations

// CreateTranslation creates a new translation for a novel
func (s *NovelService) CreateTranslation(ctx context.Context, translatorID string, dto novel.CreateNovelTranslationDTO) (*novel.NovelTranslation, error) {
	// Verify novel exists
	_, err := s.novelRepo.GetByID(ctx, dto.NovelID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("novel not found")
		}
		return nil, fmt.Errorf("failed to verify novel: %w", err)
	}

	// Check if translation already exists
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
		return nil, fmt.Errorf("failed to create translation: %w", err)
	}

	return nt, nil
}

// GetTranslation retrieves a specific translation
func (s *NovelService) GetTranslation(ctx context.Context, novelID, lang string) (*novel.NovelTranslation, error) {
	nt, err := s.novelRepo.GetTranslation(ctx, novelID, lang)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("translation not found")
		}
		return nil, fmt.Errorf("failed to get translation: %w", err)
	}
	return nt, nil
}

// UpdateTranslation updates a translation
func (s *NovelService) UpdateTranslation(ctx context.Context, id string, dto novel.UpdateNovelTranslationDTO) (*novel.NovelTranslation, error) {
	nt, err := s.novelRepo.GetTranslationByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("translation not found")
		}
		return nil, fmt.Errorf("failed to get translation: %w", err)
	}

	// Update fields
	if dto.Title != nil {
		nt.Title = *dto.Title
	}
	if dto.Description != nil {
		nt.Description = dto.Description
	}
	if dto.Summary != nil {
		nt.Summary = dto.Summary
	}

	if err := s.novelRepo.UpdateTranslation(ctx, nt); err != nil {
		return nil, fmt.Errorf("failed to update translation: %w", err)
	}

	return nt, nil
}

// DeleteTranslation deletes a translation
func (s *NovelService) DeleteTranslation(ctx context.Context, id string) error {
	if err := s.novelRepo.DeleteTranslation(ctx, id); err != nil {
		return fmt.Errorf("failed to delete translation: %w", err)
	}
	return nil
}
