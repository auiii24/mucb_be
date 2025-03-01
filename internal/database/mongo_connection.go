package database

import (
	"context"
	"fmt"
	"log"
	"mucb_be/internal/config"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMongoDB(cfg *config.Config) (*mongo.Client, error) {
	mongoConnection := fmt.Sprintf(
		"mongodb://%s:%s@%s/%s",
		cfg.DatabaseUsername,
		cfg.DatabasePassword,
		cfg.DatabaseUrl,
		cfg.DatabaseName)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(mongoConnection)

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping the database to verify the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	db := client.Database(cfg.DatabaseName)

	// ✅ Ensure indexes exist before returning the database
	if err := CreateIndexes(db); err != nil {
		log.Fatal("❌ Error creating indexes:", err)
	}

	fmt.Println("Successfully connected and pinged MongoDB!")

	return client, nil
}
