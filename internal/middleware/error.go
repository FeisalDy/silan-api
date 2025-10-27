package middleware

import (
	"fmt"
	"net/http"

	"simple-go/pkg/response"

	"github.com/gin-gonic/gin"
)

// ErrorHandler is a middleware that handles panics and errors
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				response.Error(c, http.StatusInternalServerError, "Internal server error")
				c.Abort()
			}
		}()

		c.Next()

		// Check for errors after request processing
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			response.Error(c, http.StatusInternalServerError, fmt.Sprintf("An error occurred: %v", err.Err))
		}
	}
}

// CORS middleware
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
