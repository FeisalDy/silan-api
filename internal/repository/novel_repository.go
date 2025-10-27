package repository

import (
	"context"
	"simple-go/internal/domain/novel"
)

type NovelRepository interface {
	Create(ctx context.Context, n *novel.Novel) error
	GetByID(ctx context.Context, id string) (*novel.Novel, error)
	GetByIDWithTranslations(ctx context.Context, id, lang string) (*novel.Novel, error)
	GetAll(ctx context.Context, limit, offset int, title, lang string) ([]novel.Novel, error)
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context) (int64, error)

	CreateTranslation(ctx context.Context, nt *novel.NovelTranslation) error
	GetTranslation(ctx context.Context, novelID, lang string) (*novel.NovelTranslation, error)
	DeleteTranslation(ctx context.Context, translationID string) error
}
