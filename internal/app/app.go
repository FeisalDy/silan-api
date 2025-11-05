package app

import (
	"simple-go/internal/handler"
	"simple-go/internal/repository/gormrepo"
	"simple-go/internal/service"
	"simple-go/pkg/auth"
	casbinpkg "simple-go/pkg/casbin"
	"simple-go/pkg/config"
	"simple-go/pkg/database"
	"simple-go/pkg/logger"
	"simple-go/pkg/queue"

	"github.com/casbin/casbin/v2"
	"gorm.io/gorm"
)

type App struct {
	Config                *config.Config
	DB                    *gorm.DB
	AuthHandler           *handler.AuthHandler
	UserHandler           *handler.UserHandler
	NovelHandler          *handler.NovelHandler
	TranslationJobHandler *handler.TranslationJobHandler
	ChapterHandler        *handler.ChapterHandler
	MiscellaneousHandler  *handler.MiscellaneousHandler
	UserService           *service.UserService
	MediaService          *service.MediaService
	JWTManager            *auth.JWTManager
	Enforcer              *casbin.Enforcer
}

func Initialize() (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	db, err := database.Connect(&cfg.Database)
	if err != nil {
		return nil, err
	}

	userRepo := gormrepo.NewUserRepository(db)
	roleRepo := gormrepo.NewRoleRepository(db)
	novelRepo := gormrepo.NewNovelRepository(db)
	volumeRepo := gormrepo.NewVolumeRepository(db)
	chapterRepo := gormrepo.NewChapterRepository(db)
	mediaRepo := gormrepo.NewMediaRepository(db)
	jobRepo := gormrepo.NewTranslationJobRepository(db)
	uow := gormrepo.NewUnitOfWork(db)

	enforcer, err := casbinpkg.NewEnforcer(db, cfg.Casbin.ModelPath)
	if err != nil {
		return nil, err
	}

	if err := casbinpkg.InitializeDefaultPolicies(enforcer); err != nil {
		return nil, err
	}

	jwtManager := auth.NewJWTManager(cfg.JWT.Secret, cfg.JWT.ExpirationHours)

	uploadService := service.NewUploadService(nil, cfg.Media.ImgBBAPIKey, cfg.Media.ImgBBTTL)
	mediaService := service.NewMediaService(mediaRepo, uploadService)
	epubService := service.NewEpubService()

	// Initialize Redis queue (optional - gracefully handle failure)
	var redisQueue *queue.RedisQueue
	redisQueue, err = queue.NewRedisQueue(cfg.Redis.URL, cfg.Redis.QueueName)
	if err != nil {
		logger.Warn("Failed to connect to Redis, translation jobs will not be queued: " + err.Error())
		redisQueue = nil // Service will work without queue
	}

	// Initialize services
	authService := service.NewAuthService(uow, userRepo, roleRepo, jwtManager)
	userService := service.NewUserService(userRepo, roleRepo)
	volumeService := service.NewVolumeService(uow, volumeRepo, chapterRepo, mediaService)
	novelService := service.NewNovelService(uow, novelRepo, mediaService, volumeService, epubService)
	chapterService := service.NewChapterService(uow, chapterRepo)
	jobService := service.NewTranslationJobService(uow, jobRepo, novelRepo, volumeRepo, chapterRepo, redisQueue)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	novelHandler := handler.NewNovelHandler(novelService)
	chapterHandler := handler.NewChapterHandler(chapterService, volumeService)
	translationJobHandler := handler.NewTranslationJobHandler(jobService)
	miscellaneousHandler := handler.NewMiscellaneousHandler()

	return &App{
		Config:                cfg,
		DB:                    db,
		AuthHandler:           authHandler,
		UserHandler:           userHandler,
		NovelHandler:          novelHandler,
		ChapterHandler:        chapterHandler,
		UserService:           userService,
		MediaService:          mediaService,
		TranslationJobHandler: translationJobHandler,
		MiscellaneousHandler:  miscellaneousHandler,
		JWTManager:            jwtManager,
		Enforcer:              enforcer,
	}, nil
}
