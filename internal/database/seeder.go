package database

import (
	"context"
	"log"
	"time"

	"mucb_be/internal/domain/admin"
	"mucb_be/internal/infrastructure/security"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func SeedAdmin(db *mongo.Database, hashService security.HashServiceInterface) {
	collection := db.Collection(AdminsCollection)

	count, err := collection.CountDocuments(context.Background(), bson.M{})
	if err != nil {
		log.Fatalf("Error checking admin collection: %v", err)
	}

	if count > 0 {
		log.Println("Admin already exists. Skipping seeding.")
		return
	}

	hashedPassword, _ := hashService.HashPassword("admin1234#")
	superAdmin := admin.Admin{
		Name:      "Super Admin",
		Email:     "admin@mucb.com",
		Password:  hashedPassword,
		Role:      admin.RoleSuperAdmin,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = collection.InsertOne(context.Background(), superAdmin)
	if err != nil {
		log.Fatalf("Error inserting SUPER_ADMIN: %v", err)
	}

	log.Println("SUPER_ADMIN seeded successfully!")
}
