package seeds

import (
	"log"
	"simple-go/internal/domain/role"

	"gorm.io/gorm"
)

// SeedRoles seeds default roles
func SeedRoles(db *gorm.DB) error {
	log.Println("🌱 Seeding roles...")

	roles := []role.Role{
		{Name: "admin", Description: strPtr("Administrator with full access to all resources")},
		{Name: "author", Description: strPtr("Content creator who can write and manage novels and chapters")},
		{Name: "translator", Description: strPtr("Translator who can create and manage translations")},
		{Name: "user", Description: strPtr("Regular user with read access")},
	}

	for _, r := range roles {
		var existing role.Role
		result := db.Where("name = ?", r.Name).First(&existing)

		if result.Error == gorm.ErrRecordNotFound {
			if err := db.Create(&r).Error; err != nil {
				log.Printf("⚠️  Failed to seed role %s: %v", r.Name, err)
				return err
			}
			log.Printf("✅ Created role: %s", r.Name)
		} else if result.Error != nil {
			return result.Error
		} else {
			log.Printf("⏭️  Role already exists: %s", r.Name)
		}
	}

	log.Println("✅ Roles seeding completed")
	return nil
}
