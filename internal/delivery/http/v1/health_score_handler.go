package v1

import (
	"mucb_be/internal/errors"
	"mucb_be/internal/usecase/health_score"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HealthScoreHandler struct {
	healthScoreUseCase health_score.HealthScoreInterface
}

func NewHealthScoreHandler(healthScoreUseCase health_score.HealthScoreInterface) *HealthScoreHandler {
	return &HealthScoreHandler{
		healthScoreUseCase: healthScoreUseCase,
	}
}

func (h HealthScoreHandler) CreateHealthScore(c *gin.Context) {
	var request health_score.CreateHealthScore
	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(
			errors.NewCustomError(http.StatusBadRequest, "VE001001", err.Error(), err.Error()),
		)
		return
	}

	err := h.healthScoreUseCase.CreateHealthScore(&request)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h HealthScoreHandler) GetAllHealthScore(c *gin.Context) {
	response, err := h.healthScoreUseCase.FindAllHealthScore()
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h HealthScoreHandler) GetHealthScore(c *gin.Context) {
	healthScoreID := c.Param("healthScoreId")

	response, err := h.healthScoreUseCase.FindHealthScoreById(healthScoreID)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h HealthScoreHandler) UpdateHealthScoreById(c *gin.Context) {
	var request health_score.UpdateHealthScoreByIdRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(
			errors.NewCustomError(http.StatusBadRequest, "VE001001", err.Error(), err.Error()),
		)
		return
	}

	err := h.healthScoreUseCase.UpdateHealthScoreById(&request)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h HealthScoreHandler) GetContentByScore(c *gin.Context) {
	var request health_score.GetContentByScoreRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(
			errors.NewCustomError(http.StatusBadRequest, "VE001001", err.Error(), err.Error()),
		)
		return
	}

	response, err := h.healthScoreUseCase.FindContentByScore(&request)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response)
}
