package middleware

import (
	"elderly-care-system/pkg/jwt"
	"elderly-care-system/pkg/response"
	"strings"

	"github.com/gin-gonic/gin"
)

func Auth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "Missing authorization header")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "Invalid authorization format")
			c.Abort()
			return
		}

		claims, err := jwt.ParseToken(parts[1], secret)
		if err != nil {
			response.Unauthorized(c, "Invalid token")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("phone", claims.Phone)
		c.Next()
	}
}
