package repository

import (
	"context"
	"simple-go/internal/domain/tag"
)

type TagRepository interface {
	Create(ctx context.Context, t *tag.Tag) (*tag.Tag, error)
	GetByName(ctx context.Context, name string) (*tag.Tag, error)
	GetByNames(ctx context.Context, names []string) ([]tag.Tag, error)
	GetBySlug(ctx context.Context, slug string) (*tag.Tag, error)
	FindOrCreateByNames(ctx context.Context, names []string) ([]tag.Tag, error)
}
