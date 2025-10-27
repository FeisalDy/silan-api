package handler

import (
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

// Create creates a novel with its initial translation
func (h *NovelHandler) Create(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	var req struct {
		Novel       novel.CreateNovelDTO            `json:"novel" binding:"required"`
		Translation novel.CreateNovelTranslationDTO `json:"translation" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	newNovel, newTranslation, err := h.novelService.CreateNovelWithTranslation(
		c.Request.Context(),
		userID,
		req.Novel,
		req.Translation,
	)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Failed to create novel", err)
		return
	}

	result := map[string]interface{}{
		"novel":       newNovel,
		"translation": newTranslation,
	}

	response.Success(c, http.StatusCreated, "Novel created successfully", result)
}

// GetByID retrieves a novel by ID with translations
func (h *NovelHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	result, err := h.novelService.GetByIDWithTranslations(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Novel not found", err)
		return
	}

	response.Success(c, http.StatusOK, "Novel retrieved successfully", result)
}

// GetAll retrieves all novels with pagination
func (h *NovelHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	novels, total, err := h.novelService.GetAll(c.Request.Context(), limit, offset)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to retrieve novels", err)
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

// GetMyNovels retrieves novels created by the authenticated user
func (h *NovelHandler) GetMyNovels(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	novels, err := h.novelService.GetByCreator(c.Request.Context(), userID, limit, offset)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to retrieve novels", err)
		return
	}

	response.Success(c, http.StatusOK, "Novels retrieved successfully", novels)
}

// Update updates a novel
func (h *NovelHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var dto novel.UpdateNovelDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	result, err := h.novelService.Update(c.Request.Context(), id, dto)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Failed to update novel", err)
		return
	}

	response.Success(c, http.StatusOK, "Novel updated successfully", result)
}

// Delete deletes a novel
func (h *NovelHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	err := h.novelService.Delete(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Failed to delete novel", err)
		return
	}

	response.Success(c, http.StatusOK, "Novel deleted successfully", nil)
}

// CreateTranslation creates a new translation for a novel
func (h *NovelHandler) CreateTranslation(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Unauthorized", nil)
		return
	}

	var dto novel.CreateNovelTranslationDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	result, err := h.novelService.CreateTranslation(c.Request.Context(), userID, dto)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Failed to create translation", err)
		return
	}

	response.Success(c, http.StatusCreated, "Translation created successfully", result)
}

// GetTranslation retrieves a translation by novel ID and language
func (h *NovelHandler) GetTranslation(c *gin.Context) {
	novelID := c.Param("id")
	lang := c.Param("lang")

	result, err := h.novelService.GetTranslation(c.Request.Context(), novelID, lang)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Translation not found", err)
		return
	}

	response.Success(c, http.StatusOK, "Translation retrieved successfully", result)
}

// UpdateTranslation updates a translation
func (h *NovelHandler) UpdateTranslation(c *gin.Context) {
	translationID := c.Param("translation_id")

	var dto novel.UpdateNovelTranslationDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	result, err := h.novelService.UpdateTranslation(c.Request.Context(), translationID, dto)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Failed to update translation", err)
		return
	}

	response.Success(c, http.StatusOK, "Translation updated successfully", result)
}

// DeleteTranslation deletes a translation
func (h *NovelHandler) DeleteTranslation(c *gin.Context) {
	translationID := c.Param("translation_id")

	err := h.novelService.DeleteTranslation(c.Request.Context(), translationID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Failed to delete translation", err)
		return
	}

	response.Success(c, http.StatusOK, "Translation deleted successfully", nil)
}
