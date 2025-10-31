package gormrepo

import (
	"context"
	"simple-go/internal/domain/volume"

	"gorm.io/gorm"
)

type volumeRepository struct {
	db *gorm.DB
}

func NewVolumeRepository(db *gorm.DB) *volumeRepository {
	return &volumeRepository{db: db}
}

func (r *volumeRepository) Create(ctx context.Context, v *volume.Volume) (*volume.Volume, error) {
	if err := r.db.WithContext(ctx).Create(v).Error; err != nil {
		return nil, err
	}
	return v, nil
}

func (r *volumeRepository) CreateTranslation(ctx context.Context, vt *volume.VolumeTranslation) (*volume.VolumeTranslation, error) {
	if err := r.db.WithContext(ctx).Create(vt).Error; err != nil {
		return nil, err
	}
	return vt, nil
}

func (r *volumeRepository) Update(ctx context.Context, v *volume.Volume) (*volume.Volume, error) {
	if err := r.db.WithContext(ctx).
		Model(&volume.Volume{}).
		Where("id = ?", v.ID).
		Updates(v).Error; err != nil {
		return nil, err
	}

	var updated volume.Volume
	if err := r.db.WithContext(ctx).First(&updated, "id = ?", v.ID).Error; err != nil {
		return nil, err
	}

	return &updated, nil
}

func (r *volumeRepository) GetByID(ctx context.Context, id string) (*volume.Volume, error) {
	var v volume.Volume
	err := r.db.WithContext(ctx).Preload("Translations").First(&v, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *volumeRepository) GetAllWithChaptersByNovelID(ctx context.Context, novelID string) ([]volume.Volume, error) {
	var volumes []volume.Volume
	query := r.db.WithContext(ctx).
		Preload("Chapters").
		Preload("Chapters.Translations").
		Preload("Translations").
		Where("novel_id = ?", novelID).
		Order("number ASC")

	if err := query.Find(&volumes).Error; err != nil {
		return nil, err
	}
	return volumes, nil
}

func (r *volumeRepository) GetAllWithChaptersByNovelIDAndLang(ctx context.Context, novelID, lang string) ([]volume.Volume, error) {
	var volumes []volume.Volume
	q := r.db.WithContext(ctx).
		Preload("Chapters").
		Preload("Translations", "lang = ?", lang).
		Preload("Chapters.Translations", "lang = ?", lang).
		Where("novel_id = ?", novelID).
		Order("number ASC")

	if err := q.Find(&volumes).Error; err != nil {
		return nil, err
	}
	return volumes, nil
}

func (r *volumeRepository) Delete(ctx context.Context, id string) (int64, error) {
	result := r.db.WithContext(ctx).Delete(&volume.Volume{}, "id = ?", id)
	return result.RowsAffected, result.Error
}
