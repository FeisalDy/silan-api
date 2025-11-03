package seeds

import (
	"fmt"
	"log"
	"simple-go/internal/domain/chapter"
	"simple-go/internal/domain/user"
	"simple-go/internal/domain/volume"

	"gorm.io/gorm"
)

// SeedChapters seeds sample chapters for existing novels
func SeedChapters(db *gorm.DB) error {
	log.Println("🌱 Seeding chapters...")

	// Get translator user
	var translator user.User
	if err := db.Where("email = ?", "translator1@example.com").First(&translator).Error; err != nil {
		// Fallback to author if translator not found
		if err := db.Where("email = ?", "author1@example.com").First(&translator).Error; err != nil {
			log.Println("⚠️  No suitable user found for chapter seeding")
			return nil
		}
	}

	// Get all novels
	var volumes []volume.Volume
	if err := db.Find(&volumes).Error; err != nil {
		log.Println("⚠️  No volumes found, skipping chapter seeding")
		return nil
	}

	// Seed 5 chapters for each volume
	for _, v := range volumes {
		for chapterNum := 1; chapterNum <= 5; chapterNum++ {
			// Check if chapter exists
			var existingChapter chapter.Chapter
			result := db.Where("volume_id = ? AND number = ?", v.ID, chapterNum).First(&existingChapter)

			if result.Error == gorm.ErrRecordNotFound {
				ch := chapter.Chapter{
					VolumeID:  v.ID,
					Number:    chapterNum,
					WordCount: intPtr(2500 + (chapterNum * 100)),
				}

				if err := db.Create(&ch).Error; err != nil {
					log.Printf("⚠️  Failed to seed chapter %d for volume %s: %v", chapterNum, v.ID, err)
					continue
				}

				// Create translations for this chapter
				translations := []chapter.ChapterTranslation{
					{
						ChapterID: ch.ID,
						Lang:      v.OriginalLanguage,
						Title:     getChapterTitle(v.OriginalLanguage, chapterNum),
						Content:   generateChapterContent(v.OriginalLanguage, chapterNum),
					},
					{
						ChapterID: ch.ID,
						Lang:      "en",
						Title:     fmt.Sprintf("Chapter %d: The Beginning of the Journey", chapterNum),
						Content:   generateChapterContent("en", chapterNum),
					},
					{
						ChapterID: ch.ID,
						Lang:      "id",
						Title:     fmt.Sprintf("Bab %d: Permulaan Perjalanan", chapterNum),
						Content:   generateChapterContent("id", chapterNum),
					},
				}

				for _, trans := range translations {
					if err := db.Create(&trans).Error; err != nil {
						log.Printf("⚠️  Failed to seed chapter translation: %v", err)
					}
				}

				log.Printf("✅ Created chapter %d for volume ID %s", chapterNum, v.ID)
			} else if result.Error != nil {
				return result.Error
			} else {
				log.Printf("⏭️  Chapter %d for volume %s already exists", chapterNum, v.ID)
			}
		}
	}

	log.Println("✅ Chapters seeding completed")
	return nil
}

// Helper function to get chapter title in different languages
func getChapterTitle(lang string, chapterNum int) string {
	switch lang {
	case "ko":
		return fmt.Sprintf("제%d화: 여정의 시작", chapterNum)
	case "ja":
		return fmt.Sprintf("第%d章: 旅の始まり", chapterNum)
	case "zh-CN":
		return fmt.Sprintf("第%d章: 旅程的开始", chapterNum)
	case "id":
		return fmt.Sprintf("Bab %d: Permulaan Perjalanan", chapterNum)
	default: // en
		return fmt.Sprintf("Chapter %d: The Beginning of the Journey", chapterNum)
	}
}

// Helper function to generate sample chapter content
func generateChapterContent(language string, chapterNum int) string {
	switch language {
	case "ko":
		return fmt.Sprintf(`이것은 %d화의 내용입니다. 이 장에서 주인공은 그들의 능력과 결의를 시험할 새로운 도전에 직면합니다.

그는 절벽 가장자리에 서서 아래의 어둠을 응시했습니다. 밤바람이 거세게 불어 그의 로브를 휘날리게 했습니다. 그는 한 번의 잘못된 발걸음이 그의 여정의 끝을 의미할 수 있다는 것을 알고 있었습니다.

"나는 계속해야 해," 그는 혼잣말을 했습니다. "너무 많은 것이 나에게 달려 있어."

확고한 발걸음으로 그는 자신의 운명을 향해 나아갔습니다...

[장의 나머지 부분은 흥미진진한 모험과 깊은 캐릭터 개발로 계속됩니다. 독자들은 주인공이 다양한 장애물과 점점 더 강력한 적들에 맞서면서 그의 성장을 목격할 것입니다.]`, chapterNum)

	case "ja":
		return fmt.Sprintf(`これは第%d章の内容です。この章では、主人公は彼らの能力と決意を試す新しい挑戦に直面します。

彼は崖の端に立ち、下の暗闇を見つめていました。夜風が激しく吹き、彼のローブをはためかせました。一歩間違えば、旅の終わりを意味することを彼は知っていました。

「続けなければならない」と彼は独り言を言いました。「多くのことが私にかかっている」

確固たる足取りで、彼は自分の運命に向かって前進しました...

[章の残りの部分は、エキサイティングな冒険と深いキャラクター開発で続きます。読者は、主人公がさまざまな障害とますます強力な敵に立ち向かう中で、彼らの成長を目撃するでしょう。]`, chapterNum)

	case "zh-CN":
		return fmt.Sprintf(`这是第%d章的内容。在本章中，主角面临着考验他们能力和决心的新挑战。

他站在悬崖边缘，凝视着下面的黑暗。夜风猛烈地吹着，使他的长袍飘动。他知道一步错误可能意味着他旅程的终结。

"我必须继续，"他自言自语道。"太多事情依赖于我。"

以坚定的步伐，他朝着自己的命运前进...

[本章的其余部分将继续精彩的冒险和深刻的角色发展。读者将见证主角在面对各种障碍和日益强大的敌人时的成长。]`, chapterNum)

	case "id":
		return fmt.Sprintf(`Ini adalah konten dari bab %d. Dalam bab ini, protagonis menghadapi tantangan baru yang akan menguji kemampuan dan tekadnya.

Dia berdiri di tepi jurang, menatap ke kegelapan di bawah. Angin malam berhembus kencang, membuat jubahnya berkibar. Dia tahu bahwa satu langkah salah bisa berarti akhir dari perjalanannya.

"Aku harus melanjutkan," gumamnya pada diri sendiri. "Terlalu banyak yang bergantung padaku."

Dengan langkah mantap, dia melangkah maju menuju takdirnya...

[Konten bab selanjutnya akan berlanjut dengan petualangan yang mendebarkan dan pengembangan karakter yang mendalam. Pembaca akan melihat pertumbuhan protagonis saat ia menghadapi berbagai rintangan dan musuh yang semakin kuat.]`, chapterNum)

	default: // en
		return fmt.Sprintf(`This is the content of chapter %d. In this chapter, the protagonist faces new challenges that will test their abilities and resolve.

He stood at the edge of the cliff, staring into the darkness below. The night wind blew fiercely, making his robe flutter. He knew that one wrong step could mean the end of his journey.

"I must continue," he muttered to himself. "Too much depends on me."

With steady steps, he moved forward toward his destiny...

[The rest of the chapter will continue with exciting adventures and deep character development. Readers will witness the protagonist's growth as they face various obstacles and increasingly powerful enemies.]`, chapterNum)
	}
}
