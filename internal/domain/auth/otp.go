package auth

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Otp struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	User         primitive.ObjectID `bson:"user" json:"user"`
	PhoneNumber  string             `bson:"phone_number" json:"phoneNumber"`
	RefCode      string             `bson:"ref_code" json:"refCode"`
	Code         string             `bson:"code" json:"code"`
	IsUsed       bool               `bson:"is_used" json:"isUsed"`
	AttemptCount int                `bson:"attempt_count" json:"attemptCount"`
	CreatedAt    time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updatedAt"`
}

func NewOtp(user primitive.ObjectID, phoneNumber, refCode, code string) *Otp {
	return &Otp{
		ID:           primitive.NewObjectID(),
		User:         user,
		PhoneNumber:  phoneNumber,
		RefCode:      refCode,
		Code:         code,
		IsUsed:       false,
		AttemptCount: 0,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}
