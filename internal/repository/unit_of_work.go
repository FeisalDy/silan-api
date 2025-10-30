package repository

import (
	"context"
)

// UnitOfWork defines the interface for transaction management
type UnitOfWork interface {
	Do(ctx context.Context, fn func(provider RepositoryProvider) error) error
}

// RepositoryProvider provides access to all repositories within a transaction context
type RepositoryProvider interface {
	User() UserRepository
	Role() RoleRepository
	Novel() NovelRepository
	Chapter() ChapterRepository
	Media() MediaRepository
}
