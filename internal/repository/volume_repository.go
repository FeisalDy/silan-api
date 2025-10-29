package repository

import (
	"context"
	"simple-go/internal/domain/volume"
)

type VolumeRepository interface {
	Create(ctx context.Context, v *volume.Volume) (*volume.Volume, error)
	Update(ctx context.Context, v *volume.Volume) (*volume.Volume, error)

	GetByID(ctx context.Context, id string) (*volume.Volume, error)
	GetAllWithChaptersByNovelID(ctx context.Context, novelID string) ([]volume.Volume, error)
	GetAllWithChaptersByNovelIDAndLang(ctx context.Context, novelID, lang string) ([]volume.Volume, error)

	Delete(ctx context.Context, id string) (int64, error)
}
