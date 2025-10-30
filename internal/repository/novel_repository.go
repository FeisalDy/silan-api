package repository

import (
	"context"
	"simple-go/internal/domain/novel"
)

type NovelRepository interface {
	Create(ctx context.Context, n *novel.Novel) (*novel.Novel, error)
	GetByID(ctx context.Context, id string) (*novel.Novel, error)
	GetAll(ctx context.Context, limit, offset int, title, lang string) ([]novel.Novel, error)
	GetAllByLang(ctx context.Context, lang string, limit, offset int) ([]novel.Novel, error)
	Delete(ctx context.Context, id string) (int64, error)
	UpdateCoverMedia(ctx context.Context, novelID, mediaID string) (*novel.Novel, error)
	Count(ctx context.Context) (int64, error)

	CreateTranslation(ctx context.Context, nt *novel.NovelTranslation) (*novel.NovelTranslation, error)
	GetTranslation(ctx context.Context, novelID, lang string) (*novel.NovelTranslation, error)
	DeleteTranslation(ctx context.Context, translationID string) (int64, error)
}
