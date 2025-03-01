package card

import "mucb_be/internal/infrastructure/security"

type CardInterface interface {
	CreateCard(req *CreateCardRequest) error
	FindCardById(id string) (*FindCardOutput, error)
	FindAllCard(req *GetCardsRequest, claims *security.AccessTokenModel) (*GetCardsOutput, error)
	FindCardByIdAndActivate(req *ActivateCard) error
	UpdateCardById(req *UpdateCardRequest) error
}
