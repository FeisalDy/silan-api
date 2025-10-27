# Database Seeders

This directory contains database seeders for the Simple Go application. Seeders populate the database with sample data for development and testing purposes.

## Available Seeders

1. **Roles** (`seed_roles.go`) - Creates default user roles
   - admin
   - author
   - translator
   - user

2. **Users** (`seed_users.go`) - Creates sample users with assigned roles
   - admin@example.com (admin role)
   - author1@example.com (author role)
   - translator1@example.com (translator role)
   - john@example.com (user role)
   - Default password for all: `password123`

3. **Genres** (`seed_genres.go`) - Creates novel genre categories
   - Fantasy, Romance, Action, Mystery, Sci-Fi, Horror, Comedy, Drama, etc.

4. **Tags** (`seed_tags.go`) - Creates novel tags
   - Reincarnation, Overpowered MC, System, Magic, Cultivation, etc.

5. **Novels** (`seed_novels.go`) - Creates sample novels with translations
   - 3 sample novels with English and Indonesian translations
   - Associated genres and tags

6. **Chapters** (`seed_chapters.go`) - Creates sample chapters for novels
   - 5 chapters per novel with translations in English and Indonesian

## Running Seeders

### Method 1: Using Make Command
```bash
make seed
```

### Method 2: Direct Command
```bash
go run cmd/seed/main.go
```

### Method 3: Build and Run
```bash
go build -o bin/seed cmd/seed/main.go
./bin/seed
```

## Seeding Order

Seeders run in the following order (dependency-based):
1. Roles
2. Users
3. Genres
4. Tags
5. Novels
6. Chapters

## Features

### Idempotent Seeding
- Seeders check if data already exists before creating
- Safe to run multiple times
- Won't create duplicate records

### Sample Data Quality
- Realistic novel titles and descriptions
- Multi-language support (English and Indonesian)
- Proper relationships (users, roles, novels, chapters, genres, tags)
- Sample chapter content

## Architecture Integration

The seeders follow the project's domain-driven architecture:

```
seeds/
├── seeds.go              # Main orchestrator
├── seed_roles.go         # Uses internal/domain/role
├── seed_users.go         # Uses internal/domain/user
├── seed_genres.go        # Uses internal/domain/genre
├── seed_tags.go          # Uses internal/domain/tag
├── seed_novels.go        # Uses internal/domain/novel
└── seed_chapters.go      # Uses internal/domain/chapter
```

All seeders import and use the actual domain models from `internal/domain/*`, ensuring consistency with the main application.

## Customization

### Adding New Sample Data

Edit the respective seed file and add new entries to the slice:

```go
// In seed_novels.go
novels := []struct {
    Novel        novel.Novel
    Translations []novel.NovelTranslation
    GenreSlugs   []string
    TagSlugs     []string
}{
    {
        Novel: novel.Novel{
            OriginalLanguage: "en",
            OriginalAuthor:   strPtr("Your Name"),
            Status:           strPtr("ongoing"),
            CreatedBy:        author.ID,
        },
        Translations: []novel.NovelTranslation{
            {
                Lang:        "en",
                Title:       "Your Novel Title",
                Description: strPtr("Description here"),
            },
        },
        GenreSlugs: []string{"fantasy", "action"},
        TagSlugs:   []string{"magic", "adventure"},
    },
}
```

### Creating New Seeders

1. Create a new file: `seed_yourmodel.go`
2. Implement the seeder function:

```go
package seeds

import (
    "log"
    "simple-go/internal/domain/yourmodel"
    "gorm.io/gorm"
)

func SeedYourModel(db *gorm.DB) error {
    log.Println("🌱 Seeding your model...")
    
    // Your seeding logic here
    
    log.Println("✅ Your model seeding completed")
    return nil
}
```

3. Add it to `seeds.go`:

```go
seeders := []struct {
    name string
    fn   func(*gorm.DB) error
}{
    {"Roles", SeedRoles},
    // ... other seeders
    {"Your Model", SeedYourModel},
}
```

## Database Requirements

Seeders require:
- PostgreSQL database running (see docker-compose.yml)
- Database migrations completed (run automatically on app startup)
- Valid .env configuration

## Testing

After seeding, you can verify the data:

```bash
# Connect to PostgreSQL
psql -h localhost -U postgres -d simple_go

# Check seeded data
SELECT * FROM users;
SELECT * FROM roles;
SELECT * FROM novels;
SELECT * FROM novel_translations;
SELECT * FROM chapters;
SELECT * FROM chapter_translations;
```

## Common Issues

### "Author user not found"
- Ensure roles and users are seeded before novels
- The seeder will skip novel seeding if author user doesn't exist

### "No novels found"
- Run novel seeder before chapter seeder
- Chapters require existing novels

### Duplicate Key Errors
- Seeders are idempotent but may fail if data was manually modified
- Drop and recreate the database if needed:
  ```bash
  make docker-down
  make docker-up
  make seed
  ```

## Helper Functions

Available in `seeds.go`:

- `strPtr(s string) *string` - Creates string pointer
- `intPtr(i int) *int` - Creates int pointer

Use these for nullable fields:

```go
User{
    Bio: strPtr("Sample bio"),
    WordCount: intPtr(5000),
}
```
