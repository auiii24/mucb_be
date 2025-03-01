package auth

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OtpAttempt struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	PhoneNumber string             `bson:"phone_number" json:"phone_number"`
	Attempts    int                `bson:"attempts" json:"attempts"`
	LastAttempt time.Time          `bson:"last_attempt" json:"last_attempt"`
}
