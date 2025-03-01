package auth

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Token struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	User      primitive.ObjectID `bson:"user" json:"user"`
	UniqueID  string             `bson:"unique_id" json:"uniqueId"`
	Info      string             `bson:"info" json:"info"`
	UserType  string             `bson:"user_type" json:"userType"`
	CreatedAt time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updatedAt"`
}

func NewToken(user primitive.ObjectID, uniqueId, info, userType string) *Token {
	return &Token{
		ID:        primitive.NewObjectID(),
		User:      user,
		UniqueID:  uniqueId,
		Info:      info,
		UserType:  userType,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
