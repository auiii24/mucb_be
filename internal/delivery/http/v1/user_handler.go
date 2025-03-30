package v1

import (
	"mucb_be/internal/errors"
	"mucb_be/internal/usecase/user"
	"mucb_be/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUseCase user.UserUseCaseInterface
}

func NewUserHandler(userUseCase user.UserUseCaseInterface) *UserHandler {
	return &UserHandler{userUseCase: userUseCase}
}

func (h UserHandler) UpdateUserInfo(c *gin.Context) {
	claims, err := utils.GetUserClaims(c)
	if err != nil {
		c.Error(err)
		return
	}

	var request user.UpdateUserInfoRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(
			errors.NewCustomError(http.StatusBadRequest, "VE001001", err.Error(), err.Error()),
		)
		return
	}

	response, err := h.userUseCase.UpdateUserInfo(&request, claims)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h UserHandler) GetUserInfo(c *gin.Context) {
	claims, err := utils.GetUserClaims(c)
	if err != nil {
		c.Error(err)
		return
	}
	response, err := h.userUseCase.GetUserInfo(claims)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h UserHandler) RemoveUser(c *gin.Context) {
	claims, err := utils.GetUserClaims(c)
	if err != nil {
		c.Error(err)
		return
	}
	err = h.userUseCase.RemoveUserAndInfo(claims)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
