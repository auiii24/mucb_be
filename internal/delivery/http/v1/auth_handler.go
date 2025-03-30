package v1

import (
	"mucb_be/internal/errors"
	"mucb_be/internal/usecase/auth"
	"mucb_be/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authUseCase auth.AuthUseCaseInterface
}

func NewAuthHandler(authUseCase auth.AuthUseCaseInterface) *AuthHandler {
	return &AuthHandler{authUseCase: authUseCase}
}

func (h AuthHandler) SignInAdmin(c *gin.Context) {
	var request auth.SignInAdminRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(
			errors.NewCustomError(http.StatusBadRequest, "VE001001", err.Error(), err.Error()),
		)
		return
	}

	var userAgent = c.GetHeader("User-Agent")
	response, err := h.authUseCase.SignInAdmin(&request, userAgent)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h AuthHandler) RenewAdmin(c *gin.Context) {
	var request auth.RenewAdminRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(
			errors.NewCustomError(http.StatusBadRequest, "VE001001", err.Error(), err.Error()),
		)
		return
	}

	var userAgent = c.GetHeader("User-Agent")
	response, err := h.authUseCase.RenewAdmin(&request, userAgent)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h AuthHandler) SignInUser(c *gin.Context) {
	var request auth.SignInUserRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(
			errors.NewCustomError(http.StatusBadRequest, "VE001001", err.Error(), err.Error()),
		)
		return
	}

	response, err := h.authUseCase.SignInUser(&request)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h AuthHandler) VerifyOtpUser(c *gin.Context) {
	var request auth.VerifyOtpRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(
			errors.NewCustomError(http.StatusBadRequest, "VE001001", err.Error(), err.Error()),
		)
		return
	}

	response, err := h.authUseCase.VerifyOtp(&request, c)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h AuthHandler) RenewUser(c *gin.Context) {
	var request auth.RenewUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(
			errors.NewCustomError(http.StatusUnauthorized, "VE001001", err.Error(), err.Error()),
		)
		return
	}

	response, err := h.authUseCase.RenewUser(&request, c)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h AuthHandler) SignOut(c *gin.Context) {
	var request auth.SignOutRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(
			errors.NewCustomError(http.StatusBadRequest, "VE001001", err.Error(), err.Error()),
		)
		return
	}

	err := h.authUseCase.SignOut(&request)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

func (h AuthHandler) GetAvailableTokens(c *gin.Context) {
	var request auth.FindAllTokensRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(
			errors.NewCustomError(http.StatusBadRequest, "VE001001", err.Error(), err.Error()),
		)
		return
	}

	claims, err := utils.GetUserClaims(c)
	if err != nil {
		c.Error(err)
		return
	}

	response, err := h.authUseCase.FindAllToken(&request, claims)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h AuthHandler) RevokeToken(c *gin.Context) {
	var request auth.RevokeTokenRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.Error(
			errors.NewCustomError(http.StatusBadRequest, "VE001001", err.Error(), err.Error()),
		)
		return
	}

	claims, err := utils.GetUserClaims(c)
	if err != nil {
		c.Error(err)
		return
	}

	err = h.authUseCase.RevokeToken(&request, claims)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
