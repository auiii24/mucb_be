package middleware

import (
	"mucb_be/internal/errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Let the handler execute
		c.Next()

		// Check if there are any errors
		if len(c.Errors) > 0 {
			// Get the last error (you can handle multiple errors if needed)
			err := c.Errors.Last().Err

			// Handle custom errors
			if customErr, ok := err.(*errors.CustomError); ok {
				c.JSON(customErr.StatusCode, customErr)
				return
			}

			// For generic errors, return a default 500 response
			c.JSON(http.StatusInternalServerError, gin.H{
				"code":    "ERR000000001",
				"message": "Internal Server Error",
			})
		}
	}
}
