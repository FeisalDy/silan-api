// package service

// import (
// 	"context"
// 	"errors"
// 	"simple-go/internal/domain/chapter"
// 	"simple-go/internal/repository"
// 	"simple-go/pkg/logger"

// 	"gorm.io/gorm"
// )

// // ChapterService handles chapter-related business logic
// type ChapterService struct {
// 	uow         repository.UnitOfWork
// 	chapterRepo repository.ChapterRepository
// 	novelRepo   repository.NovelRepository
// }

// func NewChapterService(
// 	uow repository.UnitOfWork,
// 	chapterRepo repository.ChapterRepository,
// 	novelRepo repository.NovelRepository,
// ) *ChapterService {
// 	return &ChapterService{
// 		uow:         uow,
// 		chapterRepo: chapterRepo,
// 		novelRepo:   novelRepo,
// 	}
// }

// // CreateChapterWithTranslation creates a chapter and its initial translation in one transaction
// func (s *ChapterService) CreateChapterWithTranslation(
// 	ctx context.Context,
// 	translatorID string,
// 	chapterDTO chapter.CreateChapterDTO,
// 	translationDTO chapter.CreateChapterTranslationDTO,
// ) (*chapter.Chapter, *chapter.ChapterTranslation, error) {
// 	var newChapter *chapter.Chapter
// 	var newTranslation *chapter.ChapterTranslation

// 	err := s.uow.Do(ctx, func(provider repository.RepositoryProvider) error {
// 		// Verify novel exists
// 		_, err := provider.Novel().GetByID(ctx, chapterDTO.VolumeID)
// 		if err != nil {
// 			if errors.Is(err, gorm.ErrRecordNotFound) {
// 				return errors.New("novel not found")
// 			}
// 			logger.Error(err, "failed to verify novel exists")
// 			return errors.New("unable to create chapter")
// 		}

// 		// Check if chapter with this number already exists
// 		existing, err := provider.Chapter().GetByVolumeAndNumber(ctx, chapterDTO.VolumeID, chapterDTO.Number)
// 		if err == nil && existing != nil {
// 			return errors.New("chapter with this number already exists for this volume")
// 		}

// 		// Create chapter
// 		newChapter = &chapter.Chapter{
// 			VolumeID:  chapterDTO.VolumeID,
// 			Number:    chapterDTO.Number,
// 			WordCount: chapterDTO.WordCount,
// 		}

// 		if err := provider.Chapter().Create(ctx, newChapter); err != nil {
// 			logger.Error(err, "failed to create chapter in database")
// 			return errors.New("unable to create chapter")
// 		}

// 		// Create translation
// 		newTranslation = &chapter.ChapterTranslation{
// 			ChapterID:    newChapter.ID,
// 			Lang:         translationDTO.Lang,
// 			Title:        translationDTO.Title,
// 			Content:      translationDTO.Content,
// 			TranslatorID: translatorID,
// 		}

// 		if err := provider.Chapter().CreateTranslation(ctx, newTranslation); err != nil {
// 			logger.Error(err, "failed to create chapter translation")
// 			return errors.New("unable to create chapter translation")
// 		}

// 		return nil
// 	})

// 	if err != nil {
// 		return nil, nil, err
// 	}

// 	return newChapter, newTranslation, nil
// }

// // GetByID retrieves a chapter by ID
// func (s *ChapterService) GetByID(ctx context.Context, id string) (*chapter.Chapter, error) {
// 	c, err := s.chapterRepo.GetByID(ctx, id)
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, errors.New("chapter not found")
// 		}
// 		logger.Error(err, "failed to get chapter by ID")
// 		return nil, errors.New("unable to retrieve chapter")
// 	}
// 	return c, nil
// }

// // GetByIDWithTranslations retrieves a chapter with all its translations
// func (s *ChapterService) GetByIDWithTranslations(ctx context.Context, id string) (*chapter.ChapterResponse, error) {
// 	c, err := s.chapterRepo.GetByID(ctx, id)
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, errors.New("chapter not found")
// 		}
// 		logger.Error(err, "failed to get chapter by ID")
// 		return nil, errors.New("unable to retrieve chapter")
// 	}

// 	translations, err := s.chapterRepo.GetTranslations(ctx, id)
// 	if err != nil {
// 		logger.Error(err, "failed to get chapter translations")
// 		return nil, errors.New("unable to retrieve chapter translations")
// 	}

// 	return &chapter.ChapterResponse{
// 		ID:           c.ID,
// 		VolumeID:     c.VolumeID,
// 		Number:       c.Number,
// 		WordCount:    c.WordCount,
// 		Translations: translations,
// 		CreatedAt:    c.CreatedAt,
// 		UpdatedAt:    c.UpdatedAt,
// 	}, nil
// }

// // GetByNovel retrieves all chapters for a novel
// func (s *ChapterService) GetByNovel(ctx context.Context, novelID string, limit, offset int) ([]chapter.Chapter, int64, error) {
// 	// Verify novel exists
// 	_, err := s.novelRepo.GetByID(ctx, novelID)
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, 0, errors.New("novel not found")
// 		}
// 		logger.Error(err, "failed to verify novel exists")
// 		return nil, 0, errors.New("unable to retrieve chapters")
// 	}

// 	chapters, err := s.chapterRepo.GetByNovel(ctx, novelID, limit, offset)
// 	if err != nil {
// 		logger.Error(err, "failed to get chapters by novel")
// 		return nil, 0, errors.New("unable to retrieve chapters")
// 	}

