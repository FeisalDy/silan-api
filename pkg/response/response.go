package response

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// Response is the standard API response structure
type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Error   any    `json:"error,omitempty"`
}

// PaginatedResponse is for paginated list responses
type PaginatedResponse struct {
	Success bool       `json:"success"`
	Message string     `json:"message,omitempty"`
	Data    any        `json:"data"`
	Meta    Pagination `json:"meta"`
}

// Pagination contains pagination metadata
type Pagination struct {
	CurrentPage int   `json:"current_page"`
	PerPage     int   `json:"per_page"`
	Total       int64 `json:"total"`
	TotalPages  int   `json:"total_pages"`
}

// Success sends a successful response
func Success(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Error sends an error response
func Error(c *gin.Context, statusCode int, message string, errorDetails ...any) {
	var errDetail any
	if len(errorDetails) > 0 {
		errDetail = errorDetails[0]
	} else {
		errDetail = nil
	}

	c.JSON(statusCode, Response{
		Success: false,
		Message: message,
		Error:   errDetail,
	})
}

// PaginatedSuccess sends a paginated successful response
func PaginatedSuccess(c *gin.Context, statusCode int, message string, data interface{}, pagination Pagination) {
	c.JSON(statusCode, PaginatedResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta:    pagination,
	})
}

func MapValidationErrors(err error) map[string]string {
	errors := make(map[string]string)

	if err == nil {
		return errors
	}

	// Check if it's a validator.ValidationErrors
	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range ve {
			// Convert StructNamespace to JSON-friendly key, e.g. "Novel.OriginalLanguage" -> "original_language"
			key := fe.StructNamespace() // e.g. "Novel.OriginalLanguage"
			keyParts := strings.Split(key, ".")
			key = strings.ToLower(keyParts[len(keyParts)-1]) // take the last part and lowercase
			errors[key] = fmt.Sprintf("Field validation for '%s' failed on the '%s' tag", fe.Field(), fe.Tag())
		}
	} else {
		// fallback
		errors["error"] = err.Error()
	}

	return errors
}
