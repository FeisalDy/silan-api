package repository

import (
	"context"
	"simple-go/internal/domain/role"
	"simple-go/internal/domain/user"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	Create(ctx context.Context, user *user.User) error
	GetByID(ctx context.Context, id string) (*user.User, error)
	GetByEmail(ctx context.Context, email string) (*user.User, error)
	GetByUsername(ctx context.Context, username string) (*user.User, error)
	GetAll(ctx context.Context, limit, offset int) ([]user.User, error)
	Update(ctx context.Context, user *user.User) error
	Delete(ctx context.Context, id string) error
	Count(ctx context.Context) (int64, error)

	// Role association operations - User owns the relationship
	AssignRole(ctx context.Context, userID, roleID string) error
	RemoveRole(ctx context.Context, userID, roleID string) error
	GetRoles(ctx context.Context, userID string) ([]role.Role, error)
}
