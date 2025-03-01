package card

import "go.mongodb.org/mongo-driver/bson/primitive"

type CardRepository interface {
	CreateCard(card *Card) error
	FindCardById(id string) (*Card, error)
	FindAllCardByRole(page, limit int, isAdmin bool) (*[]Card, int, error)
	FindCardByIdAndActivate(id string) error
	UpdateCardById(id, name, description string, image primitive.ObjectID) error
}
