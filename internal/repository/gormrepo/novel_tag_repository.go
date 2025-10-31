package gormrepo

import (
	"context"
	noveltag "simple-go/internal/domain/novel_tag"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type novelTagRepository struct {
	db *gorm.DB
}

func NewNovelTagRepository(db *gorm.DB) *novelTagRepository {
	return &novelTagRepository{db: db}
}

func (r *novelTagRepository) LinkTagsToNovel(ctx context.Context, novelID string, tagIDs []string) error {
	if len(tagIDs) == 0 {
		return nil
	}

	uniqueIDs := make(map[string]struct{}, len(tagIDs))
	var deduped []string
	for _, id := range tagIDs {
		if _, exists := uniqueIDs[id]; exists {
			continue
		}
		uniqueIDs[id] = struct{}{}
		deduped = append(deduped, id)
	}

	// Create novel_tag records
	novelTags := make([]noveltag.NovelTag, len(deduped))
	for i, tagID := range deduped {
		novelTags[i] = noveltag.NovelTag{
			NovelID: novelID,
			TagID:   tagID,
		}
	}

	// Batch insert, ignoring duplicates
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(&novelTags).Error
}

func (r *novelTagRepository) UnlinkTagsFromNovel(ctx context.Context, novelID string, tagIDs []string) error {
	if len(tagIDs) == 0 {
		return nil
	}

	return r.db.WithContext(ctx).
		Where("novel_id = ? AND tag_id IN ?", novelID, tagIDs).
		Delete(&noveltag.NovelTag{}).Error
}

func (r *novelTagRepository) GetTagIDsByNovelID(ctx context.Context, novelID string) ([]string, error) {
	var tagIDs []string
	err := r.db.WithContext(ctx).
		Model(&noveltag.NovelTag{}).
		Where("novel_id = ?", novelID).
		Pluck("tag_id", &tagIDs).Error

	return tagIDs, err
}
