package repository

import (
	"context"
	"simple-go/internal/domain/chapter"
)

type ChapterRepository interface {
	Create(ctx context.Context, c *chapter.Chapter) (*chapter.Chapter, error)
	GetByID(ctx context.Context, id string) (*chapter.Chapter, error)
	GetByIDAndLang(ctx context.Context, id, lang string) (*chapter.Chapter, error)
	GetNextChapterID(ctx context.Context, volumeID string, currentNumber int) (*string, error)
	GetPreviousChapterID(ctx context.Context, volumeID string, currentNumber int) (*string, error)
	GetFirstChapterIDOfVolume(ctx context.Context, volumeID string) (*string, error)
	GetLastChapterIDOfVolume(ctx context.Context, volumeID string) (*string, error)
	Delete(ctx context.Context, id string) (int64, error)
	CreateTranslation(ctx context.Context, ct *chapter.ChapterTranslation) (*chapter.ChapterTranslation, error)
	GetTranslation(ctx context.Context, chapterID, lang string) (*chapter.ChapterTranslation, error)
	DeleteTranslation(ctx context.Context, translationID string) (int64, error)
}
