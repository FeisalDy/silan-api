package casbin

import (
	"log"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

// NewEnforcer creates a new Casbin enforcer with GORM adapter
func NewEnforcer(db *gorm.DB, modelPath string) (*casbin.Enforcer, error) {
	// Initialize GORM adapter - this will create casbin_rule table in Postgres
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, err
	}

	// Create enforcer
	enforcer, err := casbin.NewEnforcer(modelPath, adapter)
	if err != nil {
		return nil, err
	}

	// Load policies from database
	if err := enforcer.LoadPolicy(); err != nil {
		return nil, err
	}

	log.Println("Casbin enforcer initialized successfully")
	return enforcer, nil
}

// InitializeDefaultPolicies sets up default RBAC policies
func InitializeDefaultPolicies(enforcer *casbin.Enforcer) error {
	// Define default policies
	policies := [][]string{
		// Admin can do everything
		{"admin", "/api/v1/users", "GET"},
		{"admin", "/api/v1/users/*", "GET"},
		{"admin", "/api/v1/users/*", "PUT"},
		{"admin", "/api/v1/users/*", "DELETE"},

		// Admin can manage all novels and chapters
		{"admin", "/api/v1/novels", "POST"},
		{"admin", "/api/v1/novels/*", "PUT"},
		{"admin", "/api/v1/novels/*", "DELETE"},
		{"admin", "/api/v1/novels/translations", "POST"},
		{"admin", "/api/v1/novels/translations/*", "PUT"},
		{"admin", "/api/v1/novels/translations/*", "DELETE"},
		{"admin", "/api/v1/chapters", "POST"},
		{"admin", "/api/v1/chapters/*", "PUT"},
		{"admin", "/api/v1/chapters/*", "DELETE"},
		{"admin", "/api/v1/chapters/translations", "POST"},
		{"admin", "/api/v1/chapters/translations/*", "PUT"},
		{"admin", "/api/v1/chapters/translations/*", "DELETE"},

		// Users can view their own profile
		{"user", "/api/v1/users/*", "GET"},

		// Users can create novels and chapters
		{"user", "/api/v1/novels", "POST"},
		{"user", "/api/v1/novels/translations", "POST"},
		{"user", "/api/v1/chapters", "POST"},
		{"user", "/api/v1/chapters/translations", "POST"},

		// Author role - can manage their own novels and chapters
		{"author", "/api/v1/novels", "POST"},
		{"author", "/api/v1/novels/*", "PUT"},
		{"author", "/api/v1/novels/*", "DELETE"},
		{"author", "/api/v1/novels/translations", "POST"},
		{"author", "/api/v1/novels/translations/*", "PUT"},
		{"author", "/api/v1/novels/translations/*", "DELETE"},
		{"author", "/api/v1/chapters", "POST"},
		{"author", "/api/v1/chapters/*", "PUT"},
		{"author", "/api/v1/chapters/*", "DELETE"},
		{"author", "/api/v1/chapters/translations", "POST"},
		{"author", "/api/v1/chapters/translations/*", "PUT"},
		{"author", "/api/v1/chapters/translations/*", "DELETE"},

		// Translator role - can create and manage translations
		{"translator", "/api/v1/novels/translations", "POST"},
		{"translator", "/api/v1/novels/translations/*", "PUT"},
		{"translator", "/api/v1/novels/translations/*", "DELETE"},
		{"translator", "/api/v1/chapters/translations", "POST"},
		{"translator", "/api/v1/chapters/translations/*", "PUT"},
		{"translator", "/api/v1/chapters/translations/*", "DELETE"},
	}

	// Add policies
	for _, policy := range policies {
		_, err := enforcer.AddPolicy(policy)
		if err != nil {
			log.Printf("Warning: Failed to add policy %v: %v", policy, err)
			// Continue even if policy already exists
		}
	}

	// Save policies to database
	if err := enforcer.SavePolicy(); err != nil {
		return err
	}

	log.Println("Default Casbin policies initialized")
	return nil
}
