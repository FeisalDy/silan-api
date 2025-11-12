package response

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// Response is the standard API response structure
type Response struct {
	Success   bool   `json:"success"`
	Message   string `json:"message,omitempty"`
	Data      any    `json:"data,omitempty"`
	ErrorCode string `json:"error_code,omitempty"`
	Error     any    `json:"error,omitempty"`
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
	Limit       int   `json:"limit"`
	Total       int64 `json:"total"`
	TotalPages  int   `json:"total_pages"`
}

const (
	ErrCodeUserAlreadyExists      = "USER002"
	ErrCodeUserCreationFailed     = "USER003"
	ErrCodeUserUpdateFailed       = "USER004"
	ErrCodeUserDeletionFailed     = "USER005"
	ErrCodeUserInvalidCredentials = "USER006"
	ErrCodeUserValidation         = "USER007"
)

// Success sends a successful response
func Success(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

// Error is DEPRECATED and will be removed in a future release.
// Use ErrorWithCode or a custom response function instead.
func Error(c *gin.Context, statusCode int, messageOrErr any, errorDetails ...any) {
	var (
		message   string
		errDetail any
	)

	if err, ok := messageOrErr.(error); ok {
		message = err.Error()
		errDetail = nil
	} else if msg, ok := messageOrErr.(string); ok {
		message = msg
	} else {
		message = "Unknown error"
	}

	if len(errorDetails) > 0 {
		detail := errorDetails[0]
		if e, ok := detail.(error); ok {
			errDetail = e.Error()
		} else {
			errDetail = detail
		}
	}

	c.JSON(statusCode, Response{
		Success: false,
		Message: message,
		Error:   errDetail,
	})
}

func ErrorWithCode(c *gin.Context, statusCode int, errorCode string, messageOrErr any, errorDetails ...any) {
	var (
		message   string
		errDetail any
	)

	if err, ok := messageOrErr.(error); ok {
		message = err.Error()
		errDetail = nil
	} else if msg, ok := messageOrErr.(string); ok {
		message = msg
	} else {
		message = "Unknown error"
	}

	if len(errorDetails) > 0 {
		detail := errorDetails[0]
		if e, ok := detail.(error); ok {
			errDetail = e.Error()
		} else {
			errDetail = detail
		}
	}

	c.JSON(statusCode, Response{
		Success:   false,
		Message:   message,
		ErrorCode: errorCode,
		Error:     errDetail,
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

func MapValidationErrors(err error, dto any) map[string]string {
	errors := make(map[string]string)

	if err == nil {
		return errors
	}

	if ve, ok := err.(validator.ValidationErrors); ok {
		t := reflect.TypeOf(dto)
		if t.Kind() == reflect.Pointer {
			t = t.Elem()
		}

		for _, fe := range ve {
			fieldName := fe.StructField() // e.g. "OriginalLanguage"
			field, _ := t.FieldByName(fieldName)

			// Get json tag name
			jsonTag := field.Tag.Get("json")
			jsonKey := strings.Split(jsonTag, ",")[0]
			if jsonKey == "" || jsonKey == "-" {
				jsonKey = strings.ToLower(fieldName)
			}

			errors[jsonKey] = fmt.Sprintf("Field validation for '%s' failed on the '%s' tag", jsonKey, fe.Tag())
		}
	} else {
		errors["error"] = err.Error()
	}

	return errors
}
