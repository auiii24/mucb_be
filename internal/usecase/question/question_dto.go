package question

import "mucb_be/internal/domain/question"

type CreateQuestionGroupRequest struct {
	ColumnName  string `json:"columnName" binding:"required,max=64"`
	Description string `json:"description" binding:"required,max=256"`
	Limit       int    `json:"limit" binding:"required,max=10"`
}

type GetQuestionGroupsRequest struct {
	Page  int `json:"page" binding:"required,min=1"`
	Limit int `json:"limit" binding:"required,min=1,max=50"`
}

type GetQuestionGroupsOutput struct {
	Total int                       `json:"total"`
	Page  int                       `json:"page"`
	Items *[]question.QuestionGroup `json:"items"`
}

type CreateQuestionChoiceRequest struct {
	QuestionGroup string `json:"questionGroup" binding:"required,max=64"`
	Question      string `json:"question" binding:"required,max=256"`
	ShouldInvert  bool   `json:"ShouldInvert"`
}

type GetAllQuestionChoiceByQuestionGroupOutput struct {
	Items *[]question.QuestionChoice `json:"items"`
}

type GetQuestionWithRandomChoicesOutout struct {
	Items *[]question.GroupsWithRandomChoices `json:"items"`
}

type UpdateQuestionRequest struct {
	ID           string `json:"id"`
	Question     string `json:"question" binding:"required,max=256"`
	ShouldInvert bool   `json:"ShouldInvert"`
}

type RemoveChoiceRequest struct {
	ID string `json:"id"`
}

type RemoveQuestionGroupRequest struct {
	ID string `json:"id"`
}

type UpdateQuestionGroupRequest struct {
	ID          string `json:"id"`
	ColumnName  string `json:"columnName" binding:"required,max=64"`
	Description string `json:"description" binding:"required,max=256"`
	Limit       int    `json:"limit" binding:"required,max=10"`
}
