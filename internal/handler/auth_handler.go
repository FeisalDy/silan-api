package handler

import (
	"fmt"
	"net/http"
	"simple-go/internal/service"
	"simple-go/pkg/response"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req service.RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, fmt.Sprintf("Invalid request body: %v", err))
		return
	}

	result, err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, fmt.Sprintf("Registration failed: %v", err))
		return
	}

	response.Success(c, http.StatusCreated, "User registered successfully", result)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req service.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, fmt.Sprintf("Invalid request body: %v", err))
		return
	}

	result, err := h.authService.Login(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, fmt.Sprintf("Login failed: %v", err))
		return
	}

	response.Success(c, http.StatusOK, "Login successful", result)
}
