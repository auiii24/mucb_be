package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func InvalidEndpointMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    "MWE002001001",
			"message": "The requested URL does not exist",
		})
	}
}
