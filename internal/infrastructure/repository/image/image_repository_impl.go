package repository

import (
	"context"
	"fmt"
	"mucb_be/internal/domain/image"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ImageRepositoryMongo struct {
	imageCollection *mongo.Collection
}

func NewImageRepositoryMongo(imageCollection *mongo.Collection) image.ImageRepository {
	return &ImageRepositoryMongo{
		imageCollection: imageCollection,
	}
}

func (r *ImageRepositoryMongo) CreateImage(image *image.Image) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.imageCollection.InsertOne(ctx, image)
	return err
}

func (r *ImageRepositoryMongo) FindImageByID(id string) (*image.Image, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var result image.Image
	err = r.imageCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *ImageRepositoryMongo) UpdateImageStatusById(id string, currentStats bool) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"is_active":  currentStats,
			"updated_at": time.Now(),
		},
	}

	result, err := r.imageCollection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
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
