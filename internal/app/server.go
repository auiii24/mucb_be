package app

import (
	"context"
	"log"
	"mucb_be/internal/config"
	"mucb_be/internal/database"
)

func Start() (*Dependencies, *config.Config) {
	cfg, err := config.LoadConfig()

	if err != nil {
		log.Fatalf("can not load config %v", err)
	}

	dbClient, err := database.ConnectMongoDB(cfg)
	if err != nil {
		log.Fatalf("can not connect mongodb %v", err)
	}

	deps := NewDependencies(cfg, dbClient)

	database.SeedAdmin(dbClient.Database(cfg.DatabaseName), deps.HashService)

	return deps, cfg
}

func Stop(ctx context.Context, deps *Dependencies) {
	log.Println("Closing MongoDB connection...")
	if err := deps.DBClient.Disconnect(ctx); err != nil {
		log.Printf("Error disconnecting MongoDB: %v", err)
	} else {
		log.Println("MongoDB disconnected successfully")
	}
}
