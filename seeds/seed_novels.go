package seeds

import (
	"log"
	"simple-go/internal/domain/genre"
	"simple-go/internal/domain/novel"
	"simple-go/internal/domain/tag"
	"simple-go/internal/domain/user"

	"gorm.io/gorm"
)

// SeedNovels seeds sample novels with translations
func SeedNovels(db *gorm.DB) error {
	log.Println("🌱 Seeding novels...")

	// Get author user
	var author user.User
	if err := db.Where("email = ?", "author1@example.com").First(&author).Error; err != nil {
		log.Println("⚠️  Author user not found, skipping novel seeding")
		return nil
	}

	// Get translator user
	var translator user.User
	if err := db.Where("email = ?", "translator1@example.com").First(&translator).Error; err != nil {
		log.Println("⚠️  Translator user not found, using author as translator")
		translator = author
	}

	novels := []struct {
		Novel        novel.Novel
		Translations []novel.NovelTranslation
		GenreSlugs   []string
		TagSlugs     []string
	}{
		{
			Novel: novel.Novel{
				OriginalLanguage: "zh-CN",
				OriginalAuthor:   strPtr("李明"),
				Status:           strPtr("ongoing"),
				Source:           strPtr("https://example.com/novel1"),
				WordCount:        intPtr(500000),
				CreatedBy:        author.ID,
			},
			Translations: []novel.NovelTranslation{
				{
					Lang:        "en",
					Title:       "The Legendary Cultivator",
					Description: strPtr("A young man embarks on a journey to become the strongest cultivator in the realm. Facing countless trials and powerful enemies, he must master ancient techniques and forge his own path to immortality."),
				},
				{
					Lang:        "id",
					Title:       "Kultivator Legendaris",
					Description: strPtr("Seorang pemuda memulai perjalanan untuk menjadi kultivator terkuat di alam semesta. Menghadapi berbagai ujian dan musuh yang kuat, ia harus menguasai teknik kuno dan membentuk jalannya sendiri menuju keabadian."),
				},
			},
			GenreSlugs: []string{"fantasy", "action", "martial-arts"},
			TagSlugs:   []string{"cultivation", "weak-to-strong", "overpowered-mc", "magic"},
		},
		{
			Novel: novel.Novel{
				OriginalLanguage: "ja",
				OriginalAuthor:   strPtr("田中太郎"),
				Status:           strPtr("ongoing"),
				Source:           strPtr("https://example.com/novel2"),
				WordCount:        intPtr(300000),
				CreatedBy:        author.ID,
			},
			Translations: []novel.NovelTranslation{
				{
					Lang:        "en",
					Title:       "Reborn in Another World",
					Description: strPtr("After dying in a tragic accident, a salary worker finds himself reborn in a fantasy world with game-like mechanics. Armed with knowledge from his previous life, he sets out to live freely and enjoy his second chance."),
				},
				{
					Lang:        "id",
					Title:       "Terlahir Kembali di Dunia Lain",
					Description: strPtr("Setelah meninggal dalam kecelakaan tragis, seorang pekerja kantoran mendapati dirinya terlahir kembali di dunia fantasi dengan mekanik seperti game. Dipersenjatai dengan pengetahuan dari kehidupan sebelumnya, ia bertekad untuk hidup bebas dan menikmati kesempatan keduanya."),
				},
			},
			GenreSlugs: []string{"fantasy", "comedy", "slice-of-life"},
			TagSlugs:   []string{"isekai", "reincarnation", "system", "adventure"},
		},
		{
			Novel: novel.Novel{
				OriginalLanguage: "ko",
				OriginalAuthor:   strPtr("김철수"),
				Status:           strPtr("completed"),
				Source:           strPtr("https://example.com/novel3"),
				WordCount:        intPtr(800000),
				CreatedBy:        author.ID,
			},
			Translations: []novel.NovelTranslation{
				{
					Lang:        "en",
					Title:       "Shadow Monarch",
					Description: strPtr("In a world where dungeons and monsters have become reality, the weakest hunter receives a mysterious power that allows him to rise from the dead and command an army of shadows. His journey from the weakest to the strongest begins."),
				},
				{
					Lang:        "id",
					Title:       "Raja Bayangan",
					Description: strPtr("Di dunia di mana dungeon dan monster telah menjadi kenyataan, pemburu terlemah menerima kekuatan misterius yang memungkinkannya bangkit dari kematian dan mengendalikan pasukan bayangan. Perjalanannya dari yang terlemah ke yang terkuat dimulai."),
				},
			},
			GenreSlugs: []string{"action", "fantasy", "horror"},
			TagSlugs:   []string{"weak-to-strong", "system", "dungeon", "monster", "overpowered-mc"},
		},
	}

	for i, novelData := range novels {
		// Check if novel exists (by author and language combination)
		var existingNovel novel.Novel
		result := db.Where("original_language = ? AND original_author = ?",
			novelData.Novel.OriginalLanguage,
			*novelData.Novel.OriginalAuthor).First(&existingNovel)

		if result.Error == gorm.ErrRecordNotFound {
			// Create novel
			if err := db.Create(&novelData.Novel).Error; err != nil {
				log.Printf("⚠️  Failed to seed novel %d: %v", i+1, err)
				return err
			}

			// Create translations
			for _, trans := range novelData.Translations {
				trans.NovelID = novelData.Novel.ID
				if err := db.Create(&trans).Error; err != nil {
					log.Printf("⚠️  Failed to seed translation for novel %d: %v", i+1, err)
				}
			}

			// Add genres using GORM Association API (Novel owns the M2M relationship)
			var genres []genre.Genre
			for _, genreSlug := range novelData.GenreSlugs {
				var g genre.Genre
				if err := db.Where("slug = ?", genreSlug).First(&g).Error; err == nil {
					genres = append(genres, g)
				}
			}
			if len(genres) > 0 {
				if err := db.Model(&novelData.Novel).Association("Genres").Append(genres); err != nil {
					log.Printf("⚠️  Failed to assign genres to novel %d: %v", i+1, err)
				}
			}

			// Add tags using GORM Association API (Novel owns the M2M relationship)
			var tags []tag.Tag
			for _, tagSlug := range novelData.TagSlugs {
				var t tag.Tag
				if err := db.Where("slug = ?", tagSlug).First(&t).Error; err == nil {
					tags = append(tags, t)
				}
			}
			if len(tags) > 0 {
				if err := db.Model(&novelData.Novel).Association("Tags").Append(tags); err != nil {
					log.Printf("⚠️  Failed to assign tags to novel %d: %v", i+1, err)
				}
			}

			log.Printf("✅ Created novel: %s", novelData.Translations[0].Title)
		} else if result.Error != nil {
			return result.Error
		} else {
			log.Printf("⏭️  Novel already exists: %s", *novelData.Novel.OriginalAuthor)
		}
	}

	log.Println("✅ Novels seeding completed")
	return nil
}
