# Server Architecture Documentation

## Overview

The application has been refactored to use a **framework-agnostic server abstraction pattern**. This design allows you to easily swap web frameworks (Gin, Echo, Fiber, etc.) without changing business logic.

## Architecture Layers

```
cmd/api/main.go
    ↓
internal/app/app.go (Initialization)
    ↓
internal/server/server.go (Interface)
    ↓
internal/server/gin/gin_server.go (Implementation)
```

### 1. Entry Point (`cmd/api/main.go`)

**Purpose**: Clean entry point with minimal code (36 lines)

**Responsibilities**:
- Load configuration
- Initialize application dependencies
- Create server instance
- Start server
- Handle graceful shutdown

**Key Code**:
```go
func main() {
    // Initialize app and all dependencies
    app, err := app.Initialize(cfg)
    
    // Create server config
    serverConfig := &server.Config{
        AuthHandler:    app.AuthHandler,
        UserHandler:    app.UserHandler,
        // ... other handlers
    }
    
    // Create Gin server (can be swapped with Echo)
    srv := ginserver.NewGinServer(serverConfig)
    
    // Start server
    srv.Start(":8080")
}
```

### 2. Application Initialization (`internal/app/app.go`)

**Purpose**: Bootstrap all application dependencies

**Responsibilities**:
- Database connection setup
- Run database migrations
- Initialize repositories
- Initialize services
- Initialize handlers
- Setup JWT manager
- Setup Casbin enforcer

**Key Code**:
```go
func Initialize(cfg *config.Config) (*App, error) {
    // Connect to database
    db, err := database.InitDB(cfg)
    
    // Run migrations
    migrateDatabase(db)
    
    // Initialize repositories
    userRepo := repository.NewUserRepository(db)
    novelRepo := repository.NewNovelRepository(db)
    // ...
    
    // Initialize services
    userService := service.NewUserService(userRepo)
    // ...
    
    // Initialize handlers
    authHandler := handler.NewAuthHandler(userService, jwtManager)
    // ...
    
    return &App{
        AuthHandler: authHandler,
        // ... all dependencies
    }, nil
}
```

### 3. Server Interface (`internal/server/server.go`)

**Purpose**: Define framework-agnostic server contract

**Interface Definition**:
```go
type Server interface {
    Start(addr string) error
    Shutdown() error
}
```

**Config Struct**:
```go
type Config struct {
    AuthHandler    *handler.AuthHandler
    UserHandler    *handler.UserHandler
    NovelHandler   *handler.NovelHandler
    ChapterHandler *handler.ChapterHandler
    JWTManager     *auth.JWTManager
    Enforcer       *casbin.Enforcer
    UserService    service.UserService
}
```

### 4. Framework Implementation (`internal/server/gin/gin_server.go`)

**Purpose**: Gin-specific implementation of Server interface

**Structure**:
```go
type GinServer struct {
    engine *gin.Engine
    config *server.Config
}

func NewGinServer(cfg *server.Config) server.Server {
    gin.SetMode(gin.ReleaseMode)
    r := gin.Default()
    
    setupRoutes(r, cfg)
    
    return &GinServer{
        engine: r,
        config: cfg,
    }
}

func (s *GinServer) Start(addr string) error {
    return s.engine.Run(addr)
}

func (s *GinServer) Shutdown() error {
    // Graceful shutdown logic
    return nil
}
```

## How to Swap Frameworks

### Example: Switching from Gin to Echo

**Step 1**: Install Echo
```bash
go get github.com/labstack/echo/v4
```

**Step 2**: Create Echo implementation
Create `internal/server/echo/echo_server.go`:

```go
package echoserver

import (
    "context"
    "simple-go/internal/server"
    "github.com/labstack/echo/v4"
)

type EchoServer struct {
    engine *echo.Echo
    config *server.Config
}

func NewEchoServer(cfg *server.Config) server.Server {
    e := echo.New()
    setupRoutes(e, cfg)
    return &EchoServer{engine: e, config: cfg}
}

func (s *EchoServer) Start(addr string) error {
    return s.engine.Start(addr)
}

func (s *EchoServer) Shutdown() error {
    return s.engine.Shutdown(context.Background())
}

func setupRoutes(e *echo.Echo, cfg *server.Config) {
    // Configure all routes with Echo syntax
    e.GET("/health", healthCheckHandler)
    // ... other routes
}
```

