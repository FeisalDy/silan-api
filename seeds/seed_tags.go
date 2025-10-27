package seeds

import (
	"log"
	"simple-go/internal/domain/tag"

	"gorm.io/gorm"
)

// SeedTags seeds default tags
func SeedTags(db *gorm.DB) error {
	log.Println("🌱 Seeding tags...")

	tags := []tag.Tag{
		{Name: "Reincarnation", Slug: "reincarnation", Description: strPtr("Characters reborn into new lives")},
		{Name: "Overpowered MC", Slug: "overpowered-mc", Description: strPtr("Main character with exceptional power")},
		{Name: "System", Slug: "system", Description: strPtr("Game-like system mechanics")},
		{Name: "Magic", Slug: "magic", Description: strPtr("Supernatural powers and spells")},
		{Name: "Cultivation", Slug: "cultivation", Description: strPtr("Eastern martial arts progression")},
		{Name: "Transmigration", Slug: "transmigration", Description: strPtr("Soul transfer to another body or world")},
		{Name: "Harem", Slug: "harem", Description: strPtr("Multiple romantic interests")},
		{Name: "Revenge", Slug: "revenge", Description: strPtr("Quest for vengeance")},
		{Name: "Weak to Strong", Slug: "weak-to-strong", Description: strPtr("Character progression from weakness to strength")},
		{Name: "Academy", Slug: "academy", Description: strPtr("School or academy setting")},
		{Name: "Dungeon", Slug: "dungeon", Description: strPtr("Underground labyrinth exploration")},
		{Name: "Isekai", Slug: "isekai", Description: strPtr("Transported to another world")},
		{Name: "Virtual Reality", Slug: "virtual-reality", Description: strPtr("Game or VR world setting")},
		{Name: "Monster", Slug: "monster", Description: strPtr("Creatures and beasts")},
		{Name: "Adventure", Slug: "adventure", Description: strPtr("Journey and exploration")},
		{Name: "Kingdom Building", Slug: "kingdom-building", Description: strPtr("Creating and managing territories")},
		{Name: "Time Travel", Slug: "time-travel", Description: strPtr("Moving through different time periods")},
		{Name: "Anti-Hero", Slug: "anti-hero", Description: strPtr("Protagonist with questionable morals")},
		{Name: "Female Lead", Slug: "female-lead", Description: strPtr("Female protagonist")},
		{Name: "Male Lead", Slug: "male-lead", Description: strPtr("Male protagonist")},
	}

	for _, t := range tags {
		var existing tag.Tag
		result := db.Where("slug = ?", t.Slug).First(&existing)

		if result.Error == gorm.ErrRecordNotFound {
			if err := db.Create(&t).Error; err != nil {
				log.Printf("⚠️  Failed to seed tag %s: %v", t.Name, err)
				return err
			}
			log.Printf("✅ Created tag: %s", t.Name)
		} else if result.Error != nil {
			return result.Error
		} else {
			log.Printf("⏭️  Tag already exists: %s", t.Name)
		}
	}

	log.Println("✅ Tags seeding completed")
	return nil
}
