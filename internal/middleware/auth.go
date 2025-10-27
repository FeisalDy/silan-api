package middleware

import (
	"net/http"
	"strings"

	"simple-go/pkg/auth"
	"simple-go/pkg/logger"
	"simple-go/pkg/response"

	"github.com/gin-gonic/gin"
)

// JWTAuth is a middleware that validates JWT tokens
func JWTAuth(jwtManager *auth.JWTManager) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "Authorization header required")
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Error(c, http.StatusUnauthorized, "Invalid authorization header format")
			c.Abort()
			return
		}

		token := parts[1]

		claims, err := jwtManager.ValidateToken(token)
		if err != nil {
			logger.Error(err, "failed to validate JWT token")
			response.Error(c, http.StatusUnauthorized, "Invalid or expired token")
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("userEmail", claims.Email)

		c.Next()
	}
}

func GetUserID(c *gin.Context) (string, bool) {
	userID, exists := c.Get("userID")
	if !exists {
		return "", false
	}
	id, ok := userID.(string)
	return id, ok
}

func GetUserEmail(c *gin.Context) (string, bool) {
	userEmail, exists := c.Get("userEmail")
	if !exists {
		return "", false
	}
	email, ok := userEmail.(string)
	return email, ok
}
