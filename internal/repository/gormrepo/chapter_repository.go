package gormrepo

import (
	"context"
	"simple-go/internal/domain/chapter"
	"simple-go/internal/repository"

	"gorm.io/gorm"
)

type chapterRepository struct {
	db *gorm.DB
}

func NewChapterRepository(db *gorm.DB) repository.ChapterRepository {
	return &chapterRepository{db: db}
}

func (r *chapterRepository) Create(ctx context.Context, c *chapter.Chapter) (*chapter.Chapter, error) {
	if err := r.db.WithContext(ctx).Create(c).Error; err != nil {
		return nil, err
	}
	return c, nil
}

func (r *chapterRepository) GetByID(ctx context.Context, id string) (*chapter.Chapter, error) {
	var c chapter.Chapter
	err := r.db.WithContext(ctx).Preload("Translations").First(&c, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *chapterRepository) GetByIDAndLang(ctx context.Context, id, lang string) (*chapter.Chapter, error) {
	var c chapter.Chapter
	err := r.db.WithContext(ctx).Preload("Translations", "lang = ?", lang).First(&c, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *chapterRepository) Delete(ctx context.Context, id string) (int64, error) {
	result := r.db.WithContext(ctx).Delete(&chapter.Chapter{}, "id = ?", id)
	return result.RowsAffected, result.Error
}

func (r *chapterRepository) CreateTranslation(ctx context.Context, ct *chapter.ChapterTranslation) (*chapter.ChapterTranslation, error) {
	if err := r.db.WithContext(ctx).Create(ct).Error; err != nil {
		return nil, err
	}

	return ct, nil
}

func (r *chapterRepository) GetTranslation(ctx context.Context, chapterID, lang string) (*chapter.ChapterTranslation, error) {
	var ct chapter.ChapterTranslation
	err := r.db.WithContext(ctx).First(&ct, "chapter_id = ? AND lang = ?", chapterID, lang).Error
	if err != nil {
		return nil, err
	}
	return &ct, nil
}

func (r *chapterRepository) DeleteTranslation(ctx context.Context, translationID string) (int64, error) {
	result := r.db.WithContext(ctx).Delete(&chapter.ChapterTranslation{}, "id = ?", translationID)
	return result.RowsAffected, result.Error
}
