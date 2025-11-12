package seeds

import (
	"log"
	"simple-go/internal/domain/media"
	"simple-go/internal/domain/novel"

	"gorm.io/gorm"
)

// SeedMedia ensures a default cover media exists and assigns it to novels without a cover
func SeedMedia(db *gorm.DB) error {
	log.Println("🌱 Seeding media...")

	// Default media data
	url := "https://picsum.photos/200/300"
	mediaType := "image"

	var m media.Media
	// Check if media with the URL already exists (idempotent)
	result := db.Where("url = ?", url).First(&m)

	if result.Error == gorm.ErrRecordNotFound {
		m = media.Media{
			URL:  &url,
			Type: &mediaType,
		}
		if err := db.Create(&m).Error; err != nil {
			log.Printf("❌ Failed to create default media: %v", err)
			return err
		}
		log.Printf("✅ Created default media with ID: %s", m.ID)
	} else if result.Error != nil {
		return result.Error
	} else {
		log.Printf("⏭️  Default media already exists with ID: %s", m.ID)
	}

	// Assign this media to all novels missing a cover
	if err := db.Model(&novel.Novel{}).Where("cover_media_id IS NULL").Update("cover_media_id", m.ID).Error; err != nil {
		log.Printf("⚠️  Failed to assign media to novels: %v", err)
		// Don't fail the whole seeding because of assignment; return nil to proceed
		return nil
	}
	log.Println("✅ Assigned default media to novels without a cover")

	log.Println("✅ Media seeding completed")
	return nil
}
