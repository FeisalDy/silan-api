package gormrepo

import (
	"context"
	"fmt"
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

func (r *chapterRepository) Create(ctx context.Context, c *chapter.Chapter) error {
	return r.db.WithContext(ctx).Create(c).Error
}

func (r *chapterRepository) GetByID(ctx context.Context, id string) (*chapter.Chapter, error) {
	var c chapter.Chapter
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&c).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *chapterRepository) GetByNovel(ctx context.Context, novelID string, limit, offset int) ([]chapter.Chapter, error) {
	var chapters []chapter.Chapter
	query := r.db.WithContext(ctx).Where("novel_id = ?", novelID).Order("number ASC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&chapters).Error
	return chapters, err
}

func (r *chapterRepository) GetByNovelAndNumber(ctx context.Context, novelID string, number int) (*chapter.Chapter, error) {
	var c chapter.Chapter
	err := r.db.WithContext(ctx).
		Where("novel_id = ? AND number = ?", novelID, number).
		First(&c).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *chapterRepository) GetAll(ctx context.Context, limit, offset int) ([]chapter.Chapter, error) {
	var chapters []chapter.Chapter
	query := r.db.WithContext(ctx).Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&chapters).Error
	return chapters, err
}

func (r *chapterRepository) Update(ctx context.Context, c *chapter.Chapter) error {
	result := r.db.WithContext(ctx).Model(c).Where("id = ?", c.ID).Updates(c)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("chapter not found")
	}
	return nil
}

func (r *chapterRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&chapter.Chapter{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("chapter not found")
	}
	return nil
}

func (r *chapterRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&chapter.Chapter{}).Count(&count).Error
	return count, err
}

func (r *chapterRepository) CountByNovel(ctx context.Context, novelID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&chapter.Chapter{}).Where("novel_id = ?", novelID).Count(&count).Error
	return count, err
}

// Translation operations
func (r *chapterRepository) CreateTranslation(ctx context.Context, ct *chapter.ChapterTranslation) error {
	return r.db.WithContext(ctx).Create(ct).Error
}

func (r *chapterRepository) GetTranslation(ctx context.Context, chapterID, lang string) (*chapter.ChapterTranslation, error) {
	var ct chapter.ChapterTranslation
	err := r.db.WithContext(ctx).
		Where("chapter_id = ? AND lang = ?", chapterID, lang).
		First(&ct).Error
	if err != nil {
		return nil, err
	}
	return &ct, nil
}

func (r *chapterRepository) GetTranslationByID(ctx context.Context, id string) (*chapter.ChapterTranslation, error) {
	var ct chapter.ChapterTranslation
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&ct).Error
	if err != nil {
		return nil, err
	}
	return &ct, nil
}

func (r *chapterRepository) GetTranslations(ctx context.Context, chapterID string) ([]chapter.ChapterTranslation, error) {
	var translations []chapter.ChapterTranslation
	err := r.db.WithContext(ctx).
		Where("chapter_id = ?", chapterID).
		Order("lang ASC").
		Find(&translations).Error
	return translations, err
}

func (r *chapterRepository) UpdateTranslation(ctx context.Context, ct *chapter.ChapterTranslation) error {
	result := r.db.WithContext(ctx).Model(ct).Where("id = ?", ct.ID).Updates(ct)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("chapter translation not found")
	}
	return nil
}

func (r *chapterRepository) DeleteTranslation(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&chapter.ChapterTranslation{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("chapter translation not found")
	}
	return nil
}
