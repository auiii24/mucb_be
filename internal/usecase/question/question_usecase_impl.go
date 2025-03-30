package question

import (
	"mucb_be/internal/domain/question"
	"mucb_be/internal/domain/record"
	"mucb_be/internal/domain/user"
	"mucb_be/internal/errors"
	"mucb_be/internal/infrastructure/security"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type QuestionUseCaseImpl struct {
	questionGroupRepo  question.QuestionGroupRepository
	questionChoiceRepo question.QuestionChoiceRepository
	groupRecordRepo    record.GroupRecordRepository
}

func NewAdminUseCase(
	questionGroupRepo question.QuestionGroupRepository,
	questionChoiceRepo question.QuestionChoiceRepository,
	groupRecordRepo record.GroupRecordRepository,
) QuestionInterface {
	return &QuestionUseCaseImpl{
		questionGroupRepo:  questionGroupRepo,
		questionChoiceRepo: questionChoiceRepo,
		groupRecordRepo:    groupRecordRepo,
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

func (u *QuestionUseCaseImpl) FindAllQuestionChoiceByQuestionGroup(idStr string, claims *security.AccessTokenModel) (*GetAllQuestionChoiceByQuestionGroupOutput, error) {
	if claims.Role == user.RoleUser {
		return nil, errors.NewCustomError(
			http.StatusForbidden,
			"UCE004004001",
			"Internal server error.",
			"",
		)
	}

	objectId, err := primitive.ObjectIDFromHex(idStr)
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
	err := u.CheckAvailableQuestion(claims)
	if err != nil {
		return nil, err
	}

	questionGroup, err := u.questionGroupRepo.FindGroupsWithRandomChoices()
	if err != nil {
		return nil, errors.NewCustomError(
			http.StatusBadRequest,
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
			http.StatusBadRequest,
			"UCE004006001",
			err.Error(),
			err.Error(),
		)
	}

	return nil
}

func (u *QuestionUseCaseImpl) CheckAvailableQuestion(claims *security.AccessTokenModel) error {
	userObjectId, err := primitive.ObjectIDFromHex(claims.ID)
	if err != nil {
		return errors.NewCustomError(
			http.StatusBadRequest,
			"UCE004007001",
			"Failed to check user.",
			err.Error(),
		)
	}

	isAlreadySubmit, err := u.groupRecordRepo.HasSubmittedToday(userObjectId)
	if err != nil {
		return errors.NewCustomError(
			http.StatusBadRequest,
			"UCE004007002",
			"Failed to check submission status.",
			err.Error(),
		)
	}

	if isAlreadySubmit {
		return errors.NewCustomError(
			http.StatusBadRequest,
			"UCE004007003",
			"You can only submit once per day",
			"",
		)
	}

	return nil
}

func (u *QuestionUseCaseImpl) RemoveChoice(req *RemoveChoiceRequest) error {
	err := u.questionChoiceRepo.RemoveChoiceById(req.ID)
	if err != nil {
		return errors.NewCustomError(
			http.StatusBadRequest,
			"UCE004008001",
			"Failed to remove choice.",
			err.Error(),
		)
	}

	return nil
}

func (u *QuestionUseCaseImpl) RemoveQuestionGroup(req *RemoveQuestionGroupRequest) error {
	err := u.questionGroupRepo.RemoveQuestionGroupById(req.ID)
	if err != nil {
		return errors.NewCustomError(
			http.StatusBadRequest,
			"UCE004009001",
			"Failed to remove question group.",
			err.Error(),
		)
	}

	u.questionChoiceRepo.RemoveChoicesByQuestionGroupId(req.ID)

	return nil
}

func (u *QuestionUseCaseImpl) GetQuestionGroupById(id string) (*question.QuestionGroup, error) {
	existQuestionGroup, err := u.questionGroupRepo.FindQuestionGroupById(id)
	if err != nil {
		return nil, errors.NewCustomError(
			http.StatusBadRequest,
			"UCE004010001",
			"Failed to get question group.",
			err.Error(),
		)
	}

	return existQuestionGroup, nil
}

func (u *QuestionUseCaseImpl) UpdateQuestionGroup(req *UpdateQuestionGroupRequest) error {
	err := u.questionGroupRepo.UpdateQuestionGroupById(req.ID, req.ColumnName, req.Description, req.Limit)
	if err != nil {
		return errors.NewCustomError(
			http.StatusBadRequest,
			"UCE004011001",
			"Failed to update question group.",
			err.Error(),
		)
	}

	return nil
}
