package service

import (
	"context"
	"errors"
	"simple-go/internal/domain/role"
	"simple-go/internal/domain/user"
	"simple-go/internal/repository"
	"simple-go/pkg/auth"
	"simple-go/pkg/logger"

	"gorm.io/gorm"
)

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

type RegisterRequest struct {
	Email    string  `json:"email" binding:"required,email"`
	Password string  `json:"password" binding:"required,min=6"`
	Username *string `json:"username"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string             `json:"token"`
	User  *user.UserResponse `json:"user"`
}

func (s *AuthService) Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error) {
	existingUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Error(err, "failed to check existing user")
		return nil, errors.New("unable to process registration request")
	}

	if req.Username != nil && *req.Username != "" {
		existingUsername, err := s.userRepo.GetByUsername(ctx, *req.Username)
		if err == nil && existingUsername != nil {
			return nil, errors.New("user with this username already exists")
		}

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error(err, "failed to check existing username")
			return nil, errors.New("unable to process registration request")
		}
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		logger.Error(err, "failed to hash password")
		return nil, errors.New("unable to process registration request")
	}

	var newUser *user.User
	var token string

	err = s.uow.Do(ctx, func(provider repository.RepositoryProvider) error {
		newUser = &user.User{
			Email:    req.Email,
			Password: &hashedPassword,
			Username: req.Username,
			Status:   "active",
		}

		if err := provider.User().Create(ctx, newUser); err != nil {
			logger.Error(err, "failed to create user in database")
			return errors.New("unable to create user account")
		}

		userRole, err := provider.Role().GetByName(ctx, role.RoleUser)
		if err != nil {
			logger.Error(err, "failed to get default role")
			return errors.New("unable to setup user account")
		}

		if err := provider.Role().AssignRoleToUser(ctx, newUser.ID, userRole.ID); err != nil {
			logger.Error(err, "failed to assign role to user")
			return errors.New("unable to setup user account")
		}

		// Generate JWT token
		token, err = s.jwtManager.GenerateToken(newUser.ID, newUser.Email)
		if err != nil {
			logger.Error(err, "failed to generate JWT token")
			return errors.New("unable to complete registration")
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

func (s *AuthService) Login(ctx context.Context, req LoginRequest) (*AuthResponse, error) {
	u, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Warn("login attempt with non-existent email")
			return nil, errors.New("invalid email or password")
		}
		logger.Error(err, "failed to get user by email")
		return nil, errors.New("unable to process login request")
	}

	// Check if user has a password (OAuth users might not have one)
	if u.Password == nil {
		logger.Warn("login attempt for user without password")
		return nil, errors.New("invalid email or password")
	}

	// Verify password
	if !auth.CheckPassword(*u.Password, req.Password) {
		logger.Warn("failed login attempt - invalid password")
		return nil, errors.New("invalid email or password")
	}

	// Check user status
	if u.Status != "active" {
		logger.Warn("login attempt for inactive user")
		return nil, errors.New("user account is not active")
	}

	// Generate JWT token
	token, err := s.jwtManager.GenerateToken(u.ID, u.Email)
	if err != nil {
		logger.Error(err, "failed to generate JWT token")
		return nil, errors.New("unable to complete login")
	}

	return &AuthResponse{
		Token: token,
		User:  u.ToResponse(),
	}, nil
}
