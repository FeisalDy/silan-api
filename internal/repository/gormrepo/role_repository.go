package gormrepo

import (
	"context"
	"fmt"
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

func (r *roleRepository) AssignRoleToUser(ctx context.Context, userID, roleID string) error {
	userRole := &role.UserRole{
		UserID: userID,
		RoleID: roleID,
	}
	return r.db.WithContext(ctx).Create(userRole).Error
}

func (r *roleRepository) RemoveRoleFromUser(ctx context.Context, userID, roleID string) error {
	result := r.db.WithContext(ctx).
		Where("user_id = ? AND role_id = ?", userID, roleID).
		Delete(&role.UserRole{})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("user role not found")
	}
	return nil
}

func (r *roleRepository) GetUserRoles(ctx context.Context, userID string) ([]role.Role, error) {
	var roles []role.Role
	err := r.db.WithContext(ctx).
		Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userID).
		Find(&roles).Error
	return roles, err
}

func (r *roleRepository) GetRoleUsers(ctx context.Context, roleID string) ([]string, error) {
	var userIDs []string
	err := r.db.WithContext(ctx).
		Model(&role.UserRole{}).
		Where("role_id = ?", roleID).
		Pluck("user_id", &userIDs).Error
	return userIDs, err
}
