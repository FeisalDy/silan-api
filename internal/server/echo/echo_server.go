package echoserver

import (
	"context"
	"log"
	"net/http"
	"simple-go/internal/server"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
)

// EchoServer implements the Server interface using Echo framework
type EchoServer struct {
	engine *echo.Echo
	config *server.Config
}

// NewEchoServer creates a new Echo server with all routes configured
func NewEchoServer(cfg *server.Config) server.Server {
	e := echo.New()

	// Global middleware
	e.Use(echomiddleware.Logger())
	e.Use(echomiddleware.Recover())
	e.Use(echomiddleware.CORS())

	// Setup routes
	setupRoutes(e, cfg)

	return &EchoServer{
		engine: e,
		config: cfg,
	}
}

// Start starts the Echo server
func (s *EchoServer) Start(addr string) error {
	log.Printf("Starting Echo server on %s", addr)
	return s.engine.Start(addr)
}

// Shutdown gracefully shuts down the Echo server
func (s *EchoServer) Shutdown() error {
	log.Println("Shutting down Echo server...")
	return s.engine.Shutdown(context.Background())
}

// setupRoutes configures all routes for the Echo server
func setupRoutes(e *echo.Echo, cfg *server.Config) {
	// Health check
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	// API v1 routes
	v1 := e.Group("/api/v1")

	// Public routes (no authentication)
	auth := v1.Group("/auth")
	{
		// Note: You would need to create Echo-specific adapters for handlers
		// or modify handlers to work with both frameworks
		// This is just a structure example
		auth.POST("/register", echoAdapter(cfg.AuthHandler.Register))
		auth.POST("/login", echoAdapter(cfg.AuthHandler.Login))
	}

	// Role getter function for Casbin
	roleGetter := func(c echo.Context) ([]string, error) {
		// Note: This would need to be adapted for Echo context
		userID := c.Get("user_id").(string)
		return cfg.UserService.GetUserRoles(context.Background(), userID)
	}

	// Protected routes (require authentication)
	users := v1.Group("/users")
	users.Use(echoJWTMiddleware(cfg.JWTManager))
	users.Use(echoCasbinMiddleware(cfg.Enforcer, roleGetter))
	{
		users.GET("", echoAdapter(cfg.UserHandler.GetAll))
		users.GET("/:id", echoAdapter(cfg.UserHandler.GetByID))
		users.PUT("/:id", echoAdapter(cfg.UserHandler.Update))
		users.DELETE("/:id", echoAdapter(cfg.UserHandler.Delete))
	}

	// Similar pattern for novels and chapters...
}

// echoAdapter adapts Gin handlers to Echo handlers
// This is a placeholder - you would need to implement proper adapter logic
func echoAdapter(ginHandler interface{}) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Implement adapter logic here
		return c.JSON(http.StatusNotImplemented, map[string]string{
			"error": "Handler adapter not implemented yet",
		})
	}
}

// echoJWTMiddleware adapts JWT middleware for Echo
func echoJWTMiddleware(jwtManager interface{}) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Implement JWT validation logic here
			return next(c)
		}
	}
}

// echoCasbinMiddleware adapts Casbin middleware for Echo
func echoCasbinMiddleware(enforcer interface{}, roleGetter func(echo.Context) ([]string, error)) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Implement Casbin authorization logic here
			return next(c)
		}
	}
}
