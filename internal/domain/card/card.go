package card

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Card struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Image       primitive.ObjectID `bson:"image" json:"image"`
	Description string             `bson:"description" json:"description"`
	IsActive    bool               `bson:"is_active" json:"isActive"`
	CreatedAt   time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updatedAt"`
}

func NewCard(name, description string, image primitive.ObjectID) *Card {
	return &Card{
		ID:          primitive.NewObjectID(),
		Name:        name,
		Image:       image,
		Description: description,
		IsActive:    false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}
