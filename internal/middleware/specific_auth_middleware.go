package middleware

import (
	"errors"
	"mucb_be/internal/infrastructure/security"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// SpecificAuthMiddleware validates JWT and checks role permissions
func SpecificAuthMiddleware(jwtService security.JwtServiceInterface, requiredRoles []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// ✅ Extract JWT token from Authorization header
		authHeader := c.GetHeader("X-Authorization")

		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    "MWE001001001",
				"message": "Missing authorization header.",
			})
			return
		}

		// ✅ Ensure token format is "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    "MWE001001002",
				"message": "Invalid authorization header format.",
			})
			return
		}
		token := tokenParts[1]

		// ✅ Decode JWT token
		claims, err := jwtService.ValidateAccessToken(token)
		if err != nil {
			var statusCode int
			var code string
			var message string = "Invalid token."

			if errors.Is(err, jwt.ErrTokenExpired) {
				err = security.ErrJWTTokenExpired
			}

			switch {
			case errors.Is(err, security.ErrJWTTokenExpired):
				statusCode = http.StatusUnauthorized
				code = "MWE001001003"
				message = "Token has expired"
			case errors.Is(err, security.ErrJWTMissingIDClaim), errors.Is(err, security.ErrJWTMissingRoleClaim):
				statusCode = http.StatusUnauthorized
				code = "MWE001001004"
			default:
				statusCode = http.StatusUnauthorized
				code = "MWE001001005"
			}

			c.AbortWithStatusJSON(statusCode, gin.H{
				"code":    code,
				"message": message,
			})
			return
		}

		// ✅ Check required role if specified
		if !roleAllowed(claims.Role, requiredRoles) {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    "MWE001001006",
				"message": "Insufficient permissions.",
			})
			return
		}

		// ✅ Store claims in context
		c.Set("user", claims)
		c.Next()
	}
}

// ✅ Function to check if user role is allowed
func roleAllowed(userRole string, requiredRoles []string) bool {
	// ✅ If no roles are specified, allow all users
	if len(requiredRoles) == 0 {
		return true
	}

	// ✅ Check if userRole is in the requiredRoles list
	for _, role := range requiredRoles {
		if userRole == role {
			return true
		}
	}
	return false
}
