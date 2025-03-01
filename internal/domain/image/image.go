package image

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Image struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Path         string             `bson:"path" json:"path"`
	Name         string             `bson:"name" json:"name"`
	OriginalName string             `bson:"original_name" json:"originalName"`
	Width        int                `bson:"width" json:"width"`
	Height       int                `bson:"height" json:"height"`
	ContentType  string             `bson:"content_type" json:"contentType"`
	IsActive     bool               `bson:"is_active" json:"isActive"`
	CreatedAt    time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt    time.Time          `bson:"updated_at" json:"updatedAt"`
}

func NewImage(path, name, originalName, contentType string, width, height int) *Image {
	return &Image{
		ID:           primitive.NewObjectID(),
		Path:         path,
		Name:         name,
		OriginalName: originalName,
		Width:        width,
		Height:       height,
		ContentType:  contentType,
		IsActive:     false,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}
