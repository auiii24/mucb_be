package card

import (
	"mucb_be/internal/domain/card"
	"mucb_be/internal/domain/image"
	"mucb_be/internal/domain/record"
	"mucb_be/internal/domain/user"
	"mucb_be/internal/errors"
	"mucb_be/internal/infrastructure/security"

	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CardUserCaseImpl struct {
	cardRepo       card.CardRepository
	imageRepo      image.ImageRepository
	cardRecordRepo record.CardRecordRepository
}

func NewCardUseCase(
	cardRepo card.CardRepository,
	imageRepo image.ImageRepository,
	cardRecordRepo record.CardRecordRepository,
) CardInterface {
	return &CardUserCaseImpl{
		cardRepo:       cardRepo,
		imageRepo:      imageRepo,
		cardRecordRepo: cardRecordRepo,
	}
}

func (u *CardUserCaseImpl) CreateCard(req *CreateCardRequest) error {

	imageObjectID, err := primitive.ObjectIDFromHex(req.Image)
	if err != nil {
		return errors.NewCustomError(
			http.StatusBadRequest,
			"UCE007001001",
			"Failed to convert image.",
			err.Error(),
		)
	}

	newCard := card.NewCard(req.Name, req.Description, imageObjectID)
	err = u.cardRepo.CreateCard(newCard)
	if err != nil {
		return errors.NewCustomError(
			http.StatusBadRequest,
			"UCE007001002",
			"Failed to insert card.",
			err.Error(),
		)
	}

	return nil
}

func (u *CardUserCaseImpl) FindCardById(id string) (*FindCardOutput, error) {
	cardExist, err := u.cardRepo.FindCardById(id)
	if err != nil {
		return nil, errors.NewCustomError(
			http.StatusBadRequest,
			"UCE007002001",
			"Failed to find card.",
			err.Error(),
		)
	}

	output := FindCardOutput{
		Card: *cardExist,
	}

	return &output, nil
}

func (u *CardUserCaseImpl) FindAllCard(req *GetCardsRequest, claims *security.AccessTokenModel) (*GetCardsOutput, error) {
	var isAdmin bool = true
	if claims.Role == user.RoleUser {
		isAdmin = false
	}

	if !isAdmin {
		err := u.CheckAvailableCard(claims)
		if err != nil {
			return nil, err
		}
	}

	groups, total, err := u.cardRepo.FindAllCardByRole(req.Page, req.Limit, isAdmin)
	if err != nil {
		return nil, errors.NewCustomError(
			http.StatusForbidden,
			"UCE007003002",
			"Failed to find cards.",
			err.Error(),
		)
	}

	return &GetCardsOutput{
		Total: total,
		Page:  req.Page,
		Items: groups,
	}, nil
}

func (u *CardUserCaseImpl) FindCardByIdAndActivate(req *ActivateCard) error {
	err := u.cardRepo.FindCardByIdAndActivate(req.Card)
	if err != nil {
		return errors.NewCustomError(
			http.StatusForbidden,
			"UCE007004001",
			"Failed to update cards.",
			err.Error(),
		)
	}

	return nil
}

func (u *CardUserCaseImpl) UpdateCardById(req *UpdateCardRequest) error {
	cardExist, err := u.cardRepo.FindCardById(req.Card)
	if err != nil {
		return errors.NewCustomError(
			http.StatusBadRequest,
			"UCE007005003",
			"Failed to find card.",
			err.Error(),
		)
	}

	imageObjectID, err := primitive.ObjectIDFromHex(req.Image)
	if err != nil {
		return errors.NewCustomError(
			http.StatusBadRequest,
			"UCE007005001",
			"Failed to convert image.",
			err.Error(),
		)
	}

	err = u.cardRepo.UpdateCardById(req.Card, req.Name, req.Description, imageObjectID)
	if err != nil {
		return errors.NewCustomError(
			http.StatusBadRequest,
			"UCE007005002",
			"Failed to update card.",
			err.Error(),
		)
	}

	_ = u.imageRepo.UpdateImageStatusById(cardExist.Image.Hex(), false)
	_ = u.imageRepo.UpdateImageStatusById(imageObjectID.Hex(), true)

	return nil
}

func (u *CardUserCaseImpl) CheckAvailableCard(claims *security.AccessTokenModel) error {
	userObjectId, err := primitive.ObjectIDFromHex(claims.ID)
	if err != nil {
		return errors.NewCustomError(
			http.StatusBadRequest,
			"UCE007006001",
			"Failed to check user.",
			err.Error(),
		)
	}

	isAlreadySubmit, err := u.cardRecordRepo.HasSubmittedToday(userObjectId)
	if err != nil {
		return errors.NewCustomError(
			http.StatusBadRequest,
			"UCE007006002",
			"Failed to check submission status.",
			err.Error(),
		)
	}

	if isAlreadySubmit {
		return errors.NewCustomError(
			http.StatusBadRequest,
			"UCE007006003",
			"You can only submit once per day",
			"",
		)
	}

	return nil
}
