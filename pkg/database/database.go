package database

import (
	"fmt"
	"log"
	"simple-go/internal/domain/chapter"
	"simple-go/internal/domain/genre"
	"simple-go/internal/domain/job"
	"simple-go/internal/domain/media"
	"simple-go/internal/domain/novel"
	novelgenre "simple-go/internal/domain/novel_genre"
	noveltag "simple-go/internal/domain/novel_tag"
	"simple-go/internal/domain/role"
	"simple-go/internal/domain/tag"
	"simple-go/internal/domain/user"
	userrole "simple-go/internal/domain/user_role"
	"simple-go/internal/domain/volume"

	"simple-go/pkg/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func Connect(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	dsn := cfg.DSN()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// Connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	log.Println("Database connection established")

	if err := migrateDatabase(db); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)

	}

	if err := addTrigramIndexes(db); err != nil {
		return nil, fmt.Errorf("failed to add trigram indexes: %w", err)
	}

	return db, nil
}

func migrateDatabase(db *gorm.DB) error {
	log.Println("Running database migrations...")

	err := db.AutoMigrate(
		&user.User{},
		&role.Role{},
		&userrole.UserRole{},
		&genre.Genre{},
		&tag.Tag{},
		&media.Media{},
		&novel.Novel{},
		&noveltag.NovelTag{},
		&novelgenre.NovelGenre{},
		&novel.NovelTranslation{},
		&volume.Volume{},
		&volume.VolumeTranslation{},
		&chapter.Chapter{},
		&chapter.ChapterTranslation{},
		&job.TranslationJob{},
		&job.TranslationSubtask{},
	)

	if err != nil {
		return err
	}

	log.Println("Database migrations completed successfully")

	return nil
}

type TrigramIndex struct {
	Table  string
	Column string
}

var trigramIndexes = []TrigramIndex{
	{"novel_translations", "title"},
	{"users", "username"},
}

func addTrigramIndexes(db *gorm.DB) error {
	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS pg_trgm;`).Error; err != nil {
		return fmt.Errorf("failed to enable pg_trgm: %w", err)
	}

	for _, idx := range trigramIndexes {
		indexName := fmt.Sprintf("idx_%s_%s_trgm", idx.Table, idx.Column)
		sql := fmt.Sprintf(`
			DO $$
			BEGIN
				IF NOT EXISTS (
					SELECT 1 FROM pg_indexes
					WHERE schemaname = 'public'
					AND tablename = '%s'
					AND indexname = '%s'
				) THEN
					CREATE INDEX %s
					ON %s USING GIN (%s gin_trgm_ops);
				END IF;
			END$$;
		`, idx.Table, indexName, indexName, idx.Table, idx.Column)

		if err := db.Exec(sql).Error; err != nil {
			return fmt.Errorf("failed to create trigram index for %s.%s: %w",
				idx.Table, idx.Column, err)
		}
	}

	return nil
}
