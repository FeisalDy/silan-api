package repository

import (
	"context"
)

type NovelTagRepository interface {
	LinkTagsToNovel(ctx context.Context, novelID string, tagIDs []string) error
	UnlinkTagsFromNovel(ctx context.Context, novelID string, tagIDs []string) error
	GetTagIDsByNovelID(ctx context.Context, novelID string) ([]string, error)
}
