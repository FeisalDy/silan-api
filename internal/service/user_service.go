package service

import (
	"context"
	"errors"
	"fmt"
	"simple-go/internal/domain/user"
	"simple-go/internal/repository"

	"gorm.io/gorm"
)

// UserService handles user-related business logic
type UserService struct {
	userRepo repository.UserRepository
	roleRepo repository.RoleRepository
}

func NewUserService(userRepo repository.UserRepository, roleRepo repository.RoleRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

// GetByID retrieves a user by ID
func (s *UserService) GetByID(ctx context.Context, id string) (*user.UserResponse, error) {
	u, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return u.ToResponse(), nil
}

// GetAll retrieves all users with pagination
func (s *UserService) GetAll(ctx context.Context, limit, offset int) ([]user.UserResponse, int64, error) {
	users, err := s.userRepo.GetAll(ctx, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get users: %w", err)
	}

	count, err := s.userRepo.Count(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	responses := make([]user.UserResponse, len(users))
	for i, u := range users {
		responses[i] = *u.ToResponse()
	}

	return responses, count, nil
}

// Update updates a user
func (s *UserService) Update(ctx context.Context, id string, dto user.UpdateUserDTO) (*user.UserResponse, error) {
	// Get existing user
	u, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Update fields
	if dto.Username != nil {
		u.Username = dto.Username
	}
	if dto.AvatarURL != nil {
		u.AvatarURL = dto.AvatarURL
	}
	if dto.Bio != nil {
		u.Bio = dto.Bio
	}
	if dto.Status != nil {
		u.Status = *dto.Status
	}

	// Save changes
	if err := s.userRepo.Update(ctx, u); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return u.ToResponse(), nil
}

// Delete deletes a user
func (s *UserService) Delete(ctx context.Context, id string) error {
	if err := s.userRepo.Delete(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

// GetUserRoles retrieves all roles assigned to a user
func (s *UserService) GetUserRoles(ctx context.Context, userID string) ([]string, error) {
	roles, err := s.roleRepo.GetUserRoles(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user roles: %w", err)
	}

	roleNames := make([]string, len(roles))
	for i, r := range roles {
		roleNames[i] = r.Name
	}

	return roleNames, nil
}
