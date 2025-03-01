package question

import (
	"mucb_be/internal/infrastructure/security"
)

type QuestionInterface interface {
	CreateQuestionGroup(req *CreateQuestionGroupRequest, claims *security.AccessTokenModel) error
	FindAllQuestionGroup(req *GetQuestionGroupsRequest, claims *security.AccessTokenModel) (*GetQuestionGroupsOutput, error)
	CreateQuestionChoice(req *CreateQuestionChoiceRequest, claims *security.AccessTokenModel) error
	FindAllQuestionChoiceByQuestionGroup(req *GetAllQuestionChoiceByQuestionGroupRequest, claims *security.AccessTokenModel) (*GetAllQuestionChoiceByQuestionGroupOutput, error)
	GetQuestionWithRandomChoices(claims *security.AccessTokenModel) (*GetQuestionWithRandomChoicesOutout, error)
	UpdateQuestion(req *UpdateQuestionRequest) error
}
