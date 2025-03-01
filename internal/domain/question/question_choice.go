package question

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type QuestionChoice struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	QuestionGroup primitive.ObjectID `bson:"question_group" json:"questionGroup"`
	Question      string             `bson:"question" json:"question"`
	ShouldInvert  bool               `bson:"should_invert" json:"shouldInvert"`
	CreatedAt     time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt     time.Time          `bson:"updated_at" json:"updatedAt"`
}

func NewQuestionChoice(questionGroup primitive.ObjectID, question string, shouldInvert bool) *QuestionChoice {
	return &QuestionChoice{
		ID:            primitive.NewObjectID(),
		QuestionGroup: questionGroup,
		Question:      question,
		ShouldInvert:  shouldInvert,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}
