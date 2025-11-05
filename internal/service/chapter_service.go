package service

import (
	"context"
	"errors"
	"simple-go/internal/domain/chapter"
	"simple-go/internal/repository"
	"simple-go/pkg/logger"

	"gorm.io/gorm"
)

// ChapterService handles chapter-related business logic
type ChapterService struct {
	uow         repository.UnitOfWork
	chapterRepo repository.ChapterRepository
}

func NewChapterService(
	uow repository.UnitOfWork,
	chapterRepo repository.ChapterRepository,
) *ChapterService {
	return &ChapterService{
		uow:         uow,
		chapterRepo: chapterRepo,
	}
}

func (s *ChapterService) CreateChapterWithTranslation(
	ctx context.Context,
	chapterDTO chapter.CreateChapterDTO,
) (*chapter.ChapterResponseDTO, error) {
	var newChapter *chapter.Chapter
	var newTranslation *chapter.ChapterTranslation

	err := s.uow.Do(ctx, func(provider repository.RepositoryProvider) error {
		_, err := provider.Novel().GetByID(ctx, chapterDTO.VolumeID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("novel not found")
			}
			logger.Error(err, "failed to verify novel exists")
			return errors.New("unable to create chapter")
		}

		newChapter = &chapter.Chapter{
			VolumeID:  chapterDTO.VolumeID,
			Number:    chapterDTO.Number,
			WordCount: chapterDTO.WordCount,
		}

		createdChapter, err := provider.Chapter().Create(ctx, newChapter)
		if err != nil {
			logger.Error(err, "failed to create chapter")
			return errors.New("unable to create chapter")
		}
		newChapter = createdChapter

		newTranslation = &chapter.ChapterTranslation{
			ChapterID: newChapter.ID,
			Lang:      chapterDTO.Lang,
			Title:     chapterDTO.Title,
			Content:   chapterDTO.Content,
		}

		createdTranslation, err := provider.Chapter().CreateTranslation(ctx, newTranslation)
		if err != nil {
			logger.Error(err, "failed to create chapter translation")
			return errors.New("unable to create chapter translation")
		}
		newTranslation = createdTranslation

		return nil
	})

	if err != nil {
		return nil, err
	}
	res := chapter.MapChapterAndTranslationToDTO(*newChapter, *newTranslation)
	return &res, nil
}

func (s *ChapterService) GetByID(ctx context.Context, id, lang string) (*chapter.ChapterResponseDTO, error) {
	c, err := s.chapterRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("chapter not found")
		}
		logger.Error(err, "failed to get chapter by ID")
		return nil, errors.New("unable to retrieve chapter")
	}

	// Get next and previous chapter IDs within the same volume only
	nextID, err := s.chapterRepo.GetNextChapterID(ctx, c.VolumeID, c.Number)
	if err != nil {
		logger.Error(err, "failed to get next chapter ID")
		nextID = nil
	}

	prevID, err := s.chapterRepo.GetPreviousChapterID(ctx, c.VolumeID, c.Number)
	if err != nil {
		logger.Error(err, "failed to get previous chapter ID")
		prevID = nil
	}

	res := chapter.MapChapterToDTO(*c, lang, nextID, prevID)
	return &res, nil
}

// GetByIDWithNavigation returns a chapter with cross-volume navigation
// This is called by VolumeService to provide full navigation
func (s *ChapterService) GetByIDWithNavigation(ctx context.Context, id, lang string, nextID, prevID *string) (*chapter.ChapterResponseDTO, error) {
	c, err := s.chapterRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("chapter not found")
		}
		logger.Error(err, "failed to get chapter by ID")
		return nil, errors.New("unable to retrieve chapter")
	}

	res := chapter.MapChapterToDTO(*c, lang, nextID, prevID)
	return &res, nil
}

// Delete deletes a chapter
func (s *ChapterService) Delete(ctx context.Context, id string) error {
	affected, err := s.chapterRepo.Delete(ctx, id)
	if err != nil {
		logger.Error(err, "failed to delete chapter")
		return errors.New("unable to delete chapter")
	}
	if affected == 0 {
		return errors.New("chapter not found")
	}
	return nil
}

func (s *ChapterService) CreateTranslation(ctx context.Context, dto chapter.CreateChapterTranslationDTO) (*chapter.ChapterTranslation, error) {
	_, err := s.chapterRepo.GetByID(ctx, dto.ChapterID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("chapter not found")
		}
		logger.Error(err, "failed to verify chapter exists")
		return nil, errors.New("unable to create translation")
	}

	// Check if translation already exists
	existing, err := s.chapterRepo.GetTranslation(ctx, dto.ChapterID, dto.Lang)
	if err == nil && existing != nil {
		return nil, errors.New("translation for this language already exists")
	}

	ct := &chapter.ChapterTranslation{
		ChapterID: dto.ChapterID,
		Lang:      dto.Lang,
		Title:     dto.Title,
		Content:   dto.Content,
	}

	createdTranslation, err := s.chapterRepo.CreateTranslation(ctx, ct)
	if err != nil {
		logger.Error(err, "failed to create chapter translation")
		return nil, errors.New("unable to create translation")
	}

	return createdTranslation, nil
}

// DeleteTranslation deletes a translation

func (s *ChapterService) DeleteTranslation(ctx context.Context, id string) error {
	affected, err := s.chapterRepo.DeleteTranslation(ctx, id)
	if err != nil {
		logger.Error(err, "failed to delete chapter translation")
		return errors.New("unable to delete translation")
	}
	if affected == 0 {
		return errors.New("chapter translation not found")
	}
	return nil
}
