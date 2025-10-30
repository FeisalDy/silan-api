package seeds

import (
	"log"
	"simple-go/internal/domain/role"
	"simple-go/internal/domain/user"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// SeedUsers seeds default users with roles
func SeedUsers(db *gorm.DB) error {
	log.Println("🌱 Seeding users...")

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	passwordStr := string(hashedPassword)

	users := []struct {
		User     user.User
		RoleName string
	}{
		{
			User: user.User{
				Username:  strPtr("admin"),
				Email:     "admin@example.com",
				Password:  &passwordStr,
				AvatarURL: strPtr("https://ui-avatars.com/api/?name=Admin&background=random"),
				Bio:       strPtr("System administrator"),
				Status:    "active",
			},
			RoleName: "admin",
		},
		{
			User: user.User{
				Username:  strPtr("author1"),
				Email:     "author1@example.com",
				Password:  &passwordStr,
				AvatarURL: strPtr("https://ui-avatars.com/api/?name=Author+One&background=random"),
				Bio:       strPtr("Professional novelist and content creator"),
				Status:    "active",
			},
			RoleName: "author",
		},
		{
			User: user.User{
				Username:  strPtr("translator1"),
				Email:     "translator1@example.com",
				Password:  &passwordStr,
				AvatarURL: strPtr("https://ui-avatars.com/api/?name=Translator+One&background=random"),
				Bio:       strPtr("Professional translator"),
				Status:    "active",
			},
			RoleName: "translator",
		},
		{
			User: user.User{
				Username:  strPtr("john_doe"),
				Email:     "john@example.com",
				Password:  &passwordStr,
				AvatarURL: strPtr("https://ui-avatars.com/api/?name=John+Doe&background=random"),
				Bio:       strPtr("Novel enthusiast and reader"),
				Status:    "active",
			},
			RoleName: "user",
		},
	}

	for _, userData := range users {
		var existing user.User
		result := db.Where("email = ?", userData.User.Email).First(&existing)

		if result.Error == gorm.ErrRecordNotFound {
			// Create user
			if err := db.Create(&userData.User).Error; err != nil {
				log.Printf("⚠️  Failed to seed user %s: %v", userData.User.Email, err)
				return err
			}

			// Assign role using GORM association (User owns the relationship)
			var r role.Role
			if err := db.Where("name = ?", userData.RoleName).First(&r).Error; err != nil {
				log.Printf("⚠️  Role %s not found for user %s", userData.RoleName, userData.User.Email)
				continue
			}

			// Use Association on User model since User owns the M2M relationship
			if err := db.Model(&userData.User).Association("Roles").Append(&r); err != nil {
				log.Printf("⚠️  Failed to assign role to user %s: %v", userData.User.Email, err)
			}

			log.Printf("✅ Created user: %s with role: %s", userData.User.Email, userData.RoleName)
		} else if result.Error != nil {
			return result.Error
		} else {
			log.Printf("⏭️  User already exists: %s", userData.User.Email)
		}
	}

	log.Println("✅ Users seeding completed")
	return nil
}
