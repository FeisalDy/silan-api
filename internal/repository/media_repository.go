package repository

import (
	"context"
	"simple-go/internal/domain/media"
)

type MediaRepository interface {
	Create(ctx context.Context, media *media.Media) (*media.Media, error)
	GetByID(ctx context.Context, id string) (*media.Media, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, media *media.Media) error
}
