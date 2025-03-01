package record

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CardRecord struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	User      primitive.ObjectID `bson:"user" json:"user"`
	Card      primitive.ObjectID `bson:"card" json:"card"`
	GroupCode *string            `bson:"group_code" json:"groupCode"`
	CreatedAt time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updatedAt"`
}

func NewCardRecord(groupCode *string, user, card primitive.ObjectID, timestamp time.Time) *CardRecord {
	return &CardRecord{
		ID:        primitive.NewObjectID(),
		User:      user,
		Card:      card,
		GroupCode: groupCode,
		CreatedAt: timestamp,
		UpdatedAt: timestamp,
	}
}
