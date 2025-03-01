package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func BasicAuthMiddleware(apiKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestAPIKey := c.GetHeader("X-API-KEY")

		if requestAPIKey == "" || requestAPIKey != apiKey {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    "MWE004001001",
				"message": "Invalid or missing API key",
			})
			return
		}

		c.Next()
	}
}
