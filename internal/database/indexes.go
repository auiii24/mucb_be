package database

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateIndexes ensures all necessary indexes exist in MongoDB
func CreateIndexes(db *mongo.Database) error {
	indexModels := map[string][]mongo.IndexModel{
		AdminsCollection: {
			{Keys: bson.D{{Key: "email", Value: 1}}, Options: options.Index().SetUnique(true)},
		},
		UsersCollection: {
			{Keys: bson.D{{Key: "phone_number", Value: 1}}, Options: options.Index().SetUnique(true)},
		},
		TokensCollection: {
			{Keys: bson.D{{Key: "user", Value: 1}}},
		},
		OtpsCollection: {
			{Keys: bson.D{{Key: "phone_number", Value: 1}}},
			{Keys: bson.D{{Key: "created_at", Value: 1}}},
		},
		OtpAttemptsCollection: {
			{Keys: bson.D{{Key: "phone_number", Value: 1}}},
		},
		QuestionChoicesCollection: {
			{Keys: bson.D{{Key: "question_group", Value: 1}}},
		},
		GroupRecordsCollection: {
			{Keys: bson.D{{Key: "user", Value: 1}}},
			{Keys: bson.D{{Key: "group_code", Value: 1}}},
			{Keys: bson.D{{Key: "created_at", Value: 1}}},
		},
		CardRecordsCollection: {
			{Keys: bson.D{{Key: "user", Value: 1}}},
			{Keys: bson.D{{Key: "group_code", Value: 1}}},
			{Keys: bson.D{{Key: "created_at", Value: 1}}},
		},
		StoryRecordsCollection: {
			{Keys: bson.D{{Key: "user", Value: 1}}},
			{Keys: bson.D{{Key: "group_code", Value: 1}}},
			{Keys: bson.D{{Key: "created_at", Value: 1}}},
		},
		HealthScoresCollection: {
			{Keys: bson.D{{Key: "maximum_percent", Value: 1}}, Options: options.Index().SetUnique(true)},
		},
	}

	// Iterate over collections and create indexes
	for collectionName, indexes := range indexModels {
		collection := db.Collection(collectionName)
		_, err := collection.Indexes().CreateMany(context.Background(), indexes)
		if err != nil {
			log.Printf("Failed to create indexes for %s: %v", collectionName, err)
			return err
		}
		log.Printf("Indexes created for %s", collectionName)
	}

	return nil
}
