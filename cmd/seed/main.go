package main

import (
	"log"
	"os"
	"simple-go/pkg/config"
	"simple-go/pkg/database"
	"simple-go/seeds"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to database
	db, err := database.Connect(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run seeders
	if err := seeds.RunAllSeeds(db); err != nil {
		log.Fatalf("Failed to run seeders: %v", err)
		os.Exit(1)
	}

	log.Println("✅ Database seeding completed successfully!")
}
