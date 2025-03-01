package user

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	UserStatePending   = "PENDING"
	UserStateActive    = "ACTIVE"
	UserStateSuspended = "SUSPENDED"
)

type User struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	PhoneNumber string             `bson:"phone_number" json:"phoneNumber"`
	State       string             `bson:"state" json:"state"`
	GroupCode   *string            `bson:"group_code" json:"group"`
	CreatedAt   time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updatedAt"`
}

func NewUser(phoneNumber string) *User {
	return &User{
		ID:          primitive.NewObjectID(),
		PhoneNumber: phoneNumber,
		State:       UserStatePending,
		GroupCode:   nil,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}
