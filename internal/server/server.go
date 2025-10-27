package server

import (
	"simple-go/internal/handler"
	"simple-go/internal/service"
	"simple-go/pkg/auth"
	"simple-go/pkg/config"

	"github.com/casbin/casbin/v2"
)

// Server represents an HTTP server
type Server interface {
	Start(addr string) error
	Shutdown() error
}

// Config holds all the dependencies needed to create a server
type Config struct {
	Config         *config.Config
	JWTManager     *auth.JWTManager
	Enforcer       *casbin.Enforcer
	AuthHandler    *handler.AuthHandler
	UserHandler    *handler.UserHandler
	NovelHandler   *handler.NovelHandler
	ChapterHandler *handler.ChapterHandler
	UserService    *service.UserService
}
