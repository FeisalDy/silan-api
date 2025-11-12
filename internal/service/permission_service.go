package service

import (
	"context"
	"simple-go/internal/repository"

	"github.com/casbin/casbin/v2"
)

// PermissionService handles permission-related business logic
type PermissionService struct {
	userRepo repository.UserRepository
	enforcer *casbin.Enforcer
}

// NewPermissionService creates a new instance of PermissionService
func NewPermissionService(
	userRepo repository.UserRepository,
	enforcer *casbin.Enforcer,
) *PermissionService {
	return &PermissionService{
		userRepo: userRepo,
		enforcer: enforcer,
	}
}

// PermissionMap represents permissions grouped by resource
type PermissionMap map[string][]string

// GetUserPermissions retrieves all permissions for a user based on their roles
func (s *PermissionService) GetUserPermissions(ctx context.Context, userID string) (PermissionMap, error) {
	// Get user roles
	roles, err := s.userRepo.GetRoles(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Initialize permission map
	permissions := make(PermissionMap)

	// Get permissions for each role from Casbin
	for _, role := range roles {
		// Get all policies for this role
		// Format: [role, resource, action]
		policies, _ := s.enforcer.GetFilteredPolicy(0, role.Name)

		for _, policy := range policies {
			if len(policy) >= 3 {
				resource := policy[1]
				action := policy[2]

				// Initialize resource array if it doesn't exist
				if _, exists := permissions[resource]; !exists {
					permissions[resource] = []string{}
				}

				// Add action if not already present
				actionExists := false
				for _, existingAction := range permissions[resource] {
					if existingAction == action {
						actionExists = true
						break
					}
				}

				if !actionExists {
					permissions[resource] = append(permissions[resource], action)
				}
			}
		}
	}

	return permissions, nil
}
