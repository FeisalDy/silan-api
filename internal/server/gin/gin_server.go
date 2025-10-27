package gin

import (
	"context"
	"log"
	"simple-go/internal/middleware"
	"simple-go/internal/server"

	"github.com/gin-gonic/gin"
)

// GinServer implements the Server interface using Gin framework
type GinServer struct {
	engine *gin.Engine
	config *server.Config
}

// NewGinServer creates a new Gin server with all routes configured
func NewGinServer(cfg *server.Config) server.Server {
	// Set Gin mode
	if cfg.Config.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Global middleware
	router.Use(middleware.Logger())
	router.Use(middleware.ErrorHandler())
	router.Use(middleware.CORS())
	router.Use(gin.Recovery())

	// Setup routes
	setupRoutes(router, cfg)

	return &GinServer{
		engine: router,
		config: cfg,
	}
}

// Start starts the Gin server
func (s *GinServer) Start(addr string) error {
	log.Printf("Starting Gin server on %s", addr)
	return s.engine.Run(addr)
}

// Shutdown gracefully shuts down the Gin server
func (s *GinServer) Shutdown() error {
	log.Println("Shutting down Gin server...")
	// Gin doesn't have built-in graceful shutdown in the simple API
	// For graceful shutdown, you'd need to use http.Server directly
	return nil
}

// setupRoutes configures all routes for the Gin server
func setupRoutes(router *gin.Engine, cfg *server.Config) {
	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Public routes (no authentication)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", cfg.AuthHandler.Register)
			auth.POST("/login", cfg.AuthHandler.Login)
		}

		// Role getter function for Casbin (shared by all protected routes)
		roleGetter := func(c *gin.Context) ([]string, error) {
			userID, exists := middleware.GetUserID(c)
			if !exists {
				return nil, nil
			}
			return cfg.UserService.GetUserRoles(context.Background(), userID)
		}

		// Protected routes (require authentication)
		users := v1.Group("/users")
		users.Use(middleware.JWTAuth(cfg.JWTManager))
		users.Use(middleware.CasbinAuthorizer(cfg.Enforcer, roleGetter))
		{
			users.GET("", cfg.UserHandler.GetAll)
			users.GET("/:id", cfg.UserHandler.GetByID)
			users.PUT("/:id", cfg.UserHandler.Update)
			users.DELETE("/:id", cfg.UserHandler.Delete)
		}

		// Novel routes
		novels := v1.Group("/novels")
		novels.Use(middleware.JWTAuth(cfg.JWTManager))
		{
			// Public novel endpoints (authenticated users can view)
			novels.GET("", cfg.NovelHandler.GetAll)
			novels.GET("/my", cfg.NovelHandler.GetMyNovels)
			novels.GET("/:id", cfg.NovelHandler.GetByID)
			novels.GET("/:id/translations/:lang", cfg.NovelHandler.GetTranslation)

			// Protected novel endpoints (require authorization)
			novelsProtected := novels.Group("")
			novelsProtected.Use(middleware.CasbinAuthorizer(cfg.Enforcer, roleGetter))
			{
				novelsProtected.POST("", cfg.NovelHandler.Create)
				novelsProtected.PUT("/:id", cfg.NovelHandler.Update)
				novelsProtected.DELETE("/:id", cfg.NovelHandler.Delete)
				novelsProtected.POST("/translations", cfg.NovelHandler.CreateTranslation)
				novelsProtected.PUT("/translations/:translation_id", cfg.NovelHandler.UpdateTranslation)
				novelsProtected.DELETE("/translations/:translation_id", cfg.NovelHandler.DeleteTranslation)
			}
		}

		// Chapter routes
		chapters := v1.Group("/chapters")
		chapters.Use(middleware.JWTAuth(cfg.JWTManager))
		{
			// Public chapter endpoints (authenticated users can view)
			chapters.GET("", cfg.ChapterHandler.GetByNovel)
			chapters.GET("/search", cfg.ChapterHandler.GetByNovelAndNumber)
			chapters.GET("/:id", cfg.ChapterHandler.GetByID)
			chapters.GET("/:id/translations/:lang", cfg.ChapterHandler.GetTranslation)

			// Protected chapter endpoints (require authorization)
			chaptersProtected := chapters.Group("")
			chaptersProtected.Use(middleware.CasbinAuthorizer(cfg.Enforcer, roleGetter))
			{
				chaptersProtected.POST("", cfg.ChapterHandler.Create)
				chaptersProtected.PUT("/:id", cfg.ChapterHandler.Update)
				chaptersProtected.DELETE("/:id", cfg.ChapterHandler.Delete)
				chaptersProtected.POST("/translations", cfg.ChapterHandler.CreateTranslation)
				chaptersProtected.PUT("/translations/:translation_id", cfg.ChapterHandler.UpdateTranslation)
				chaptersProtected.DELETE("/translations/:translation_id", cfg.ChapterHandler.DeleteTranslation)
			}
		}
	}
}
