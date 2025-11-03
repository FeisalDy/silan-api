package main

import (
	"log"
	"simple-go/internal/app"
	"simple-go/internal/server"
	ginserver "simple-go/internal/server/gin"
)

func main() {
	// Initialize application
	application, err := app.Initialize()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Create server configuration
	serverConfig := &server.Config{
		Config:                application.Config,
		JWTManager:            application.JWTManager,
		Enforcer:              application.Enforcer,
		AuthHandler:           application.AuthHandler,
		UserHandler:           application.UserHandler,
		NovelHandler:          application.NovelHandler,
		ChapterHandler:        application.ChapterHandler,
		TranslationJobHandler: application.TranslationJobHandler,
		MiscellaneousHandler:  application.MiscellaneousHandler,
		UserService:           application.UserService,
	}

	srv := ginserver.NewGinServer(serverConfig)

	// Start server
	addr := application.Config.Server.Host + ":" + application.Config.Server.Port
	if err := srv.Start(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
