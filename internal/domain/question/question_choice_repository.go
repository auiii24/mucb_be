package question

import "go.mongodb.org/mongo-driver/bson/primitive"

type QuestionChoiceRepository interface {
	CreateQuestionChoice(questionChoice *QuestionChoice) error
	FindAllQuestionChoiceByQuestionGroup(questionGroup *primitive.ObjectID) (*[]QuestionChoice, error)
	UpdateQuestionChoiceById(id, question string, shouldInvert bool) error
}
