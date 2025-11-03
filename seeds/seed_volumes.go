package seeds

import (
	"log"
	"simple-go/internal/domain/novel"
	"simple-go/internal/domain/user"
	"simple-go/internal/domain/volume"
	"strconv"

	"gorm.io/gorm"
)

func SeedVolumes(db *gorm.DB) error {
	log.Println("Seeding Volumes...")

	var translator user.User
	if err := db.Where("email = ?", "translator1@example.com").First(&translator).Error; err != nil {
		// Fallback to author if translator not found
		if err := db.Where("email = ?", "author1@example.com").First(&translator).Error; err != nil {
			log.Println("⚠️  No suitable user found for chapter seeding")
			return nil
		}
	}

	var novels []novel.Novel
	if err := db.Find(&novels).Error; err != nil {
		log.Println("⚠️  No novels found, skipping volumes seeding")
		return nil
	}

	for i, n := range novels {
		isLastNovel := i == len(novels)-1

		if isLastNovel {
			volumeNum := 1
			var existingVolume volume.Volume
			result := db.Where("novel_id = ? AND number = ?", n.ID, volumeNum).First(&existingVolume)

			if result.Error == gorm.ErrRecordNotFound {
				v := volume.Volume{
					NovelID:          n.ID,
					Number:           volumeNum,
					CoverMediaID:     nil,
					OriginalLanguage: n.OriginalLanguage,
					IsVirtual:        true,
				}

				if err := db.Create(&v).Error; err != nil {
					log.Printf("⚠️  Failed to seed virtual volume for novel %s: %v", n.ID, err)
					continue
				}

				translations := []volume.VolumeTranslation{
					{
						VolumeID:    v.ID,
						Lang:        n.OriginalLanguage,
						Title:       getVirtualVolumeTitle(n.OriginalLanguage),
						Description: nil,
					},
					{
						VolumeID:    v.ID,
						Lang:        "en",
						Title:       "Virtual Volume: Adventures Ahead",
						Description: nil,
					},
					{
						VolumeID:    v.ID,
						Lang:        "id",
						Title:       "Volume Virtual: Petualangan di Depan",
						Description: nil,
					},
				}

				for _, trans := range translations {
					if err := db.Create(&trans).Error; err != nil {
						log.Printf("⚠️  Failed to seed virtual volume translation: %v", err)
					}
				}

				log.Printf("✅ Created virtual volume for last novel ID %s", n.ID)
			} else if result.Error != nil {
				return result.Error
			} else {
				log.Printf("⏭️  Virtual volume for last novel %s already exists", n.ID)
			}
		} else {
			for volumeNum := 1; volumeNum <= 3; volumeNum++ {
				var existingVolume volume.Volume
				result := db.Where("novel_id = ? AND number = ?", n.ID, volumeNum).First(&existingVolume)

				if result.Error == gorm.ErrRecordNotFound {
					v := volume.Volume{
						NovelID:          n.ID,
						Number:           volumeNum,
						CoverMediaID:     nil,
						OriginalLanguage: n.OriginalLanguage,
						IsVirtual:        false,
					}

					if err := db.Create(&v).Error; err != nil {
						log.Printf("⚠️  Failed to seed volume %d for novel %s: %v", volumeNum, n.ID, err)
						continue
					}

					translations := []volume.VolumeTranslation{
						{
							VolumeID:    v.ID,
							Lang:        n.OriginalLanguage,
							Title:       getVolumeTitle(n.OriginalLanguage, volumeNum),
							Description: nil,
						},
						{
							VolumeID:    v.ID,
							Lang:        "en",
							Title:       "Volume " + strconv.Itoa(int(volumeNum)) + ": Adventures Ahead",
							Description: nil,
						},
						{
							VolumeID:    v.ID,
							Lang:        "id",
							Title:       "Volume " + strconv.Itoa(int(volumeNum)) + ": Petualangan di Depan",
							Description: nil,
						},
					}

					for _, trans := range translations {
						if err := db.Create(&trans).Error; err != nil {
							log.Printf("⚠️  Failed to seed volume translation: %v", err)
						}
					}

					log.Printf("✅ Created volume %d for novel ID %s", volumeNum, n.ID)
				} else if result.Error != nil {
					return result.Error
				} else {
					log.Printf("⏭️  Volume %d for novel %s already exists", volumeNum, n.ID)
				}
			}
		}
	}

	log.Println("✅ Volumes seeding completed")
	return nil
}

// Helper function to get volume title in different languages
func getVolumeTitle(lang string, volumeNum int) string {
	volumeNumStr := strconv.Itoa(volumeNum)
	switch lang {
	case "ko":
		return "제 " + volumeNumStr + " 권: 앞으로의 모험"
	case "ja":
		return "第" + volumeNumStr + "巻: 冒険の始まり"
	case "zh-CN":
		return "第" + volumeNumStr + "卷: 即将到来的冒险"
	case "id":
		return "Volume " + volumeNumStr + ": Petualangan di Depan"
	default: // en
		return "Volume " + volumeNumStr + ": Adventures Ahead"
	}
}

// Helper function to get virtual volume title in different languages
func getVirtualVolumeTitle(lang string) string {
	switch lang {
	case "ko":
		return "가상 권: 앞으로의 모험"
	case "ja":
		return "仮想巻: 冒険の始まり"
	case "zh-CN":
		return "虚拟卷: 即将到来的冒险"
	case "id":
		return "Volume Virtual: Petualangan di Depan"
	default: // en
		return "Virtual Volume: Adventures Ahead"
	}
}
