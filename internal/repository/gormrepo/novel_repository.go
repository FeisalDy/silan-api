package gormrepo

import (
	"context"
	"fmt"
	"simple-go/internal/domain/novel"
	"simple-go/internal/repository"

	"gorm.io/gorm"
)

type novelRepository struct {
	db *gorm.DB
}

func NewNovelRepository(db *gorm.DB) repository.NovelRepository {
	return &novelRepository{db: db}
}

func (r *novelRepository) Create(ctx context.Context, n *novel.Novel) error {
	return r.db.WithContext(ctx).Create(n).Error
}

func (r *novelRepository) GetByID(ctx context.Context, id string) (*novel.Novel, error) {
	var n novel.Novel
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&n).Error
	if err != nil {
		return nil, err
	}
	return &n, nil
}

func (r *novelRepository) GetByCreator(ctx context.Context, creatorID string, limit, offset int) ([]novel.Novel, error) {
	var novels []novel.Novel
	query := r.db.WithContext(ctx).Where("created_by = ?", creatorID).Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&novels).Error
	return novels, err
}

func (r *novelRepository) GetAll(ctx context.Context, limit, offset int) ([]novel.Novel, error) {
	var novels []novel.Novel
	query := r.db.WithContext(ctx).Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&novels).Error
	return novels, err
}

func (r *novelRepository) Update(ctx context.Context, n *novel.Novel) error {
	result := r.db.WithContext(ctx).Model(n).Where("id = ?", n.ID).Updates(n)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("novel not found")
	}
	return nil
}

func (r *novelRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&novel.Novel{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("novel not found")
	}
	return nil
}

func (r *novelRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&novel.Novel{}).Count(&count).Error
	return count, err
}

// Translation operations
func (r *novelRepository) CreateTranslation(ctx context.Context, nt *novel.NovelTranslation) error {
	return r.db.WithContext(ctx).Create(nt).Error
}

func (r *novelRepository) GetTranslation(ctx context.Context, novelID, lang string) (*novel.NovelTranslation, error) {
	var nt novel.NovelTranslation
	err := r.db.WithContext(ctx).
		Where("novel_id = ? AND lang = ?", novelID, lang).
		First(&nt).Error
	if err != nil {
		return nil, err
	}
	return &nt, nil
}

func (r *novelRepository) GetTranslationByID(ctx context.Context, id string) (*novel.NovelTranslation, error) {
	var nt novel.NovelTranslation
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&nt).Error
	if err != nil {
		return nil, err
	}
	return &nt, nil
}

func (r *novelRepository) GetTranslations(ctx context.Context, novelID string) ([]novel.NovelTranslation, error) {
	var translations []novel.NovelTranslation
	err := r.db.WithContext(ctx).
		Where("novel_id = ?", novelID).
		Order("lang ASC").
		Find(&translations).Error
	return translations, err
}

func (r *novelRepository) UpdateTranslation(ctx context.Context, nt *novel.NovelTranslation) error {
	result := r.db.WithContext(ctx).Model(nt).Where("id = ?", nt.ID).Updates(nt)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("novel translation not found")
	}
	return nil
}

func (r *novelRepository) DeleteTranslation(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&novel.NovelTranslation{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("novel translation not found")
	}
	return nil
}
