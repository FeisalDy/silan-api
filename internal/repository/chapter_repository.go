package repository

import (
	"context"
	"simple-go/internal/domain/chapter"
)

// ChapterRepository defines the interface for chapter data operations
type ChapterRepository interface {
	Create(ctx context.Context, c *chapter.Chapter) error
	GetByID(ctx context.Context, id string) (*chapter.Chapter, error)
	GetByNovel(ctx context.Context, novelID string, limit, offset int) ([]chapter.Chapter, error)
	GetByNovelAndNumber(ctx context.Context, novelID string, number int) (*chapter.Chapter, error)
	GetAll(ctx context.Context, limit, offset int) ([]chapter.Chapter, error)
	Update(ctx context.Context, c *chapter.Chapter) error
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context) (int64, error)
	CountByNovel(ctx context.Context, novelID string) (int64, error)

	// Translation operations
	CreateTranslation(ctx context.Context, ct *chapter.ChapterTranslation) error
	GetTranslation(ctx context.Context, chapterID, lang string) (*chapter.ChapterTranslation, error)
	GetTranslationByID(ctx context.Context, id string) (*chapter.ChapterTranslation, error)
	GetTranslations(ctx context.Context, chapterID string) ([]chapter.ChapterTranslation, error)
	UpdateTranslation(ctx context.Context, ct *chapter.ChapterTranslation) error
	DeleteTranslation(ctx context.Context, id string) error
}
