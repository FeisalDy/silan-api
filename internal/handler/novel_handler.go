package handler

import (
	"fmt"
	"net/http"
	"simple-go/internal/domain/novel"
	"simple-go/internal/middleware"
	"simple-go/internal/service"
	"simple-go/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// NovelHandler handles novel-related HTTP requests
type NovelHandler struct {
	novelService *service.NovelService
}

func NewNovelHandler(novelService *service.NovelService) *NovelHandler {
	return &NovelHandler{
		novelService: novelService,
	}
}

func (h *NovelHandler) Create(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req struct {
		Novel       novel.CreateNovelDTO            `json:"novel" binding:"required"`
		Translation novel.CreateNovelTranslationDTO `json:"translation" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", response.MapValidationErrors(err))
		return
	}

	newNovel, newTranslation, err := h.novelService.CreateNovelWithTranslation(
		c.Request.Context(),
		userID,
		req.Novel,
		req.Translation,
	)

	if err != nil {
		response.Error(c, http.StatusBadRequest, fmt.Sprintf("Failed to create novel: %v", err))
		return
	}

	result := map[string]interface{}{
		"novel":       newNovel,
		"translation": newTranslation,
	}

	response.Success(c, http.StatusCreated, "Novel created successfully", result)
}

func (h *NovelHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	lang := c.DefaultQuery("lang", "")

	result, err := h.novelService.GetByIDWithTranslations(c.Request.Context(), id, lang)
	if err != nil {
		response.Error(c, http.StatusNotFound, fmt.Sprintf("Novel not found: %v", err))
		return
	}

	response.Success(c, http.StatusOK, "Novel retrieved successfully", result)
}

func (h *NovelHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	title := c.DefaultQuery("title", "")
	lang := c.DefaultQuery("lang", "")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	novels, total, err := h.novelService.GetAll(c.Request.Context(), limit, offset, title, lang)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, fmt.Sprintf("Failed to retrieve novels: %v", err))
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

	response.PaginatedSuccess(c, http.StatusOK, "Novels retrieved successfully", novels, pagination)
}

func (h *NovelHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	err := h.novelService.Delete(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, fmt.Sprintf("Failed to delete novel: %v", err))
		return
	}

	response.Success(c, http.StatusOK, "Novel deleted successfully", nil)
}

func (h *NovelHandler) CreateTranslation(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var dto novel.CreateNovelTranslationDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.Error(c, http.StatusBadRequest, fmt.Sprintf("Invalid request body: %v", err))
		return
	}

	result, err := h.novelService.CreateTranslation(c.Request.Context(), userID, dto)
	if err != nil {
		response.Error(c, http.StatusBadRequest, fmt.Sprintf("Failed to create translation: %v", err))
		return
	}

	response.Success(c, http.StatusCreated, "Translation created successfully", result)
}

func (h *NovelHandler) DeleteTranslation(c *gin.Context) {
	translationID := c.Param("translation_id")

	err := h.novelService.DeleteTranslation(c.Request.Context(), translationID)
	if err != nil {
		response.Error(c, http.StatusNotFound, fmt.Sprintf("Failed to delete translation: %v", err))
		return
	}

	response.Success(c, http.StatusOK, "Translation deleted successfully", nil)
}