**Step 3**: Update main.go (only 2 lines change)
```go
// Before (Gin)
import "simple-go/internal/server/gin"
srv := ginserver.NewGinServer(serverConfig)

// After (Echo)
import "simple-go/internal/server/echo"
srv := echoserver.NewEchoServer(serverConfig)
```

**That's it!** All business logic remains unchanged.

## Benefits

### 1. **Framework Independence**
- Business logic is completely decoupled from web framework
- Easy to migrate or test different frameworks
- Can run multiple framework implementations side-by-side

### 2. **Maintainability**
- Clear separation of concerns
- Each layer has a single responsibility
- Easy to locate and modify code

### 3. **Testability**
- Can mock the Server interface for testing
- Business logic can be tested without web framework
- Can create test implementations of Server

### 4. **Scalability**
- Can add new frameworks without modifying existing code
- Can support multiple servers (HTTP, gRPC, WebSocket) simultaneously
- Easy to add new features at the right layer

## Project Structure

```
simple-go/
├── cmd/
│   └── api/
│       └── main.go                    # Entry point (36 lines)
├── internal/
│   ├── app/
│   │   └── app.go                     # Application initialization
│   ├── server/
│   │   ├── server.go                  # Server interface
│   │   ├── gin/
│   │   │   └── gin_server.go          # Gin implementation
│   │   └── echo/                      # (Optional)
│   │       └── echo_server.go         # Echo implementation
│   ├── domain/                        # Business entities
│   ├── handler/                       # HTTP handlers
│   ├── service/                       # Business logic
│   ├── repository/                    # Data access
│   ├── middleware/                    # Shared middleware
│   └── auth/                          # Authentication
└── ...
```

## Handler Compatibility

Handlers are designed to be framework-agnostic but currently use Gin's context. To fully support multiple frameworks, you have two options:

### Option 1: Framework-Specific Handlers
Create separate handlers for each framework:
- `handler/gin/auth_handler.go`
- `handler/echo/auth_handler.go`

### Option 2: Adapter Pattern (Recommended)
Keep handlers framework-agnostic and create adapters:
- Handlers use a custom `Context` interface
- Each framework has an adapter that converts framework context to custom context
- Most flexible but requires more initial setup

## Middleware

Current middleware (JWT, Casbin) is Gin-specific. For multi-framework support:

1. Define middleware interfaces in `internal/middleware/`
2. Implement framework-specific middleware in respective server packages
3. Register middleware in `setupRoutes()` function

## Configuration

All configuration is loaded once in `main.go` and passed through the layers:
```
Config → app.Initialize() → App → server.Config → Server Implementation
```

This ensures consistency across all components.

## Running the Application

### Development
```bash
# Using Gin (current)
go run cmd/api/main.go

# Using Echo (after implementation)
# Just change the import and constructor in main.go
go run cmd/api/main.go
```

### Production
```bash
# Build
go build -o bin/api cmd/api/main.go

# Run
./bin/api
```

## Future Enhancements

1. **gRPC Server**: Implement `internal/server/grpc/grpc_server.go`
2. **WebSocket Server**: Implement `internal/server/websocket/ws_server.go`
3. **Multiple Servers**: Run HTTP and gRPC simultaneously
4. **Hot Reload**: Add configuration-based server selection
5. **Metrics**: Add Prometheus metrics at server layer

## Troubleshooting

### Build Errors
```bash
# Clear module cache
go clean -modcache

# Tidy dependencies
go mod tidy

# Rebuild
go build -o bin/api cmd/api/main.go
```

### Port Already in Use
```bash
# Find process using port 8080
lsof -i :8080

# Kill process
kill -9 <PID>
```

## Summary

The new architecture provides:
- ✅ Clean separation of concerns
- ✅ Framework independence
- ✅ Easy framework migration
- ✅ Better maintainability
- ✅ Improved testability
- ✅ Scalable design

Main changes:
- `main.go`: Reduced from 200+ lines to 36 lines
- `internal/app/`: New package for initialization
- `internal/server/`: New interface-based architecture
- `internal/server/gin/`: Gin-specific implementation

To switch frameworks, simply create a new implementation of the `Server` interface and change one line in `main.go`.
