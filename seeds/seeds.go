package seeds

import (
	"log"

	"gorm.io/gorm"
)

// Helper function to create string pointers
func strPtr(s string) *string {
	return &s
}

// Helper function to create int pointers
func intPtr(i int) *int {
	return &i
}

// RunAllSeeds runs all database seeders in the correct order
// Seeds are only created if they don't already exist
func RunAllSeeds(db *gorm.DB) error {
	log.Println("🌱 Starting database seeding...")

	// Order matters! Seed in dependency order
	seeders := []struct {
		name string
		fn   func(*gorm.DB) error
	}{
		{"Roles", SeedRoles},
		{"Users", SeedUsers},
		{"Genres", SeedGenres},
		{"Tags", SeedTags},
		{"Media", SeedMedia},
		{"Novels", SeedNovels},
		{"Volumes", SeedVolumes},
		{"Chapters", SeedChapters},
	}

	for _, seeder := range seeders {
		log.Printf("Running seeder: %s", seeder.name)
		if err := seeder.fn(db); err != nil {
			log.Printf("❌ Failed to run %s seeder: %v", seeder.name, err)
			return err
		}
	}

	log.Println("✅ All seeders completed successfully!")
	return nil
}
