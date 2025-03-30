package health_score

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	ContentTypeTitle       = "TITLE"
	ContentTypeImage       = "IMAGE"
	ContentTypeDescription = "DESCRIPTION"
	ContentTypeSpace       = "SPACE"
	ContentTypeURL         = "URL"
)

type HealthScoreContent struct {
	ContentType string `bson:"content_type" json:"contentType" binding:"required"`
	Content     string `bson:"content" json:"content" binding:"required"`
}

type HealthScore struct {
	ID             primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Contents       []HealthScoreContent `bson:"contents" json:"contents"`
	MaximumPercent int                  `bson:"maximum_percent" json:"maximumPercent"`
	CreatedAt      time.Time            `bson:"created_at" json:"createdAt"`
	UpdatedAt      time.Time            `bson:"updated_at" json:"updatedAt"`
}

func NewHealthScore(contents []HealthScoreContent, maximumPercent int) *HealthScore {
	return &HealthScore{
		ID:             primitive.NewObjectID(),
		Contents:       contents,
		MaximumPercent: maximumPercent,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}
