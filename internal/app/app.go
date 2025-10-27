package app

import (
	"log"
	"simple-go/internal/domain/chapter"
	"simple-go/internal/domain/genre"
	"simple-go/internal/domain/novel"
	"simple-go/internal/domain/role"
	"simple-go/internal/domain/tag"
	"simple-go/internal/domain/user"
	"simple-go/internal/handler"
	"simple-go/internal/repository/gormrepo"
	"simple-go/internal/service"
	"simple-go/pkg/auth"
	casbinpkg "simple-go/pkg/casbin"
	"simple-go/pkg/config"
	"simple-go/pkg/database"

	"github.com/casbin/casbin/v2"
	"gorm.io/gorm"
)

// App represents the application with all its dependencies
type App struct {
	Config         *config.Config
	DB             *gorm.DB
	AuthHandler    *handler.AuthHandler
	UserHandler    *handler.UserHandler
	NovelHandler   *handler.NovelHandler
	ChapterHandler *handler.ChapterHandler
	UserService    *service.UserService
	JWTManager     *auth.JWTManager
	Enforcer       *casbin.Enforcer
}

// Initialize sets up the application with all dependencies
func Initialize() (*App, error) {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	// Connect to database
	db, err := database.Connect(&cfg.Database)
	if err != nil {
		return nil, err
	}

	// Auto-migrate database schema
	if err := migrateDatabase(db); err != nil {
		return nil, err
	}

	// Initialize repositories
	userRepo := gormrepo.NewUserRepository(db)
	roleRepo := gormrepo.NewRoleRepository(db)
	novelRepo := gormrepo.NewNovelRepository(db)
	chapterRepo := gormrepo.NewChapterRepository(db)
	uow := gormrepo.NewUnitOfWork(db)

	// Initialize Casbin
	enforcer, err := casbinpkg.NewEnforcer(db, cfg.Casbin.ModelPath)
	if err != nil {
		return nil, err
	}

	// Initialize default Casbin policies
	if err := casbinpkg.InitializeDefaultPolicies(enforcer); err != nil {
		return nil, err
	}

	// Initialize JWT manager
	jwtManager := auth.NewJWTManager(cfg.JWT.Secret, cfg.JWT.ExpirationHours)

	// Initialize services
	authService := service.NewAuthService(uow, userRepo, roleRepo, jwtManager)
	userService := service.NewUserService(userRepo, roleRepo)
	novelService := service.NewNovelService(uow, novelRepo)
	chapterService := service.NewChapterService(uow, chapterRepo, novelRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	novelHandler := handler.NewNovelHandler(novelService)
	chapterHandler := handler.NewChapterHandler(chapterService)

	return &App{
		Config:         cfg,
		DB:             db,
		AuthHandler:    authHandler,
		UserHandler:    userHandler,
		NovelHandler:   novelHandler,
		ChapterHandler: chapterHandler,
		UserService:    userService,
		JWTManager:     jwtManager,
		Enforcer:       enforcer,
	}, nil
}

// migrateDatabase runs all database migrations
func migrateDatabase(db *gorm.DB) error {
	log.Println("Running database migrations...")

	err := db.AutoMigrate(
		&user.User{},
		&role.Role{},
		&role.UserRole{},
		&genre.Genre{},
		&genre.NovelGenre{},
		&tag.Tag{},
		&tag.NovelTag{},
		&novel.Novel{},
		&novel.NovelTranslation{},
		&chapter.Chapter{},
		&chapter.ChapterTranslation{},
	)

	if err != nil {
		return err
	}

	log.Println("Database migrations completed successfully")
	return nil
}
