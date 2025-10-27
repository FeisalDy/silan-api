package casbin

import (
	"fmt"

	"github.com/casbin/casbin/v2"
)

type Resource string

const (
	ResourceUser               Resource = "user"
	ResourceNovel              Resource = "novel"
	ResourceNovelTranslation   Resource = "novel_translation"
	ResourceChapter            Resource = "chapter"
	ResourceChapterTranslation Resource = "chapter_translation"
)

type Action string

const (
	ActionList   Action = "list"
	ActionRead   Action = "read"
	ActionCreate Action = "create"
	ActionUpdate Action = "update"
	ActionDelete Action = "delete"
)

// Role represents a user role
type Role string

const (
	RoleAdmin      Role = "admin"
	RoleAuthor     Role = "author"
	RoleTranslator Role = "translator"
	RoleUser       Role = "user"
)

// PolicyHelper provides utility functions for working with Casbin policies
type PolicyHelper struct {
	enforcer *casbin.Enforcer
}

// NewPolicyHelper creates a new PolicyHelper
func NewPolicyHelper(enforcer *casbin.Enforcer) *PolicyHelper {
	return &PolicyHelper{enforcer: enforcer}
}

// GrantPermission grants a role permission to perform an action on a resource
func (h *PolicyHelper) GrantPermission(role Role, resource Resource, action Action) error {
	_, err := h.enforcer.AddPolicy(string(role), string(resource), string(action))
	if err != nil {
		return fmt.Errorf("failed to grant permission: %w", err)
	}
	return h.enforcer.SavePolicy()
}

// RevokePermission revokes a role's permission to perform an action on a resource
func (h *PolicyHelper) RevokePermission(role Role, resource Resource, action Action) error {
	_, err := h.enforcer.RemovePolicy(string(role), string(resource), string(action))
	if err != nil {
		return fmt.Errorf("failed to revoke permission: %w", err)
	}
	return h.enforcer.SavePolicy()
}

// HasPermission checks if a role has permission to perform an action on a resource
func (h *PolicyHelper) HasPermission(role Role, resource Resource, action Action) (bool, error) {
	return h.enforcer.Enforce(string(role), string(resource), string(action))
}

// GetRolePermissions returns all permissions for a given role
func (h *PolicyHelper) GetRolePermissions(role Role) []Permission {
	policies, _ := h.enforcer.GetFilteredPolicy(0, string(role))
	permissions := make([]Permission, 0, len(policies))

	for _, policy := range policies {
		if len(policy) >= 3 {
			permissions = append(permissions, Permission{
				Role:     Role(policy[0]),
				Resource: Resource(policy[1]),
				Action:   Action(policy[2]),
			})
		}
	}

	return permissions
}

// GetResourcePermissions returns all permissions for a given resource
func (h *PolicyHelper) GetResourcePermissions(resource Resource) []Permission {
	policies, _ := h.enforcer.GetFilteredPolicy(1, string(resource))
	permissions := make([]Permission, 0, len(policies))

	for _, policy := range policies {
		if len(policy) >= 3 {
			permissions = append(permissions, Permission{
				Role:     Role(policy[0]),
				Resource: Resource(policy[1]),
				Action:   Action(policy[2]),
			})
		}
	}

	return permissions
}

// GetAllPermissions returns all permissions in the system
func (h *PolicyHelper) GetAllPermissions() []Permission {
	policies, _ := h.enforcer.GetPolicy()
	permissions := make([]Permission, 0, len(policies))

	for _, policy := range policies {
		if len(policy) >= 3 {
			permissions = append(permissions, Permission{
				Role:     Role(policy[0]),
				Resource: Resource(policy[1]),
				Action:   Action(policy[2]),
			})
		}
	}

	return permissions
}

// ClearRolePermissions removes all permissions for a role
func (h *PolicyHelper) ClearRolePermissions(role Role) error {
	_, err := h.enforcer.RemoveFilteredPolicy(0, string(role))
	if err != nil {
		return fmt.Errorf("failed to clear role permissions: %w", err)
	}
	return h.enforcer.SavePolicy()
}

// ClearResourcePermissions removes all permissions for a resource
func (h *PolicyHelper) ClearResourcePermissions(resource Resource) error {
	_, err := h.enforcer.RemoveFilteredPolicy(1, string(resource))
	if err != nil {
		return fmt.Errorf("failed to clear resource permissions: %w", err)
	}
	return h.enforcer.SavePolicy()
}

// Permission represents a role-resource-action triple
type Permission struct {
	Role     Role
	Resource Resource
	Action   Action
}

// String returns a string representation of the permission
func (p Permission) String() string {
	return fmt.Sprintf("%s can %s %s", p.Role, p.Action, p.Resource)
}

// ValidatePermission validates that a permission uses valid constants
func ValidatePermission(role Role, resource Resource, action Action) error {
	// Validate role
	validRoles := []Role{RoleAdmin, RoleAuthor, RoleTranslator, RoleUser}
	roleValid := false
	for _, r := range validRoles {
		if role == r {
			roleValid = true
			break
		}
	}
	if !roleValid {
		return fmt.Errorf("invalid role: %s", role)
	}

	// Validate resource
	validResources := []Resource{
		ResourceUser,
		ResourceNovel,
		ResourceNovelTranslation,
		ResourceChapter,
		ResourceChapterTranslation,
	}
	resourceValid := false
	for _, res := range validResources {
		if resource == res {
			resourceValid = true
			break
		}
	}
	if !resourceValid {
		return fmt.Errorf("invalid resource: %s", resource)
	}

	// Validate action
	validActions := []Action{ActionList, ActionRead, ActionCreate, ActionUpdate, ActionDelete}
	actionValid := false
	for _, act := range validActions {
		if action == act {
			actionValid = true
			break
		}
	}
	if !actionValid {
		return fmt.Errorf("invalid action: %s", action)
	}

	return nil
}

// GrantFullAccess grants all permissions for a resource to a role
func (h *PolicyHelper) GrantFullAccess(role Role, resource Resource) error {
	actions := []Action{ActionList, ActionRead, ActionCreate, ActionUpdate, ActionDelete}

	for _, action := range actions {
		if err := h.GrantPermission(role, resource, action); err != nil {
			return err
		}
	}

	return nil
}

// GrantReadOnlyAccess grants read permissions for a resource to a role
func (h *PolicyHelper) GrantReadOnlyAccess(role Role, resource Resource) error {
	if err := h.GrantPermission(role, resource, ActionList); err != nil {
		return err
	}
	return h.GrantPermission(role, resource, ActionRead)
}

// GrantWriteAccess grants create and update permissions for a resource to a role
func (h *PolicyHelper) GrantWriteAccess(role Role, resource Resource) error {
	if err := h.GrantPermission(role, resource, ActionCreate); err != nil {
		return err
	}
	return h.GrantPermission(role, resource, ActionUpdate)
}
