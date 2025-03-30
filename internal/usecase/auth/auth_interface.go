package auth

import (
	"mucb_be/internal/infrastructure/security"

	"github.com/gin-gonic/gin"
)

type AuthUseCaseInterface interface {
	SignInAdmin(req *SignInAdminRequest, userAgent string) (*SignInAdminOutput, error)
	RenewAdmin(req *RenewAdminRequest, userAgent string) (*RenewAdminOutput, error)
	SignInUser(req *SignInUserRequest) (*SignInUserOutput, error)
	VerifyOtp(req *VerifyOtpRequest, c *gin.Context) (*VerifyOtpOutput, error)
	RenewUser(req *RenewUserRequest, c *gin.Context) (*RenewUserOutput, error)
	FindAllToken(req *FindAllTokensRequest, claims *security.AccessTokenModel) (*FindAllTokensOutput, error)
	SignOut(req *SignOutRequest) error
	RevokeToken(req *RevokeTokenRequest, claims *security.AccessTokenModel) error
}
