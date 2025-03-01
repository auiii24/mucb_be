package record

import (
	"mucb_be/internal/domain/record"
	"mucb_be/internal/errors"
	"mucb_be/internal/infrastructure/security"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RecordUseCaseImpl struct {
	groupRecordRepo record.GroupRecordRepository
	cardRecordRepo  record.CardRecordRepository
	storyRecordRepo record.StoryRecordRepository
}

func NewRecordUseCase(
	groupRecordRepo record.GroupRecordRepository,
	cardRecordRepo record.CardRecordRepository,
	storyRecordRepo record.StoryRecordRepository,
) RecordInterface {
	return &RecordUseCaseImpl{
		groupRecordRepo: groupRecordRepo,
		cardRecordRepo:  cardRecordRepo,
		storyRecordRepo: storyRecordRepo,
	}
}

func (u *RecordUseCaseImpl) CreateManyGroupRecord(req *CreateGroupRecordRequest, claims *security.AccessTokenModel) error {
	var groupRecordList []record.GroupRecord
	var timestamp = time.Now()

	userObjectId, err := primitive.ObjectIDFromHex(claims.ID)
	if err != nil {
		return errors.NewCustomError(
			http.StatusBadRequest,
			"UCE005001001",
			"Failed to check user.",
			err.Error(),
		)
	}

	isAlreadySubmit, err := u.groupRecordRepo.HasSubmittedToday(userObjectId)
	if err != nil {
		return errors.NewCustomError(
			http.StatusBadRequest,
			"UCE005001004",
			"Failed to check submission status.",
			err.Error(),
		)
	}

	if isAlreadySubmit {
		return errors.NewCustomError(
			http.StatusBadRequest,
			"UCE005001005",
			"You can only submit once per day",
			"",
		)
	}

	for _, answer := range req.Answers {
		questionGroupObjectId, err := primitive.ObjectIDFromHex(answer.QuestionGroup)
		if err != nil {
			return errors.NewCustomError(
				http.StatusBadRequest,
				"UCE005001002",
				"Failed to check question group.",
				err.Error(),
			)
		}

		groupRecord := record.NewGroupRecord(req.GroupCode, userObjectId, questionGroupObjectId, answer.Score, answer.QuestionSize, timestamp)
		groupRecordList = append(groupRecordList, *groupRecord)
	}

	err = u.groupRecordRepo.CreateManyGroupRecord((&groupRecordList))
	if err != nil {
		return errors.NewCustomError(
			http.StatusBadRequest,
			"UCE005001003",
			"Failed to insert records.",
			err.Error(),
		)
	}

	return nil
}

func (u *RecordUseCaseImpl) CreateManyCardRecord(req *CreateManyCardRequest, claims *security.AccessTokenModel) error {
	var cardRecordList []record.CardRecord
	var timestamp = time.Now()

	userObjectId, err := primitive.ObjectIDFromHex(claims.ID)
	if err != nil {
		return errors.NewCustomError(
			http.StatusBadRequest,
			"UCE005002001",
			"Failed to check user.",
			err.Error(),
		)
	}

	isAlreadySubmit, err := u.cardRecordRepo.HasSubmittedToday(userObjectId)
	if err != nil {
		return errors.NewCustomError(
			http.StatusBadRequest,
			"UCE005002004",
			"Failed to check submission status.",
			err.Error(),
		)
	}

	if isAlreadySubmit {
		return errors.NewCustomError(
			http.StatusBadRequest,
			"UCE005002005",
			"You can only submit once per day",
			"",
		)
	}

	for _, answer := range req.Answers {
		cardObjectId, err := primitive.ObjectIDFromHex(answer.Card)
		if err != nil {
			return errors.NewCustomError(
				http.StatusBadRequest,
				"UCE005002002",
				"Failed to check card",
				err.Error(),
			)
		}

		cardRecord := record.NewCardRecord(req.GroupCode, userObjectId, cardObjectId, timestamp)
		cardRecordList = append(cardRecordList, *cardRecord)
	}

	err = u.cardRecordRepo.CreateManyGroupRecord((&cardRecordList))
	if err != nil {
		return errors.NewCustomError(
			http.StatusBadRequest,
			"UCE005002003",
			"Failed to insert records.",
			err.Error(),
		)
	}

	return nil

}

func (u *RecordUseCaseImpl) CreateStoryRecord(req *CreateStoryRequest, claims *security.AccessTokenModel) error {
	userObjectId, err := primitive.ObjectIDFromHex(claims.ID)
	if err != nil {
		return errors.NewCustomError(
			http.StatusBadRequest,
			"UCE005003001",
			"Failed to check user.",
			err.Error(),
		)
	}

	newStoryRecord := record.NewStoryRecord(req.GroupCode, userObjectId, req.Content)
	err = u.storyRecordRepo.CreateStoryRecord(newStoryRecord)
	if err != nil {
		return errors.NewCustomError(
			http.StatusBadRequest,
			"UCE005003002",
			"Failed to insert story.",
			err.Error(),
		)
	}

	return nil
}
