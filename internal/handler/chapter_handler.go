package handler

import (
	"fmt"
	"net/http"
	"simple-go/internal/domain/chapter"
	"simple-go/internal/middleware"
	"simple-go/internal/service"
	"simple-go/pkg/response"

	"github.com/gin-gonic/gin"
)

type ChapterHandler struct {
	chapterService *service.ChapterService
	volumeService  *service.VolumeService
}

func NewChapterHandler(chapterService *service.ChapterService, volumeService *service.VolumeService) *ChapterHandler {
	return &ChapterHandler{
		chapterService: chapterService,
		volumeService:  volumeService,
	}
}

func (h *ChapterHandler) Create(c *gin.Context) {
	_, exists := middleware.GetUserID(c)
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req chapter.CreateChapterDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", response.MapValidationErrors(err, chapter.CreateChapterDTO{}))
		return
	}

	newChapter, err := h.chapterService.CreateChapterWithTranslation(
		c.Request.Context(),
		req,
	)
	if err != nil {
		response.Error(c, http.StatusBadRequest, fmt.Sprintf("Failed to create chapter: %v", err))
		return
	}

	response.Success(c, http.StatusCreated, "Chapter created successfully", newChapter)
}

func (h *ChapterHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	lang := c.DefaultQuery("lang", "")

	// Use VolumeService to get chapter with cross-volume navigation
	result, err := h.volumeService.GetChapterWithCrossVolumeNavigation(c.Request.Context(), id, lang)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Failed to retrieve chapter", err)
		return
	}

	response.Success(c, http.StatusOK, "Chapter retrieved successfully", result)
}

func (h *ChapterHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	err := h.chapterService.Delete(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Failed to delete chapter", err)
		return
	}

	response.Success(c, http.StatusOK, "Chapter deleted successfully", nil)
}

func (h *ChapterHandler) CreateTranslation(c *gin.Context) {
	_, exists := middleware.GetUserID(c)
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var dto chapter.CreateChapterTranslationDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", response.MapValidationErrors(err, chapter.CreateChapterTranslationDTO{}))
		return
	}

	result, err := h.chapterService.CreateTranslation(c.Request.Context(), dto)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Failed to create translation", err)
		return
	}

	response.Success(c, http.StatusCreated, "Translation created successfully", result)
}

func (h *ChapterHandler) DeleteTranslation(c *gin.Context) {
	translationID := c.Param("id")

	err := h.chapterService.DeleteTranslation(c.Request.Context(), translationID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Failed to delete translation", err)
		return
	}

	response.Success(c, http.StatusOK, "Translation deleted successfully", nil)
}
