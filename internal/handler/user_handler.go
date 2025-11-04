package handler

import (
	"fmt"
	"net/http"
	"simple-go/internal/domain/user"
	"simple-go/internal/middleware"
	"simple-go/internal/service"
	"simple-go/pkg/logger"
	"simple-go/pkg/response"
	"slices"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	result, err := h.userService.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, fmt.Sprintf("User not found: %v", err))
		return
	}

	response.Success(c, http.StatusOK, "User retrieved successfully", result)
}

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
		response.Error(c, http.StatusInternalServerError, fmt.Sprintf("Failed to retrieve users: %v", err))
		return
	}

	totalPages := int(total) / limit
	if int(total)%limit > 0 {
		totalPages++
	}

	pagination := response.Pagination{
		CurrentPage: page,
		Limit:       limit,
		Total:       total,
		TotalPages:  totalPages,
	}

	response.PaginatedSuccess(c, http.StatusOK, "Users retrieved successfully", users, pagination)
}

func (h *UserHandler) Update(c *gin.Context) {
	id := c.Param("id")

	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if id != userID {
		roles, err := h.userService.GetUserRoles(c.Request.Context(), userID)
		if err != nil {
			logger.Error(err, "failed to get user roles")
			response.Error(c, http.StatusInternalServerError, fmt.Sprintf("Failed to check user role: %v", err))
			return
		}

		isAdmin := slices.Contains(roles, "admin")

		if !isAdmin {
			response.Error(c, http.StatusForbidden, "You are not allowed to update other users")
			return
		}
	}

	var dto user.UpdateUserDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.Error(c, http.StatusBadRequest, fmt.Sprintf("Invalid request body: %v", err))
		return
	}

	result, err := h.userService.Update(c.Request.Context(), id, dto)
	if err != nil {
		logger.Error(err, "failed to update user")
		response.Error(c, http.StatusInternalServerError, fmt.Sprintf("Failed to update user: %v", err))
		return
	}

	response.Success(c, http.StatusOK, "User updated successfully", result)
}

func (h *UserHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	err := h.userService.Delete(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, fmt.Sprintf("Failed to delete user: %v", err))
		return
	}

	response.Success(c, http.StatusOK, "User deleted successfully", nil)
}
