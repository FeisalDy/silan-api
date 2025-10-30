package gormrepo

import (
	"context"
	"simple-go/internal/domain/role"
	"simple-go/internal/repository"

	"gorm.io/gorm"
)

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) repository.RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) Create(ctx context.Context, role *role.Role) error {
	return r.db.WithContext(ctx).Create(role).Error
}

func (r *roleRepository) GetByID(ctx context.Context, id string) (*role.Role, error) {
	var ro role.Role
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&ro).Error
	if err != nil {
		return nil, err
	}
	return &ro, nil
}

func (r *roleRepository) GetByName(ctx context.Context, name string) (*role.Role, error) {
	var ro role.Role
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&ro).Error
	if err != nil {
		return nil, err
	}
	return &ro, nil
}

func (r *roleRepository) GetAll(ctx context.Context) ([]role.Role, error) {
	var roles []role.Role
	err := r.db.WithContext(ctx).Find(&roles).Error
	return roles, err
}

// GetRoleUsers retrieves all user IDs that have this role
// Note: Uses raw table query to avoid importing user domain
func (r *roleRepository) GetRoleUsers(ctx context.Context, roleID string) ([]string, error) {
	var userIDs []string
	err := r.db.WithContext(ctx).
		Table("user_roles").
		Where("role_id = ?", roleID).
		Pluck("user_id", &userIDs).Error
	return userIDs, err
}
