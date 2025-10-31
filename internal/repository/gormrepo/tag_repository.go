package gormrepo

import (
	"context"
	"errors"
	"simple-go/internal/domain/tag"
	"strings"

	"gorm.io/gorm"
)

type tagRepository struct {
	db *gorm.DB
}

func NewTagRepository(db *gorm.DB) *tagRepository {
	return &tagRepository{db: db}
}

func (r *tagRepository) Create(ctx context.Context, t *tag.Tag) (*tag.Tag, error) {
	if err := r.db.WithContext(ctx).Create(t).Error; err != nil {
		return nil, err
	}
	return t, nil
}

func (r *tagRepository) GetByName(ctx context.Context, name string) (*tag.Tag, error) {
	var t tag.Tag
	err := r.db.WithContext(ctx).Where("LOWER(name) = LOWER(?)", name).First(&t).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *tagRepository) GetByNames(ctx context.Context, names []string) ([]tag.Tag, error) {
	var tags []tag.Tag
	if len(names) == 0 {
		return tags, nil
	}

	// Convert names to lowercase for case-insensitive search
	lowerNames := make([]string, len(names))
	for i, name := range names {
		lowerNames[i] = strings.ToLower(name)
	}

	err := r.db.WithContext(ctx).
		Where("LOWER(name) IN ?", lowerNames).
		Find(&tags).Error

	return tags, err
}

func (r *tagRepository) GetBySlug(ctx context.Context, slug string) (*tag.Tag, error) {
	var t tag.Tag
	if err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&t).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *tagRepository) FindOrCreateByNames(ctx context.Context, names []string) ([]tag.Tag, error) {
	processed := uniqueNormalizedNames(names)
	if len(processed) == 0 {
		return []tag.Tag{}, nil
	}

	existing, err := r.fetchExistingTags(ctx, processed)
	if err != nil {
		return nil, err
	}

	var result []tag.Tag
	for _, info := range processed {
		if t, ok := existing.byName[info.lowerName]; ok {
			result = append(result, t)
			continue
		}
		if t, ok := existing.bySlug[info.slug]; ok {
			result = append(result, t)
			continue
		}

		newTag := &tag.Tag{
			Name: info.original,
			Slug: info.slug,
		}

		createErr := r.db.WithContext(ctx).Create(newTag).Error
		if createErr != nil {
			// If the tag already exists (slug or name), try to fetch it again
			if fetched, lookupErr := r.lookupExistingTag(ctx, info); lookupErr == nil {
				existing.add(*fetched)
				result = append(result, *fetched)
				continue
			}
			return nil, createErr
		}

		existing.add(*newTag)
		result = append(result, *newTag)
	}

	return result, nil
}

// generateSlug creates a URL-friendly slug from a tag name
func generateSlug(name string) string {
	slug := strings.ToLower(name)
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")
	// Remove special characters (keep alphanumeric and hyphens)
	var result strings.Builder
	for _, r := range slug {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}
	return result.String()
}

type normalizedTagName struct {
	original  string
	lowerName string
	slug      string
}

type tagLookupCache struct {
	byName map[string]tag.Tag
	bySlug map[string]tag.Tag
}

func (c *tagLookupCache) add(t tag.Tag) {
	if c.byName == nil {
		c.byName = make(map[string]tag.Tag)
	}
	if c.bySlug == nil {
		c.bySlug = make(map[string]tag.Tag)
	}
	c.byName[strings.ToLower(t.Name)] = t
	c.bySlug[t.Slug] = t
}

func uniqueNormalizedNames(names []string) []normalizedTagName {
	seen := make(map[string]struct{})
	var result []normalizedTagName
	for _, name := range names {
		trimmed := strings.TrimSpace(name)
		if trimmed == "" {
			continue
		}
		lower := strings.ToLower(trimmed)
		if _, exists := seen[lower]; exists {
			continue
		}
		seen[lower] = struct{}{}
		result = append(result, normalizedTagName{
			original:  trimmed,
			lowerName: lower,
			slug:      generateSlug(trimmed),
		})
	}
	return result
}

func (r *tagRepository) fetchExistingTags(ctx context.Context, names []normalizedTagName) (*tagLookupCache, error) {
	cache := &tagLookupCache{
		byName: make(map[string]tag.Tag),
		bySlug: make(map[string]tag.Tag),
	}

	if len(names) == 0 {
		return cache, nil
	}

	lowerNames := make([]string, len(names))
	slugs := make([]string, len(names))
	for i, info := range names {
		lowerNames[i] = info.lowerName
		slugs[i] = info.slug
	}

	var existing []tag.Tag
	query := r.db.WithContext(ctx).
		Where("LOWER(name) IN ? OR slug IN ?", lowerNames, slugs)
	if err := query.Find(&existing).Error; err != nil {
		return nil, err
	}

	for _, t := range existing {
		cache.add(t)
	}

	return cache, nil
}

func (r *tagRepository) lookupExistingTag(ctx context.Context, info normalizedTagName) (*tag.Tag, error) {
	if tagByName, err := r.GetByName(ctx, info.original); err == nil {
		return tagByName, nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if tagBySlug, err := r.GetBySlug(ctx, info.slug); err == nil {
		return tagBySlug, nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return nil, gorm.ErrRecordNotFound
}
