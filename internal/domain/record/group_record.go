package record

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GroupRecord struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	User          primitive.ObjectID `bson:"user" json:"user"`
	QuestionGroup primitive.ObjectID `bson:"question_group" json:"questionGroup"`
	Score         int                `bson:"score" json:"score"`
	Size          int                `bson:"question_size" json:"questionSize"`
	GroupCode     *string            `bson:"group_code" json:"groupCode"`
	CreatedAt     time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updatedAt"`
}

func NewGroupRecord(groupCode *string, user, questionGroup primitive.ObjectID, score, size int, timestamp time.Time) *GroupRecord {
	return &GroupRecord{
		ID:            primitive.NewObjectID(),
		User:          user,
		QuestionGroup: questionGroup,
		Score:         score,
		Size:          size,
		GroupCode:     groupCode,
		CreatedAt:     timestamp,
		UpdatedAt:     timestamp,
	}
}
