package seeds

import (
	"log"
	"simple-go/internal/domain/genre"

	"gorm.io/gorm"
)

// SeedGenres seeds default genres
func SeedGenres(db *gorm.DB) error {
	log.Println("🌱 Seeding genres...")

	genres := []genre.Genre{
		{
			Name:        "Fantasy",
			Slug:        "fantasy",
			Description: strPtr("Stories featuring magical elements, mythical creatures, and imaginary worlds"),
		},
		{
			Name:        "Romance",
			Slug:        "romance",
			Description: strPtr("Stories centered on romantic relationships and emotional connections"),
		},
		{
			Name:        "Action",
			Slug:        "action",
			Description: strPtr("Fast-paced stories with physical conflict, battles, and adventures"),
		},
		{
			Name:        "Mystery",
			Slug:        "mystery",
			Description: strPtr("Stories involving puzzles, crimes, and investigative plots"),
		},
		{
			Name:        "Science Fiction",
			Slug:        "science-fiction",
			Description: strPtr("Stories exploring futuristic concepts, technology, and space"),
		},
		{
			Name:        "Horror",
			Slug:        "horror",
			Description: strPtr("Stories designed to frighten and unsettle readers"),
		},
		{
			Name:        "Comedy",
			Slug:        "comedy",
			Description: strPtr("Humorous stories meant to entertain and amuse"),
		},
		{
			Name:        "Drama",
			Slug:        "drama",
			Description: strPtr("Stories focused on realistic characters and emotional themes"),
		},
		{
			Name:        "Slice of Life",
			Slug:        "slice-of-life",
			Description: strPtr("Stories depicting everyday life and mundane activities"),
		},
		{
			Name:        "Historical",
			Slug:        "historical",
			Description: strPtr("Stories set in the past with historical contexts"),
		},
		{
			Name:        "Martial Arts",
			Slug:        "martial-arts",
			Description: strPtr("Stories featuring martial arts, cultivation, and combat techniques"),
		},
		{
			Name:        "Psychological",
			Slug:        "psychological",
			Description: strPtr("Stories exploring the human mind, emotions, and behavior"),
		},
	}

	for _, g := range genres {
		var existing genre.Genre
		result := db.Where("slug = ?", g.Slug).First(&existing)

		if result.Error == gorm.ErrRecordNotFound {
			if err := db.Create(&g).Error; err != nil {
				log.Printf("⚠️  Failed to seed genre %s: %v", g.Name, err)
				return err
			}
			log.Printf("✅ Created genre: %s", g.Name)
		} else if result.Error != nil {
			return result.Error
		} else {
			log.Printf("⏭️  Genre already exists: %s", g.Name)
		}
	}

	log.Println("✅ Genres seeding completed")
	return nil
}
