package question

import (
	"mucb_be/internal/domain/question"
	"mucb_be/internal/infrastructure/security"
)

type QuestionInterface interface {
	CreateQuestionGroup(req *CreateQuestionGroupRequest, claims *security.AccessTokenModel) error
	FindAllQuestionGroup(req *GetQuestionGroupsRequest, claims *security.AccessTokenModel) (*GetQuestionGroupsOutput, error)
	CreateQuestionChoice(req *CreateQuestionChoiceRequest, claims *security.AccessTokenModel) error
	FindAllQuestionChoiceByQuestionGroup(id string, claims *security.AccessTokenModel) (*GetAllQuestionChoiceByQuestionGroupOutput, error)
	GetQuestionWithRandomChoices(claims *security.AccessTokenModel) (*GetQuestionWithRandomChoicesOutout, error)
	UpdateQuestion(req *UpdateQuestionRequest) error
	CheckAvailableQuestion(claims *security.AccessTokenModel) error
	RemoveChoice(req *RemoveChoiceRequest) error
	RemoveQuestionGroup(req *RemoveQuestionGroupRequest) error
	GetQuestionGroupById(id string) (*question.QuestionGroup, error)
	UpdateQuestionGroup(req *UpdateQuestionGroupRequest) error
}
