package service

import (
	"context"
	"errors"
	"simple-go/internal/domain/user"
	"simple-go/internal/repository"
	"simple-go/pkg/logger"

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

func (s *UserService) GetByID(ctx context.Context, id string) (*user.UserResponse, error) {
	u, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		logger.Error(err, "failed to get user by ID")
		return nil, errors.New("unable to retrieve user")
	}
	return u.ToResponse(), nil
}

// GetAll retrieves all users with pagination
func (s *UserService) GetAll(ctx context.Context, limit, offset int) ([]user.UserResponse, int64, error) {
	users, err := s.userRepo.GetAll(ctx, limit, offset)
	if err != nil {
		logger.Error(err, "failed to get all users")
		return nil, 0, errors.New("unable to retrieve users")
	}

	count, err := s.userRepo.Count(ctx)
	if err != nil {
		logger.Error(err, "failed to count users")
		return nil, 0, errors.New("unable to retrieve users")
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
		logger.Error(err, "failed to get user for update")
		return nil, errors.New("unable to update user")
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
		logger.Error(err, "failed to save user updates")
		return nil, errors.New("unable to update user")
	}

	return u.ToResponse(), nil
}

// Delete deletes a user
func (s *UserService) Delete(ctx context.Context, id string) error {
	if err := s.userRepo.Delete(ctx, id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		logger.Error(err, "failed to delete user")
		return errors.New("unable to delete user")
	}
	return nil
}

// GetUserRoles retrieves all roles assigned to a user
func (s *UserService) GetUserRoles(ctx context.Context, userID string) ([]string, error) {
	roles, err := s.userRepo.GetRoles(ctx, userID)
	if err != nil {
		logger.Error(err, "failed to get user roles")
		return nil, errors.New("unable to retrieve user roles")
	}

	roleNames := make([]string, len(roles))
	for i, r := range roles {
		roleNames[i] = r.Name
	}

	return roleNames, nil
}