// 	count, err := s.chapterRepo.CountByNovel(ctx, novelID)
// 	if err != nil {
// 		logger.Error(err, "failed to count chapters")
// 		return nil, 0, errors.New("unable to retrieve chapters")
// 	}

// 	return chapters, count, nil
// }

// // GetByVolumeAndNumber retrieves a specific chapter by volume and number
// func (s *ChapterService) GetByVolumeAndNumber(ctx context.Context, volumeID string, number int) (*chapter.Chapter, error) {
// 	c, err := s.chapterRepo.GetByVolumeAndNumber(ctx, volumeID, number)
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, errors.New("chapter not found")
// 		}
// 		logger.Error(err, "failed to get chapter by volume and number")
// 		return nil, errors.New("unable to retrieve chapter")
// 	}
// 	return c, nil
// }

// // Update updates a chapter
// func (s *ChapterService) Update(ctx context.Context, id string, dto chapter.UpdateChapterDTO) (*chapter.Chapter, error) {
// 	c, err := s.chapterRepo.GetByID(ctx, id)
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, errors.New("chapter not found")
// 		}
// 		logger.Error(err, "failed to get chapter for update")
// 		return nil, errors.New("unable to update chapter")
// 	}

// 	// Update fields
// 	if dto.Number != nil {
// 		// Check if new number conflicts
// 		existing, err := s.chapterRepo.GetByVolumeAndNumber(ctx, c.VolumeID, *dto.Number)
// 		if err == nil && existing != nil && existing.ID != id {
// 			return nil, errors.New("chapter with this number already exists")
// 		}
// 		c.Number = *dto.Number
// 	}
// 	if dto.WordCount != nil {
// 		c.WordCount = dto.WordCount
// 	}

// 	if err := s.chapterRepo.Update(ctx, c); err != nil {
// 		logger.Error(err, "failed to save chapter updates")
// 		return nil, errors.New("unable to update chapter")
// 	}

// 	return c, nil
// }

// // Delete deletes a chapter
// func (s *ChapterService) Delete(ctx context.Context, id string) error {
// 	if err := s.chapterRepo.Delete(ctx, id); err != nil {
// 		logger.Error(err, "failed to delete chapter")
// 		return errors.New("unable to delete chapter")
// 	}
// 	return nil
// }

// // Translation operations

// // CreateTranslation creates a new translation for a chapter
// func (s *ChapterService) CreateTranslation(ctx context.Context, translatorID string, dto chapter.CreateChapterTranslationDTO) (*chapter.ChapterTranslation, error) {
// 	// Verify chapter exists
// 	_, err := s.chapterRepo.GetByID(ctx, dto.ChapterID)
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, errors.New("chapter not found")
// 		}
// 		logger.Error(err, "failed to verify chapter exists")
// 		return nil, errors.New("unable to create translation")
// 	}

// 	// Check if translation already exists
// 	existing, err := s.chapterRepo.GetTranslation(ctx, dto.ChapterID, dto.Lang)
// 	if err == nil && existing != nil {
// 		return nil, errors.New("translation for this language already exists")
// 	}

// 	ct := &chapter.ChapterTranslation{
// 		ChapterID:    dto.ChapterID,
// 		Lang:         dto.Lang,
// 		Title:        dto.Title,
// 		Content:      dto.Content,
// 		TranslatorID: translatorID,
// 	}

// 	if err := s.chapterRepo.CreateTranslation(ctx, ct); err != nil {
// 		logger.Error(err, "failed to create chapter translation")
// 		return nil, errors.New("unable to create translation")
// 	}

// 	return ct, nil
// }

// // GetTranslation retrieves a specific translation
// func (s *ChapterService) GetTranslation(ctx context.Context, chapterID, lang string) (*chapter.ChapterTranslation, error) {
// 	ct, err := s.chapterRepo.GetTranslation(ctx, chapterID, lang)
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, errors.New("translation not found")
// 		}
// 		logger.Error(err, "failed to get chapter translation")
// 		return nil, errors.New("unable to retrieve translation")
// 	}
// 	return ct, nil
// }

// // UpdateTranslation updates a translation
// func (s *ChapterService) UpdateTranslation(ctx context.Context, id string, dto chapter.UpdateChapterTranslationDTO) (*chapter.ChapterTranslation, error) {
// 	ct, err := s.chapterRepo.GetTranslationByID(ctx, id)
// 	if err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, errors.New("translation not found")
// 		}
// 		logger.Error(err, "failed to get translation for update")
// 		return nil, errors.New("unable to update translation")
// 	}

// 	// Update fields
// 	if dto.Title != nil {
// 		ct.Title = *dto.Title
// 	}
// 	if dto.Content != nil {
// 		ct.Content = *dto.Content
// 	}

// 	if err := s.chapterRepo.UpdateTranslation(ctx, ct); err != nil {
// 		logger.Error(err, "failed to save translation updates")
// 		return nil, errors.New("unable to update translation")
// 	}

// 	return ct, nil
// }

// // DeleteTranslation deletes a translation
//
//	func (s *ChapterService) DeleteTranslation(ctx context.Context, id string) error {
//		if err := s.chapterRepo.DeleteTranslation(ctx, id); err != nil {
//			logger.Error(err, "failed to delete chapter translation")
//			return errors.New("unable to delete translation")
//		}
//		return nil
//	}
package service
