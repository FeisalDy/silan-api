package gormrepo

import (
	"context"
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

func (r *novelRepository) Create(ctx context.Context, n *novel.Novel) (*novel.Novel, error) {
	if err := r.db.WithContext(ctx).Create(n).Error; err != nil {
		return nil, err
	}
	return n, nil
}

func (r *novelRepository) GetByID(ctx context.Context, id string) (*novel.Novel, error) {
	var n novel.Novel
	err := r.db.WithContext(ctx).
		Preload("Media").
		Preload("Genres").
		Preload("Tags").
		Preload("Translations").
		Joins("JOIN novel_translations ON novel_translations.novel_id = novels.id").
		Where("novels.id = ?", id).
		First(&n).Error
	if err != nil {
		return nil, err
	}
	return &n, nil
}

func (r *novelRepository) GetAll(ctx context.Context, limit, offset int, title, lang string) ([]novel.Novel, error) {
	var novels []novel.Novel

	// Base query: load novels and associations without joining to translations to avoid row multiplication
	query := r.db.WithContext(ctx).
		Model(&novel.Novel{}).
		Preload("Media").
		Preload("Translations").
		Preload("Genres").
		Preload("Tags").
		Order("novels.updated_at DESC")

	// Optional title filter: when provided, filter through translations' title
	if title != "" {
		// Join only when filtering by title and ensure distinct novels
		query = query.Joins("JOIN novel_translations ON novel_translations.novel_id = novels.id").
			Where("novel_translations.title ILIKE ?", "%"+title+"%").
			Distinct("novels.id")
	}

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&novels).Error; err != nil {
		return nil, err
	}

	return novels, nil
}

func (r *novelRepository) GetAllByLang(ctx context.Context, lang string, limit, offset int) ([]novel.Novel, error) {
	var novels []novel.Novel

	// Filter novels that have a translation for the specified lang without duplicating rows
	query := r.db.WithContext(ctx).
		Model(&novel.Novel{}).
		Preload("Translations", "lang = ?", lang).
		Preload("Media").
		Preload("Genres").
		Preload("Tags").
		Where("EXISTS (SELECT 1 FROM novel_translations nt WHERE nt.novel_id = novels.id AND nt.lang = ?)", lang).
		Order("novels.updated_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&novels).Error; err != nil {
		return nil, err
	}
	return novels, nil
}

func (r *novelRepository) Delete(ctx context.Context, id string) (int64, error) {
	result := r.db.WithContext(ctx).Delete(&novel.Novel{}, "id = ?", id)
	return result.RowsAffected, result.Error
}

func (r *novelRepository) UpdateCoverMedia(ctx context.Context, novelID, mediaID string) (*novel.Novel, error) {
	if err := r.db.WithContext(ctx).
		Model(&novel.Novel{}).
		Where("id = ?", novelID).
		Update("cover_media_id", mediaID).Error; err != nil {
		return nil, err
	}

	var updated novel.Novel
	if err := r.db.WithContext(ctx).First(&updated, "id = ?", novelID).Error; err != nil {
		return nil, err
	}
	return &updated, nil
}

func (r *novelRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&novel.Novel{}).Count(&count).Error
	return count, err
}

func (r *novelRepository) CreateTranslation(ctx context.Context, nt *novel.NovelTranslation) (*novel.NovelTranslation, error) {
	if err := r.db.WithContext(ctx).Create(nt).Error; err != nil {
		return nil, err
	}
	return nt, nil
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

func (r *novelRepository) DeleteTranslation(ctx context.Context, id string) (int64, error) {
	result := r.db.WithContext(ctx).Delete(&novel.NovelTranslation{}, "id = ?", id)
	return result.RowsAffected, result.Error
}
