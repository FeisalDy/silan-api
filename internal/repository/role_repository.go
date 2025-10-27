package repository

import (
	"context"
	"simple-go/internal/domain/role"
)

// RoleRepository defines the interface for role data operations
type RoleRepository interface {
	Create(ctx context.Context, role *role.Role) error
	GetByID(ctx context.Context, id string) (*role.Role, error)
	GetByName(ctx context.Context, name string) (*role.Role, error)
	GetAll(ctx context.Context) ([]role.Role, error)

	// UserRole operations
	AssignRoleToUser(ctx context.Context, userID, roleID string) error
	RemoveRoleFromUser(ctx context.Context, userID, roleID string) error
	GetUserRoles(ctx context.Context, userID string) ([]role.Role, error)
	GetRoleUsers(ctx context.Context, roleID string) ([]string, error) // Returns user IDs
}
