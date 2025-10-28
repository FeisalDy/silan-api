package app

import (
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
	MediaService   *service.MediaService
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

	// Initialize repositories
	userRepo := gormrepo.NewUserRepository(db)
	roleRepo := gormrepo.NewRoleRepository(db)
	novelRepo := gormrepo.NewNovelRepository(db)
	chapterRepo := gormrepo.NewChapterRepository(db)
	mediaRepo := gormrepo.NewMediaRepository(db)
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

	// Initialize upload + media services (ImgBB)
	uploadService := service.NewUploadService(nil, cfg.Media.ImgBBAPIKey, cfg.Media.ImgBBTTL)
	mediaService := service.NewMediaService(mediaRepo, uploadService)

	// Initialize services
	authService := service.NewAuthService(uow, userRepo, roleRepo, jwtManager)
	userService := service.NewUserService(userRepo, roleRepo)
	novelService := service.NewNovelService(uow, novelRepo, mediaService)
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
		MediaService:   mediaService,
		JWTManager:     jwtManager,
		Enforcer:       enforcer,
	}, nil
}
