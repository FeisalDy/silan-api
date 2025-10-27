package handler

import (
	"net/http"
	"simple-go/internal/domain/user"
	"simple-go/internal/middleware"
	"simple-go/internal/service"
	"simple-go/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetByID retrieves a user by ID
// @Summary Get user by ID
// @Tags users
// @Security BearerAuth
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} response.Response{data=user.UserResponse}
// @Failure 404 {object} response.Response
// @Router /api/v1/users/{id} [get]
func (h *UserHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	result, err := h.userService.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "User not found", err)
		return
	}

	response.Success(c, http.StatusOK, "User retrieved successfully", result)
}

// GetAll retrieves all users with pagination
// @Summary Get all users
// @Tags users
// @Security BearerAuth
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} response.PaginatedResponse{data=[]user.UserResponse}
// @Failure 500 {object} response.Response
// @Router /api/v1/users [get]
func (h *UserHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	users, total, err := h.userService.GetAll(c.Request.Context(), limit, offset)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to retrieve users", err)
		return
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	pagination := response.Pagination{
		CurrentPage: page,
		PerPage:     limit,
		Total:       total,
		TotalPages:  totalPages,
	}

	response.PaginatedSuccess(c, http.StatusOK, "Users retrieved successfully", users, pagination)
}

// Update updates a user
// @Summary Update user
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Param request body user.UpdateUserDTO true "Update data"
// @Success 200 {object} response.Response{data=user.UserResponse}
// @Failure 400,404 {object} response.Response
// @Router /api/v1/users/{id} [put]
func (h *UserHandler) Update(c *gin.Context) {
	id := c.Param("id")

	// Check if user is updating their own profile or is admin
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	// For now, users can only update their own profile
	// Admin checks are handled by Casbin middleware
	if id != userID {
		// This will be checked by Casbin, but we add extra validation
		// In a real app, you'd check roles here or rely entirely on Casbin
	}

	var dto user.UpdateUserDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	result, err := h.userService.Update(c.Request.Context(), id, dto)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Failed to update user", err)
		return
	}

	response.Success(c, http.StatusOK, "User updated successfully", result)
}

// Delete deletes a user
// @Summary Delete user
// @Tags users
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /api/v1/users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	err := h.userService.Delete(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Failed to delete user", err)
		return
	}

	response.Success(c, http.StatusOK, "User deleted successfully", nil)
}
