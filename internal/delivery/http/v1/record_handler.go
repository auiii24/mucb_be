package v1

import (
	"mucb_be/internal/errors"
	"mucb_be/internal/usecase/record"
	"mucb_be/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type RecordHandler struct {
	recordUseCase record.RecordInterface
}

func NewRecordHandler(recordUseCase record.RecordInterface) *RecordHandler {
	return &RecordHandler{recordUseCase: recordUseCase}
}

func (h RecordHandler) SubmitGroupAnswer(c *gin.Context) {
	claims, err := utils.GetUserClaims(c)
	if err != nil {
		c.Error(err)
		return
	}

	var request record.CreateGroupRecordRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(
			errors.NewCustomError(http.StatusBadRequest, "VE001001", err.Error(), err.Error()),
		)
		return
	}

	err = h.recordUseCase.CreateManyGroupRecord(&request, claims)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h RecordHandler) SubmitCardAnswer(c *gin.Context) {
	claims, err := utils.GetUserClaims(c)
	if err != nil {
		c.Error(err)
		return
	}

	var request record.CreateManyCardRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(
			errors.NewCustomError(http.StatusBadRequest, "VE001001", err.Error(), err.Error()),
		)
		return
	}

	err = h.recordUseCase.CreateManyCardRecord(&request, claims)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h RecordHandler) SubmitStoryAnswer(c *gin.Context) {
	claims, err := utils.GetUserClaims(c)
	if err != nil {
		c.Error(err)
		return
	}

	var request record.CreateStoryRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(
			errors.NewCustomError(http.StatusBadRequest, "VE001001", err.Error(), err.Error()),
		)
		return
	}

	err = h.recordUseCase.CreateStoryRecord(&request, claims)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
