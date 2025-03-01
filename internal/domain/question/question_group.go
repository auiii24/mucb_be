package question

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type QuestionGroup struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ColumnName  string             `bson:"column_name" json:"columnName"`
	Description string             `bson:"description" json:"description"`
	Limit       int                `bson:"limit" json:"limit"`
	CreatedAt   time.Time          `bson:"created_at" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updatedAt"`
}

type GroupsWithRandomChoices struct {
	QuestionGroup `bson:",inline"`
	Choices       []QuestionChoice `bson:"choices" json:"choices"`
}

func NewQuestionGroup(columnName, description string, limit int) *QuestionGroup {
	return &QuestionGroup{
		ID:          primitive.NewObjectID(),
		ColumnName:  columnName,
		Description: description,
		Limit:       limit,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}
