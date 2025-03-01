package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequestLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		const maxUploadSize = 5 * 1024 * 1024

		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxUploadSize)

		if c.Request.ContentLength > maxUploadSize {
			c.AbortWithStatusJSON(http.StatusRequestEntityTooLarge, gin.H{
				"code":    "MWE003001001",
				"message": "Request size exceeds.",
			})
			return
		}

		c.Next()
	}
}
