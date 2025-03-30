package v1

import (
	"mucb_be/internal/errors"
	"mucb_be/internal/usecase/card"
	"mucb_be/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CardHandler struct {
	cardUseCase card.CardInterface
}

func NewCardHandler(cardUseCase card.CardInterface) *CardHandler {
	return &CardHandler{cardUseCase: cardUseCase}
}

func (h CardHandler) CreateCard(c *gin.Context) {
	var request card.CreateCardRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(
			errors.NewCustomError(http.StatusBadRequest, "VE001001", err.Error(), err.Error()),
		)
		return
	}

	err := h.cardUseCase.CreateCard(&request)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h CardHandler) GetCard(c *gin.Context) {
	cardID := c.Param("cardId")

	response, err := h.cardUseCase.FindCardById(cardID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *CardHandler) GetAllCards(c *gin.Context) {
	claims, err := utils.GetUserClaims(c)
	if err != nil {
		c.Error(err)
		return
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		c.Error(errors.NewCustomError(http.StatusBadRequest, "VE001001", "Invalid page number", ""))
		return
	}

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 || limit > 50 {
		c.Error(errors.NewCustomError(http.StatusBadRequest, "VE001002", "Limit must be between 1 and 50", ""))
		return
	}

	req := card.GetCardsRequest{
		Page:  page,
		Limit: limit,
	}

	response, err := h.cardUseCase.FindAllCard(&req, claims)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h CardHandler) ActivateCard(c *gin.Context) {
	var request card.ActivateCard
	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(
			errors.NewCustomError(http.StatusBadRequest, "VE001001", err.Error(), err.Error()),
		)
		return
	}

	err := h.cardUseCase.FindCardByIdAndActivate(&request)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h CardHandler) UpdateCard(c *gin.Context) {
	var request card.UpdateCardRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(
			errors.NewCustomError(http.StatusBadRequest, "VE001001", err.Error(), err.Error()),
		)
		return
	}

	err := h.cardUseCase.UpdateCardById(&request)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h CardHandler) CheckAvailableCard(c *gin.Context) {
	claims, err := utils.GetUserClaims(c)
	if err != nil {
		c.Error(err)
		return
	}

	err = h.cardUseCase.CheckAvailableCard(claims)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
