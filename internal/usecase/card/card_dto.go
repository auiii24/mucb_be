package card

import "mucb_be/internal/domain/card"

type CreateCardRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required,max=2048"`
	Image       string `json:"image" binding:"required"`
}

type FindCardOutput struct {
	card.Card
}

type GetCardsRequest struct {
	Page  int `json:"page" binding:"required,min=1"`
	Limit int `json:"limit" binding:"required,min=1,max=50"`
}

type GetCardsOutput struct {
	Total int          `json:"total"`
	Page  int          `json:"page"`
	Items *[]card.Card `json:"items"`
}

type ActivateCard struct {
	Card string `json:"card" binding:"required"`
}

type UpdateCardRequest struct {
	Card        string `json:"card" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description" binding:"required,max=2048"`
	Image       string `json:"image" binding:"required"`
}
