package repository

import (
	"context"
	"simple-go/internal/domain/chapter"
)

type ChapterRepository interface {
	Create(ctx context.Context, c *chapter.Chapter) (*chapter.Chapter, error)
	GetByID(ctx context.Context, id string) (*chapter.Chapter, error)
	GetByIDAndLang(ctx context.Context, id, lang string) (*chapter.Chapter, error)
	Delete(ctx context.Context, id string) (int64, error)
	CreateTranslation(ctx context.Context, ct *chapter.ChapterTranslation) (*chapter.ChapterTranslation, error)
	GetTranslation(ctx context.Context, chapterID, lang string) (*chapter.ChapterTranslation, error)
	DeleteTranslation(ctx context.Context, translationID string) (int64, error)
}
