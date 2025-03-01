package repository

import (
	"context"
	"mucb_be/internal/domain/user"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepositoryMongo struct {
	userCollection *mongo.Collection
}

func NewUserRepositoryMongo(userCollection *mongo.Collection) user.UserRepository {
	return &UserRepositoryMongo{
		userCollection: userCollection,
	}
}

func (r *UserRepositoryMongo) CreateUser(user *user.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.userCollection.InsertOne(ctx, user)
	return err
}

func (r *UserRepositoryMongo) FindUserByPhoneNumber(phoneNumber string) (*user.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result user.User
	err := r.userCollection.FindOne(ctx, bson.M{"phone_number": phoneNumber}).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *UserRepositoryMongo) FindUserById(id string) (*user.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var result user.User
	err = r.userCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *UserRepositoryMongo) UpdateUserInfo(id, name, group string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"name":       name,
			"group_code": group,
			"state":      user.UserStateActive,
			"updated_at": time.Now(),
		},
	}

	_, err = r.userCollection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return err
	}

	return nil
}
