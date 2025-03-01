package question

import (
	"mucb_be/internal/domain/question"
	"mucb_be/internal/domain/user"
	"mucb_be/internal/errors"
	"mucb_be/internal/infrastructure/security"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type QuestionUseCaseImpl struct {
	questionGroupRepo  question.QuestionGroupRepository
	questionChoiceRepo question.QuestionChoiceRepository
}

func NewAdminUseCase(
	questionGroupRepo question.QuestionGroupRepository,
	questionChoiceRepo question.QuestionChoiceRepository,
) QuestionInterface {
	return &QuestionUseCaseImpl{
		questionGroupRepo:  questionGroupRepo,
		questionChoiceRepo: questionChoiceRepo,
	}
}

func (u *QuestionUseCaseImpl) CreateQuestionGroup(req *CreateQuestionGroupRequest, claims *security.AccessTokenModel) error {
	if claims.Role == user.RoleUser {
		return errors.NewCustomError(
			http.StatusForbidden,
			"UCE004001001",
			"Internal server error.",
			"",
		)
	}

	questionGroup := question.NewQuestionGroup(
		req.ColumnName,
		req.Description,
		req.Limit,
	)

	err := u.questionGroupRepo.CreateQuestionGroup(questionGroup)
	if err != nil {
		return errors.NewCustomError(
			http.StatusForbidden,
			"UCE004001002",
			"Can not create.",
			"",
		)
	}

	return nil
}

func (u *QuestionUseCaseImpl) FindAllQuestionGroup(req *GetQuestionGroupsRequest, claims *security.AccessTokenModel) (*GetQuestionGroupsOutput, error) {
	if claims.Role == user.RoleUser {
		return nil, errors.NewCustomError(
			http.StatusForbidden,
			"UCE004002001",
			"Internal server error.",
			"",
		)
	}

	groups, total, err := u.questionGroupRepo.FindAllQuestionGroup(req.Page, req.Limit)
	if err != nil {
		return nil, errors.NewCustomError(
			http.StatusForbidden,
			"UCE004002002",
			"Internal server error.",
			err.Error(),
		)
	}

	return &GetQuestionGroupsOutput{
		Total: total,
		Page:  req.Page,
		Items: groups,
	}, nil
}

func (u *QuestionUseCaseImpl) CreateQuestionChoice(req *CreateQuestionChoiceRequest, claims *security.AccessTokenModel) error {
	if claims.Role == user.RoleUser {
		return errors.NewCustomError(
			http.StatusForbidden,
			"UCE004003001",
			"Internal server error.",
			"",
		)
	}

	objectId, err := primitive.ObjectIDFromHex(req.QuestionGroup)
	if err != nil {
		return errors.NewCustomError(
			http.StatusForbidden,
			"UCE004003002",
			"Question group not found.",
			err.Error(),
		)
	}

	questionGroup := question.NewQuestionChoice(
		objectId,
		req.Question,
		req.ShouldInvert,
	)

	err = u.questionChoiceRepo.CreateQuestionChoice(questionGroup)
	if err != nil {
		return errors.NewCustomError(
			http.StatusForbidden,
			"UCE004003003",
			"Can not create.",
			"",
		)
	}

	return nil
}

func (u *QuestionUseCaseImpl) FindAllQuestionChoiceByQuestionGroup(req *GetAllQuestionChoiceByQuestionGroupRequest, claims *security.AccessTokenModel) (*GetAllQuestionChoiceByQuestionGroupOutput, error) {
	if claims.Role == user.RoleUser {
		return nil, errors.NewCustomError(
			http.StatusForbidden,
			"UCE004004001",
			"Internal server error.",
			"",
		)
	}

	objectId, err := primitive.ObjectIDFromHex(req.QuestionGroup)
	if err != nil {
		return nil, errors.NewCustomError(
			http.StatusForbidden,
			"UCE004004002",
			"Question group not found.",
			err.Error(),
		)
	}

	questions, err := u.questionChoiceRepo.FindAllQuestionChoiceByQuestionGroup(&objectId)
	if err != nil {
		return nil, errors.NewCustomError(
			http.StatusForbidden,
			"UCE004004002",
			"Question group not found.",
			err.Error(),
		)
	}

	return &GetAllQuestionChoiceByQuestionGroupOutput{
		Items: questions,
	}, nil
}

func (u *QuestionUseCaseImpl) GetQuestionWithRandomChoices(claims *security.AccessTokenModel) (*GetQuestionWithRandomChoicesOutout, error) {
	questionGroup, err := u.questionGroupRepo.FindGroupsWithRandomChoices()
	if err != nil {
		return nil, errors.NewCustomError(
			http.StatusForbidden,
			"UCE004005002",
			"Question group not found.",
			err.Error(),
		)
	}

	return &GetQuestionWithRandomChoicesOutout{
		Items: questionGroup,
	}, nil
}

func (u *QuestionUseCaseImpl) UpdateQuestion(req *UpdateQuestionRequest) error {
	err := u.questionChoiceRepo.UpdateQuestionChoiceById(req.ID, req.Question, req.ShouldInvert)
	if err != nil {
		return errors.NewCustomError(
			http.StatusForbidden,
			"UCE004006001",
			err.Error(),
			err.Error(),
		)
	}

	return nil
}
