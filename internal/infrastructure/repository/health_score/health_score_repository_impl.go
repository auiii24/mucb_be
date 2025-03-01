package repository

import (
	"context"
	"fmt"
	"mucb_be/internal/domain/health_score"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type HealthScoreRepositoryMongo struct {
	healthScoreCollection *mongo.Collection
}

func NewHealthScoreRepositoryMongo(healthScoreCollection *mongo.Collection) health_score.HealthScoreRepository {
	return &HealthScoreRepositoryMongo{
		healthScoreCollection: healthScoreCollection,
	}
}

func (r *HealthScoreRepositoryMongo) CreateHealthScore(healthScore *health_score.HealthScore) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.healthScoreCollection.InsertOne(ctx, healthScore)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf("duplicate maximum percent")
		}
		return err
	}

	return nil
}

func (r *HealthScoreRepositoryMongo) FindAllHealthScore() (*[]health_score.HealthScore, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.healthScoreCollection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	healthScores := make([]health_score.HealthScore, 0)
	err = cursor.All(ctx, &healthScores)
	if err != nil {
		return nil, err
	}

	return &healthScores, nil
}

func (r *HealthScoreRepositoryMongo) FindHealthScoreById(id string) (*health_score.HealthScore, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var healthScore health_score.HealthScore
	err = r.healthScoreCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&healthScore)
	if err != nil {
		return nil, err
	}

	return &healthScore, nil
}

func (r *HealthScoreRepositoryMongo) UpdateHealthScoreById(id string, contents []health_score.HealthScoreContent, maximumPercent int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"contents":        contents,
			"maximum_percent": maximumPercent,
			"updated_at":      time.Now(),
		},
	}

	result, err := r.healthScoreCollection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("data not found")
	}

	if result.ModifiedCount == 0 {
		return fmt.Errorf("can not update")
	}

	return nil
}

func (r *HealthScoreRepositoryMongo) FindContentByScore(score int) (*health_score.HealthScore, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"maximum_percent": bson.M{"$gte": score}}
	findOptions := options.FindOne().SetSort(bson.D{{Key: "maximum_percent", Value: 1}})

	var result health_score.HealthScore
	err := r.healthScoreCollection.FindOne(ctx, filter, findOptions).Decode(&result)

	if err == mongo.ErrNoDocuments {
		highestFilter := bson.M{}
		highestFindOptions := options.FindOne().SetSort(bson.D{{Key: "maximum_percent", Value: -1}})

		err = r.healthScoreCollection.FindOne(ctx, highestFilter, highestFindOptions).Decode(&result)
		if err == mongo.ErrNoDocuments {
			return nil, err
		}
		if err != nil {
			return nil, err
		}
	}

	return &result, nil
}
