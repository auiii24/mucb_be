package utils

import (
	"mucb_be/internal/errors"
	"mucb_be/internal/infrastructure/security"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ✅ GetUserClaims extracts user claims from Gin context
func GetUserClaims(c *gin.Context) (*security.AccessTokenModel, error) {
	userClaims, exists := c.Get("user")
	if !exists {
		message := "Invalid token: id claim is missing or not a string"
		return nil, errors.NewCustomError(
			http.StatusUnauthorized, "VE001001002", message, message,
		)
	}

	claims, ok := userClaims.(security.AccessTokenModel)
	if !ok {
		message := "Invalid token: role claim is missing"
		return nil, errors.NewCustomError(
			http.StatusUnauthorized, "VE001001002", message, message,
		)
	}

	return &claims, nil
}
