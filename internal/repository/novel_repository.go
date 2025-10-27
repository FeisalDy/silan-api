package repository

import (
	"context"
	"simple-go/internal/domain/novel"
)

// NovelRepository defines the interface for novel data operations
type NovelRepository interface {
	Create(ctx context.Context, n *novel.Novel) error
	GetByID(ctx context.Context, id string) (*novel.Novel, error)
	GetByCreator(ctx context.Context, creatorID string, limit, offset int) ([]novel.Novel, error)
	GetAll(ctx context.Context, limit, offset int) ([]novel.Novel, error)
	Update(ctx context.Context, n *novel.Novel) error
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context) (int64, error)

	// Translation operations
	CreateTranslation(ctx context.Context, nt *novel.NovelTranslation) error
	GetTranslation(ctx context.Context, novelID, lang string) (*novel.NovelTranslation, error)
	GetTranslationByID(ctx context.Context, id string) (*novel.NovelTranslation, error)
	GetTranslations(ctx context.Context, novelID string) ([]novel.NovelTranslation, error)
	UpdateTranslation(ctx context.Context, nt *novel.NovelTranslation) error
	DeleteTranslation(ctx context.Context, id string) error
}
