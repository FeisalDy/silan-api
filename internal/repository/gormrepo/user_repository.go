package gormrepo

import (
	"context"
	"fmt"
	"simple-go/internal/domain/role"
	"simple-go/internal/domain/user"
	"simple-go/internal/repository"

	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, u *user.User) error {
	return r.db.WithContext(ctx).Create(u).Error
}

func (r *userRepository) GetByID(ctx context.Context, id string) (*user.User, error) {
	var u user.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*user.User, error) {
	var u user.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) GetByUsername(ctx context.Context, username string) (*user.User, error) {
	var u user.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&u).Error
	if err != nil {
		return nil, err
	}
	return &u, nil
}
func (r *userRepository) GetAll(ctx context.Context, limit, offset int) ([]user.User, error) {
	var users []user.User
	query := r.db.WithContext(ctx).Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	err := query.Find(&users).Error
	return users, err
}

func (r *userRepository) Update(ctx context.Context, u *user.User) error {
	result := r.db.WithContext(ctx).Model(u).Where("id = ?", u.ID).Updates(u)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

func (r *userRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&user.User{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}

func (r *userRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&user.User{}).Count(&count).Error
	return count, err
}

// AssignRole assigns a role to a user using GORM M2M association
func (r *userRepository) AssignRole(ctx context.Context, userID, roleID string) error {
	var u user.User
	if err := r.db.WithContext(ctx).First(&u, "id = ?", userID).Error; err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	var ro role.Role
	if err := r.db.WithContext(ctx).First(&ro, "id = ?", roleID).Error; err != nil {
		return fmt.Errorf("role not found: %w", err)
	}

	return r.db.WithContext(ctx).Model(&u).Association("Roles").Append(&ro)
}

// RemoveRole removes a role from a user using GORM M2M association
func (r *userRepository) RemoveRole(ctx context.Context, userID, roleID string) error {
	var u user.User
	if err := r.db.WithContext(ctx).First(&u, "id = ?", userID).Error; err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	var ro role.Role
	if err := r.db.WithContext(ctx).First(&ro, "id = ?", roleID).Error; err != nil {
		return fmt.Errorf("role not found: %w", err)
	}

	return r.db.WithContext(ctx).Model(&u).Association("Roles").Delete(&ro)
}

// GetRoles retrieves all roles assigned to a user using GORM M2M association
func (r *userRepository) GetRoles(ctx context.Context, userID string) ([]role.Role, error) {
	var u user.User
	if err := r.db.WithContext(ctx).Preload("Roles").First(&u, "id = ?", userID).Error; err != nil {
		return nil, err
	}
	return u.Roles, nil
}
