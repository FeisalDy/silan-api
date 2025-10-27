package service

import (
	"context"
	"errors"
	"fmt"
	"simple-go/internal/domain/role"
	"simple-go/internal/domain/user"
	"simple-go/internal/repository"
	"simple-go/pkg/auth"

	"gorm.io/gorm"
)

// AuthService handles authentication-related business logic
type AuthService struct {
	uow        repository.UnitOfWork
	userRepo   repository.UserRepository
	roleRepo   repository.RoleRepository
	jwtManager *auth.JWTManager
}

func NewAuthService(
	uow repository.UnitOfWork,
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
	jwtManager *auth.JWTManager,
) *AuthService {
	return &AuthService{
		uow:        uow,
		userRepo:   userRepo,
		roleRepo:   roleRepo,
		jwtManager: jwtManager,
	}
}

// RegisterRequest represents registration data
type RegisterRequest struct {
	Email    string  `json:"email" binding:"required,email"`
	Password string  `json:"password" binding:"required,min=6"`
	Username *string `json:"username"`
}

// LoginRequest represents login credentials
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse represents authentication response
type AuthResponse struct {
	Token string             `json:"token"`
	User  *user.UserResponse `json:"user"`
}

// Register creates a new user with default "user" role
func (s *AuthService) Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	var newUser *user.User
	var token string

	// Use transaction to create user and assign role
	err = s.uow.Do(ctx, func(provider repository.RepositoryProvider) error {
		// Create user
		newUser = &user.User{
			Email:    req.Email,
			Password: &hashedPassword,
			Username: req.Username,
			Status:   "active",
		}

		if err := provider.User().Create(ctx, newUser); err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}

		// Get default "user" role
		userRole, err := provider.Role().GetByName(ctx, role.RoleUser)
		if err != nil {
			return fmt.Errorf("failed to get default role: %w", err)
		}

		// Assign role to user
		if err := provider.Role().AssignRoleToUser(ctx, newUser.ID, userRole.ID); err != nil {
			return fmt.Errorf("failed to assign role: %w", err)
		}

		// Generate JWT token
		token, err = s.jwtManager.GenerateToken(newUser.ID, newUser.Email)
		if err != nil {
			return fmt.Errorf("failed to generate token: %w", err)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token: token,
		User:  newUser.ToResponse(),
	}, nil
}

// Login authenticates a user and returns a JWT token
func (s *AuthService) Login(ctx context.Context, req LoginRequest) (*AuthResponse, error) {
	// Find user by email
	u, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid email or password")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Check if user has a password (OAuth users might not have one)
	if u.Password == nil {
		return nil, errors.New("invalid email or password")
	}

	// Verify password
	if !auth.CheckPassword(*u.Password, req.Password) {
		return nil, errors.New("invalid email or password")
	}

	// Check user status
	if u.Status != "active" {
		return nil, errors.New("user account is not active")
	}

	// Generate JWT token
	token, err := s.jwtManager.GenerateToken(u.ID, u.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	return &AuthResponse{
		Token: token,
		User:  u.ToResponse(),
	}, nil
}
