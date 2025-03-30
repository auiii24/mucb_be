package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CorsMiddleware(allowOrigin string) gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:  []string{allowOrigin},
		AllowMethods:  []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "HEAD"},
		AllowHeaders:  []string{"Content-Type", "X-Authorization", "X-API-KEY"},
		ExposeHeaders: []string{"Content-Length"},
	})
}
