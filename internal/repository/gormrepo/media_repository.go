package gormrepo

import (
	"context"
	"simple-go/internal/domain/media"

	"gorm.io/gorm"
)

type MediaRepository struct {
	db *gorm.DB
}

func NewMediaRepository(db *gorm.DB) *MediaRepository {
	return &MediaRepository{db: db}
}

func (r *MediaRepository) Create(ctx context.Context, m *media.Media) (*media.Media, error) {
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return nil, err
	}
	return m, nil
}

func (r *MediaRepository) GetByID(ctx context.Context, id string) (*media.Media, error) {
	var m media.Media

	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&m).Error

	if err != nil {
		return nil, err
	}

	return &m, nil
}

func (r *MediaRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&media.Media{}, "id = ?", id).Error
}

func (r *MediaRepository) Update(ctx context.Context, media *media.Media) error {
	return r.db.WithContext(ctx).Save(media).Error
}
