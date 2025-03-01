package repository

import (
	"context"
	"mucb_be/internal/domain/auth"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OtpRepositoryMongo struct {
	otpCollection *mongo.Collection
}

func NewOtpRepositoryMongo(otpCollection *mongo.Collection) auth.OtpRepository {
	return &OtpRepositoryMongo{
		otpCollection: otpCollection,
	}
}

func (r *OtpRepositoryMongo) CreateOtp(otp *auth.Otp) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.otpCollection.InsertOne(ctx, otp)
	return err
}

func (r *OtpRepositoryMongo) FindLatestOtpByPhoneNumber(phoneNumber string) (*auth.Otp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var result auth.Otp
	filter := bson.M{"phone_number": phoneNumber}
	opts := options.FindOne().SetSort(bson.D{{Key: "created_at", Value: -1}})

	err := r.otpCollection.FindOne(ctx, filter, opts).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *OtpRepositoryMongo) MarkOtpAsUsedById(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"is_used":    true,
			"updated_at": time.Now(),
		},
	}

	_, err = r.otpCollection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return err
	}

	return nil
}

func (r *OtpRepositoryMongo) IncrementOtpAttemptsById(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$inc": bson.M{"attempt_count": 1},
		"$set": bson.M{"updated_at": time.Now()},
	}

	_, err = r.otpCollection.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	if err != nil {
		return err
	}

	return nil
}
