package handler

import (
	"fmt"
	"net/http"
	"simple-go/internal/domain/chapter"
	"simple-go/internal/middleware"
	"simple-go/internal/service"
	"simple-go/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ChapterHandler handles chapter-related HTTP requests
type ChapterHandler struct {
	chapterService *service.ChapterService
}

func NewChapterHandler(chapterService *service.ChapterService) *ChapterHandler {
	return &ChapterHandler{
		chapterService: chapterService,
	}
}

// Create creates a chapter with its initial translation
func (h *ChapterHandler) Create(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req struct {
		Chapter     chapter.CreateChapterDTO            `json:"chapter" binding:"required"`
		Translation chapter.CreateChapterTranslationDTO `json:"translation" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, fmt.Sprintf("Invalid request body: %v", err))
		return
	}

	newChapter, newTranslation, err := h.chapterService.CreateChapterWithTranslation(
		c.Request.Context(),
		userID,
		req.Chapter,
		req.Translation,
	)
	if err != nil {
		response.Error(c, http.StatusBadRequest, fmt.Sprintf("Failed to create chapter: %v", err))
		return
	}

	result := map[string]interface{}{
		"chapter":     newChapter,
		"translation": newTranslation,
	}

	response.Success(c, http.StatusCreated, "Chapter created successfully", result)
}

// GetByID retrieves a chapter by ID with translations
func (h *ChapterHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	result, err := h.chapterService.GetByIDWithTranslations(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, fmt.Sprintf("Chapter not found: %v", err))
		return
	}

	response.Success(c, http.StatusOK, "Chapter retrieved successfully", result)
}

// GetByNovel retrieves all chapters for a novel
func (h *ChapterHandler) GetByNovel(c *gin.Context) {
	novelID := c.Query("novel_id")
	if novelID == "" {
		response.Error(c, http.StatusBadRequest, "novel_id query parameter is required")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 50
	}

	offset := (page - 1) * limit

	chapters, total, err := h.chapterService.GetByNovel(c.Request.Context(), novelID, limit, offset)
	if err != nil {
		response.Error(c, http.StatusBadRequest, fmt.Sprintf("Failed to retrieve chapters: %v", err))
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

	response.PaginatedSuccess(c, http.StatusOK, "Chapters retrieved successfully", chapters, pagination)
}

// GetByNovelAndNumber retrieves a specific chapter by novel and number
func (h *ChapterHandler) GetByNovelAndNumber(c *gin.Context) {
	novelID := c.Query("novel_id")
	numberStr := c.Query("number")

	if novelID == "" || numberStr == "" {
		response.Error(c, http.StatusBadRequest, "novel_id and number query parameters are required")
		return
	}

	number, err := strconv.Atoi(numberStr)
	if err != nil {
		response.Error(c, http.StatusBadRequest, fmt.Sprintf("Invalid number format: %v", err))
		return
	}

	result, err := h.chapterService.GetByNovelAndNumber(c.Request.Context(), novelID, number)
	if err != nil {
		response.Error(c, http.StatusNotFound, fmt.Sprintf("Chapter not found: %v", err))
		return
	}

	response.Success(c, http.StatusOK, "Chapter retrieved successfully", result)
}

// Update updates a chapter
func (h *ChapterHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var dto chapter.UpdateChapterDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.Error(c, http.StatusBadRequest, fmt.Sprintf("Invalid request body: %v", err))
		return
	}

	result, err := h.chapterService.Update(c.Request.Context(), id, dto)
	if err != nil {
		response.Error(c, http.StatusBadRequest, fmt.Sprintf("Failed to update chapter: %v", err))
		return
	}

	response.Success(c, http.StatusOK, "Chapter updated successfully", result)
}

// Delete deletes a chapter
func (h *ChapterHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	err := h.chapterService.Delete(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, fmt.Sprintf("Failed to delete chapter: %v", err))
		return
	}

	response.Success(c, http.StatusOK, "Chapter deleted successfully", nil)
}

// CreateTranslation creates a new translation for a chapter
func (h *ChapterHandler) CreateTranslation(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var dto chapter.CreateChapterTranslationDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.Error(c, http.StatusBadRequest, fmt.Sprintf("Invalid request body: %v", err))
		return
	}

	result, err := h.chapterService.CreateTranslation(c.Request.Context(), userID, dto)
	if err != nil {
		response.Error(c, http.StatusBadRequest, fmt.Sprintf("Failed to create translation: %v", err))
		return
	}

	response.Success(c, http.StatusCreated, "Translation created successfully", result)
}

// GetTranslation retrieves a translation by chapter ID and language
func (h *ChapterHandler) GetTranslation(c *gin.Context) {
	chapterID := c.Param("id")
	lang := c.Param("lang")

	result, err := h.chapterService.GetTranslation(c.Request.Context(), chapterID, lang)
	if err != nil {
		response.Error(c, http.StatusNotFound, fmt.Sprintf("Translation not found: %v", err))
		return
	}

	response.Success(c, http.StatusOK, "Translation retrieved successfully", result)
}

// UpdateTranslation updates a translation
func (h *ChapterHandler) UpdateTranslation(c *gin.Context) {
	translationID := c.Param("translation_id")

	var dto chapter.UpdateChapterTranslationDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.Error(c, http.StatusBadRequest, fmt.Sprintf("Invalid request body: %v", err))
		return
	}

	result, err := h.chapterService.UpdateTranslation(c.Request.Context(), translationID, dto)
	if err != nil {
		response.Error(c, http.StatusBadRequest, fmt.Sprintf("Failed to update translation: %v", err))
		return
	}

	response.Success(c, http.StatusOK, "Translation updated successfully", result)
}

// DeleteTranslation deletes a translation
func (h *ChapterHandler) DeleteTranslation(c *gin.Context) {
	translationID := c.Param("translation_id")

	err := h.chapterService.DeleteTranslation(c.Request.Context(), translationID)
	if err != nil {
		response.Error(c, http.StatusNotFound, fmt.Sprintf("Failed to delete translation: %v", err))
		return
	}

	response.Success(c, http.StatusOK, "Translation deleted successfully", nil)
}
