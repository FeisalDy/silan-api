package gormrepo

import (
	"context"
	"simple-go/internal/repository"

	"gorm.io/gorm"
)

type unitOfWork struct {
	db *gorm.DB
}

func NewUnitOfWork(db *gorm.DB) repository.UnitOfWork {
	return &unitOfWork{db: db}
}

func (u *unitOfWork) Do(ctx context.Context, fn func(provider repository.RepositoryProvider) error) error {
	return u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		provider := &repoProvider{db: tx}
		return fn(provider)
	})
}

type repoProvider struct {
	db *gorm.DB
}

func (rp *repoProvider) User() repository.UserRepository {
	return NewUserRepository(rp.db)
}

func (rp *repoProvider) Role() repository.RoleRepository {
	return NewRoleRepository(rp.db)
}

func (rp *repoProvider) Novel() repository.NovelRepository {
	return NewNovelRepository(rp.db)
}

// func (rp *repoProvider) Chapter() repository.ChapterRepository {
// 	return NewChapterRepository(rp.db)
// }

func (rp *repoProvider) Media() repository.MediaRepository {
	return NewMediaRepository(rp.db)
}
