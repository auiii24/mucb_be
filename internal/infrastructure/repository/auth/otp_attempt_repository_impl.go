package repository

import (
	"context"
	"errors"
	"mucb_be/internal/domain/auth"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OtpAttemptRepositoryMongo struct {
	otpAttemptCollection *mongo.Collection
}

func NewOtpAttemptRepositoryMongo(otpAttemptCollection *mongo.Collection) auth.OtpAttemptRepository {
	return &OtpAttemptRepositoryMongo{
		otpAttemptCollection: otpAttemptCollection,
	}
}

func (r *OtpAttemptRepositoryMongo) CheckOtpRateLimit(phoneNumber string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var attempt auth.OtpAttempt
	filter := bson.M{"phone_number": phoneNumber}
	err := r.otpAttemptCollection.FindOne(ctx, filter).Decode(&attempt)
	if err == mongo.ErrNoDocuments {
		// ✅ No previous attempts, allow OTP request
		return false, nil
	} else if err != nil {
		return false, err
	}

	// ✅ Step 1: Prevent multiple OTPs within 1 minute
	if time.Since(attempt.LastAttempt) < time.Minute {
		return true, errors.New("OTP request too frequent. Please wait 1 minute.")
	}

	// ✅ Step 2: Block user after 5 attempts in 1 hour
	if attempt.Attempts >= 5 && time.Since(attempt.LastAttempt) < time.Hour {
		return true, errors.New("Too many OTP requests. Please try again in 1 hour.")
	}

	return false, nil
}

// ✅ Step 3: Increment OTP request count
func (r *OtpAttemptRepositoryMongo) IncrementOtpAttempt(phoneNumber string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"phone_number": phoneNumber}
	update := bson.M{
		"$inc": bson.M{"attempts": 1},
		"$set": bson.M{"last_attempt": time.Now()},
	}

	_, err := r.otpAttemptCollection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	return err
}

// ✅ Clear OTP attempts after successful login
func (r *OtpAttemptRepositoryMongo) ClearOtpAttempts(phoneNumber string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.otpAttemptCollection.DeleteOne(ctx, bson.M{"phone_number": phoneNumber})
	return err
}
