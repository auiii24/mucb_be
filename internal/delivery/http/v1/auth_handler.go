package v1

import (
	"mucb_be/internal/errors"
	"mucb_be/internal/usecase/auth"
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

	// TODO:
	// if err := request.Validate(); err != nil {
	// 	c.Error(
	// 		errors.NewCustomError(http.StatusBadRequest, "VE001002", "Invalid phone number.", err.Error()),
	// 	)
	// 	return
	// }

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

	response, err := h.authUseCase.VerifyOtp(&request)
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
			errors.NewCustomError(http.StatusBadRequest, "VE001001", err.Error(), err.Error()),
		)
		return
	}

	var userAgent = c.GetHeader("User-Agent")
	response, err := h.authUseCase.RenewUser(&request, userAgent)
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
