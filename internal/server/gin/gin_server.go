package gin

import (
	"context"
	"log"
	"simple-go/internal/middleware"
	"simple-go/internal/server"

	"github.com/gin-gonic/gin"
)

type GinServer struct {
	engine *gin.Engine
	config *server.Config
}

func NewGinServer(cfg *server.Config) server.Server {
	if cfg.Config.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	router.Use(middleware.Logger())
	router.Use(middleware.ErrorHandler())
	router.Use(middleware.CORS())
	router.Use(gin.Recovery())

	setupRoutes(router, cfg)

	return &GinServer{
		engine: router,
		config: cfg,
	}
}

func (s *GinServer) Start(addr string) error {
	log.Printf("Starting Gin server on %s", addr)
	return s.engine.Run(addr)
}

func (s *GinServer) Shutdown() error {
	log.Println("Shutting down Gin server...")
	return nil
}

func setupRoutes(router *gin.Engine, cfg *server.Config) {
	v1 := router.Group("/api/v1")
	{
		router.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})

		auth := v1.Group("/auth")
		{
			auth.POST("/register", cfg.AuthHandler.Register)
			auth.POST("/login", cfg.AuthHandler.Login)
			auth.GET("/profile", middleware.JWTAuth(cfg.JWTManager), cfg.AuthHandler.GetProfile)
		}

		// Role getter function for Casbin (shared by all protected routes)
		roleGetter := func(c *gin.Context) ([]string, error) {
			userID, exists := middleware.GetUserID(c)
			if !exists {
				return nil, nil
			}
			return cfg.UserService.GetUserRoles(context.Background(), userID)
		}

		users := v1.Group("/users")
		users.Use(middleware.JWTAuth(cfg.JWTManager))
		{
			users.GET("", middleware.RequirePermission("user", "list", cfg.Enforcer, roleGetter), cfg.UserHandler.GetAll)
			users.GET("/:id", middleware.RequirePermission("user", "read", cfg.Enforcer, roleGetter), cfg.UserHandler.GetByID)
			users.PATCH("/:id", middleware.RequirePermission("user", "update", cfg.Enforcer, roleGetter), cfg.UserHandler.Update)
			users.DELETE("/:id", middleware.RequirePermission("user", "delete", cfg.Enforcer, roleGetter), cfg.UserHandler.Delete)
		}

		// Novel routes
		novels := v1.Group("/novels")
		novels.GET("", cfg.NovelHandler.GetAll)
		novels.GET("/:id", cfg.NovelHandler.GetByID)
		novels.GET("/:id/volumes", cfg.NovelHandler.GetNovelVolumes)
		novels.Use(middleware.JWTAuth(cfg.JWTManager))
		{
			novels.POST("", middleware.RequirePermission("novel", "create", cfg.Enforcer, roleGetter), cfg.NovelHandler.Create)
			novels.POST("/epub", middleware.RequirePermission("novel", "create", cfg.Enforcer, roleGetter), cfg.NovelHandler.UploadEpub)
			novels.DELETE("/:id", middleware.RequirePermission("novel", "delete", cfg.Enforcer, roleGetter), cfg.NovelHandler.Delete)
			novels.PATCH("/:id/cover", middleware.RequirePermission("novel", "update", cfg.Enforcer, roleGetter), cfg.NovelHandler.UpdateCoverMedia)

			novels.POST("/translations", middleware.RequirePermission("novel_translation", "create", cfg.Enforcer, roleGetter), cfg.NovelHandler.CreateTranslation)
			novels.DELETE("/translations/:translation_id", middleware.RequirePermission("novel_translation", "delete", cfg.Enforcer, roleGetter), cfg.NovelHandler.DeleteTranslation)
		}

		chapters := v1.Group("/chapters")
		chapters.GET("/:id", cfg.ChapterHandler.GetByID)
		chapters.Use(middleware.JWTAuth(cfg.JWTManager))
		{

			chapters.POST("", middleware.RequirePermission("chapter", "create", cfg.Enforcer, roleGetter), cfg.ChapterHandler.Create)
			chapters.DELETE("/:id", middleware.RequirePermission("chapter", "delete", cfg.Enforcer, roleGetter), cfg.ChapterHandler.Delete)

			chapters.POST("/translations", middleware.RequirePermission("chapter_translation", "create", cfg.Enforcer, roleGetter), cfg.ChapterHandler.CreateTranslation)
			chapters.DELETE("/translations/:id", middleware.RequirePermission("chapter_translation", "delete", cfg.Enforcer, roleGetter), cfg.ChapterHandler.DeleteTranslation)
		}

		jobs := v1.Group("/translation-jobs")
		jobs.Use(middleware.JWTAuth(cfg.JWTManager))
		{
			jobs.POST("", middleware.RequirePermission("translation_job", "create", cfg.Enforcer, roleGetter), cfg.TranslationJobHandler.CreateTranslationJob)
			jobs.GET("", middleware.RequirePermission("translation_job", "list", cfg.Enforcer, roleGetter), cfg.TranslationJobHandler.GetAllJobs)
			jobs.GET("/:id", middleware.RequirePermission("translation_job", "read", cfg.Enforcer, roleGetter), cfg.TranslationJobHandler.GetJobByID)
			jobs.PUT("/:id/cancel", middleware.RequirePermission("translation_job", "update", cfg.Enforcer, roleGetter), cfg.TranslationJobHandler.CancelJob)
		}

		// Miscellaneous routes
		misc := v1.Group("/miscellaneous")
		{
			misc.GET("/languages", cfg.MiscellaneousHandler.GetLanguages)
		}
	}
}
