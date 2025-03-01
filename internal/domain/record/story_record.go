package record

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StoryRecord struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	User      primitive.ObjectID `bson:"user" json:"user"`
	Content   string             `bson:"content" json:"content"`
	GroupCode *string            `bson:"group_code" json:"groupCode"`
	CreatedAt time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updatedAt"`
}

func NewStoryRecord(groupCode *string, user primitive.ObjectID, content string) *StoryRecord {
	return &StoryRecord{
		ID:        primitive.NewObjectID(),
		User:      user,
		Content:   content,
		GroupCode: groupCode,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
