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

// InitializeDefaultPolicies sets up default RBAC policies using resource-action approach
func InitializeDefaultPolicies(enforcer *casbin.Enforcer) error {
	// Define default policies using resource:action format
	// Format: {role, resource, action}
	policies := [][]string{
		// ============ ADMIN ROLE ============
		// Admin can do everything
		{"admin", "user", "list"},
		{"admin", "user", "read"},
		{"admin", "user", "update"},
		{"admin", "user", "delete"},

		{"admin", "novel", "create"},
		{"admin", "novel", "read"},
		{"admin", "novel", "update"},
		{"admin", "novel", "delete"},

		{"admin", "novel_translation", "create"},
		{"admin", "novel_translation", "read"},
		{"admin", "novel_translation", "update"},
		{"admin", "novel_translation", "delete"},

		{"admin", "chapter", "create"},
		{"admin", "chapter", "read"},
		{"admin", "chapter", "update"},
		{"admin", "chapter", "delete"},

		{"admin", "chapter_translation", "create"},
		{"admin", "chapter_translation", "read"},
		{"admin", "chapter_translation", "update"},
		{"admin", "chapter_translation", "delete"},

		{"admin", "translation_job", "create"},
		{"admin", "translation_job", "read"},
		{"admin", "translation_job", "list"},
		{"admin", "translation_job", "update"},
		{"admin", "translation_job", "delete"},

		// ============ USER ROLE ============
		// Basic user can read their own profile
		{"user", "user", "read"},

		// Basic user can create novels and chapters
		{"user", "novel", "read"},
		{"user", "chapter", "read"},
		{"user", "novel", "list"},
		{"user", "chapter", "list"},

		// ============ AUTHOR ROLE ============
		// Author can manage novels and chapters
		{"author", "novel", "create"},
		{"author", "novel", "read"},
		{"author", "novel", "update"},
		{"author", "novel", "delete"},

		{"author", "chapter", "create"},
		{"author", "chapter", "read"},
		{"author", "chapter", "update"},
		{"author", "chapter", "delete"},

		{"author", "novel_translation", "create"},
		{"author", "novel_translation", "update"},
		{"author", "novel_translation", "delete"},

		{"author", "chapter_translation", "create"},
		{"author", "chapter_translation", "update"},
		{"author", "chapter_translation", "delete"},

		// ============ TRANSLATOR ROLE ============
		// Translator can manage translations only
		{"translator", "novel_translation", "create"},
		{"translator", "novel_translation", "read"},
		{"translator", "novel_translation", "update"},
		{"translator", "novel_translation", "delete"},

		{"translator", "chapter_translation", "create"},
		{"translator", "chapter_translation", "read"},
		{"translator", "chapter_translation", "update"},
		{"translator", "chapter_translation", "delete"},

		// Translators need to read novels and chapters to translate them
		{"translator", "novel", "read"},
		{"translator", "chapter", "read"},
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

	log.Println("Default Casbin policies initialized with resource-action approach")
	return nil
}
