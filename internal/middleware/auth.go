package middleware

import (
	"strings"

	"san/pkg/apperr"
	"san/pkg/response"
	"san/pkg/token"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(tokenManager token.TokenManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, apperr.Unauthorized("Authorization header is required"))
			c.Abort()
			return
		}

		fields := strings.Fields(authHeader)
		if len(fields) != 2 || strings.ToLower(fields[0]) != "bearer" {
			response.Error(c, apperr.Unauthorized("Invalid authorization header format"))
			c.Abort()
			return
		}

		tokenString := fields[1]
		claims, err := tokenManager.VerifyToken(tokenString)
		if err != nil {
			response.Error(c, apperr.Unauthorized(err.Error()))
			c.Abort()
			return
		}

		if claims.TokenType != "access" {
			response.Error(c, apperr.Unauthorized("Invalid token type"))
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Next()
	}
}
