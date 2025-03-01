package repository

import (
	"context"
	"mucb_be/internal/domain/auth"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuthRepositoryMongo struct {
	tokenCollection *mongo.Collection
}

func NewAuthRepositoryMongo(tokenCollection *mongo.Collection) auth.AuthRepository {
	return &AuthRepositoryMongo{
		tokenCollection: tokenCollection,
	}
}

func (r *AuthRepositoryMongo) FindTokenById(id string) (*auth.Token, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var result auth.Token
	filter := bson.M{"_id": objectID}
	err = r.tokenCollection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *AuthRepositoryMongo) CreateToken(token auth.Token) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.tokenCollection.InsertOne(ctx, token)
	if err != nil {
		return err
	}

	return nil
}

func (r *AuthRepositoryMongo) RemoveTokenById(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objectID}
	_, err = r.tokenCollection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}
