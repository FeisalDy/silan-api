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

	// GetRoleUsers returns user IDs that have this role
	// Note: For assigning roles to users, use UserRepository.AssignRole()
	GetRoleUsers(ctx context.Context, roleID string) ([]string, error)
}
