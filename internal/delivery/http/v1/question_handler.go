package v1

import (
	"mucb_be/internal/errors"
	"mucb_be/internal/usecase/question"
	"mucb_be/internal/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type QuestionHandler struct {
	questionUseCase question.QuestionInterface
}

func NewQuestionHandler(questionUseCase question.QuestionInterface) *QuestionHandler {
	return &QuestionHandler{questionUseCase: questionUseCase}
}

func (h QuestionHandler) CreateQuestionGroup(c *gin.Context) {
	claims, err := utils.GetUserClaims(c)
	if err != nil {
		c.Error(err)
		return
	}

	var request question.CreateQuestionGroupRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(
			errors.NewCustomError(http.StatusBadRequest, "VE001001", err.Error(), err.Error()),
		)
		return
	}

	err = h.questionUseCase.CreateQuestionGroup(&request, claims)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// ✅ Get all question groups with pagination
func (h *QuestionHandler) GetAllQuestionGroups(c *gin.Context) {
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

	req := question.GetQuestionGroupsRequest{
		Page:  page,
		Limit: limit,
	}

	response, err := h.questionUseCase.FindAllQuestionGroup(&req, claims)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h QuestionHandler) CreateQuestionChoice(c *gin.Context) {
	claims, err := utils.GetUserClaims(c)
	if err != nil {
		c.Error(err)
		return
	}

	var request question.CreateQuestionChoiceRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(
			errors.NewCustomError(http.StatusBadRequest, "VE001001", err.Error(), err.Error()),
		)
		return
	}

	err = h.questionUseCase.CreateQuestionChoice(&request, claims)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h QuestionHandler) FindAllQuestionChoiceByQuestionGroup(c *gin.Context) {
	idStr := c.Param("id")

	claims, err := utils.GetUserClaims(c)
	if err != nil {
		c.Error(err)
		return
	}

	response, err := h.questionUseCase.FindAllQuestionChoiceByQuestionGroup(idStr, claims)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h QuestionHandler) GetQuestionWithRandomChoices(c *gin.Context) {
	claims, err := utils.GetUserClaims(c)
	if err != nil {
		c.Error(err)
		return
	}

	response, err := h.questionUseCase.GetQuestionWithRandomChoices(claims)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h QuestionHandler) UpdateQuestion(c *gin.Context) {
	var request question.UpdateQuestionRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(
			errors.NewCustomError(http.StatusBadRequest, "VE001001", err.Error(), err.Error()),
		)
		return
	}

	err := h.questionUseCase.UpdateQuestion(&request)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h QuestionHandler) CheckAvailableQuestion(c *gin.Context) {
	claims, err := utils.GetUserClaims(c)
	if err != nil {
		c.Error(err)
		return
	}

	err = h.questionUseCase.CheckAvailableQuestion(claims)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h QuestionHandler) RemoveChoice(c *gin.Context) {
	var request question.RemoveChoiceRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(
			errors.NewCustomError(http.StatusBadRequest, "VE001001", err.Error(), err.Error()),
		)
		return
	}

	err := h.questionUseCase.RemoveChoice(&request)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h QuestionHandler) RemoveQuestionGroup(c *gin.Context) {
	var request question.RemoveQuestionGroupRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(
			errors.NewCustomError(http.StatusBadRequest, "VE001001", err.Error(), err.Error()),
		)
		return
	}

	err := h.questionUseCase.RemoveQuestionGroup(&request)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h QuestionHandler) GetQuestionGroupById(c *gin.Context) {
	idStr := c.Param("id")

	response, err := h.questionUseCase.GetQuestionGroupById(idStr)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h QuestionHandler) UpdateQuestionGroup(c *gin.Context) {
	var request question.UpdateQuestionGroupRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(
			errors.NewCustomError(http.StatusBadRequest, "VE001001", err.Error(), err.Error()),
		)
		return
	}

	err := h.questionUseCase.UpdateQuestionGroup(&request)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
