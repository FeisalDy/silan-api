package handler

import (
	"fmt"
	"io"
	"net/http"
	"simple-go/internal/domain/novel"
	"simple-go/internal/middleware"
	"simple-go/internal/service"
	"simple-go/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

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

	var req novel.CreateNovelDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", response.MapValidationErrors(err, novel.CreateNovelDTO{}))
		return
	}

	newNovel, newTranslation, err := h.novelService.CreateNovelWithTranslation(
		c.Request.Context(),
		userID,
		req,
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

	result, err := h.novelService.GetByID(c.Request.Context(), id, lang)
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
		response.Error(c, http.StatusInternalServerError, fmt.Sprintf("Failed to delete novel: %v", err))
		return
	}

	response.Success(c, http.StatusOK, "Novel deleted successfully", nil)
}

func (h *NovelHandler) UpdateCoverMedia(c *gin.Context) {
	id := c.Param("id")

	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	fileHeader, err := c.FormFile("cover_media")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Missing cover_media file")
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to open file")
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to read file", nil)
		return
	}

	var req novel.UpdateCoverMediaDTO

	req.FileName = fileHeader.Filename
	req.FileBytes = fileBytes
	req.UploaderID = userID

	if err := h.novelService.UpdateCoverMedia(c.Request.Context(), id, req); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update cover media", err)
		return
	}

	response.Success(c, http.StatusOK, "Cover media updated successfully", nil)
}

func (h *NovelHandler) CreateTranslation(c *gin.Context) {
	_, exists := middleware.GetUserID(c)
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var dto novel.CreateNovelTranslationDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		response.Error(c, http.StatusBadRequest, fmt.Sprintf("Invalid request body: %v", err))
		return
	}

	result, err := h.novelService.CreateTranslation(c.Request.Context(), dto)
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

func (h *NovelHandler) GetNovelVolumes(c *gin.Context) {
	id := c.Param("id")
	lang := c.DefaultQuery("lang", "")

	volumes, err := h.novelService.GetNovelVolumes(c.Request.Context(), id, lang)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, fmt.Sprintf("Failed to retrieve novel volumes: %v", err))
		return
	}

	response.Success(c, http.StatusOK, "Novel volumes retrieved successfully", volumes)
}

func (h *NovelHandler) UploadEpub(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	fileHeader, err := c.FormFile("epub_file")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Missing epub_file in form data")
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to open epub file")
		return
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to read epub file")
		return
	}

	if len(fileBytes) == 0 {
		response.Error(c, http.StatusBadRequest, "Epub file is empty")
		return
	}

	result, err := h.novelService.ProcessAndSaveEpubUpload(c.Request.Context(), fileBytes, userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to process epub file", err)
		return
	}

	// Prepare volume summary
	volumeSummary := []map[string]interface{}{}
	for _, vol := range result.Volumes {
		volumeSummary = append(volumeSummary, map[string]interface{}{
			"number":     vol.Number,
			"title":      vol.Title,
			"is_virtual": vol.IsVirtual,
		})
	}

	// Prepare response with parsed data from the processing result
	responseData := map[string]interface{}{
		"source_type": result.SourceType,
		"novel_data": map[string]interface{}{
			"title":             result.NovelData.Title,
			"original_author":   result.NovelData.OriginalAuthor,
			"original_language": result.NovelData.OriginalLanguage,
			"description":       result.NovelData.Description,
			"publisher":         result.NovelData.Publisher,
			"tags":              result.NovelData.Tags,
			"has_cover_image":   len(result.NovelData.CoverImage) > 0,
		},
		"volumes":        volumeSummary,
		"total_volumes":  result.TotalVolumes,
		"total_chapters": result.TotalChapters,
		"total_files":    len(result.RawContent.RawFiles),
	}

	response.Success(c, http.StatusOK, "EPUB file parsed successfully. Check console for detailed output.", responseData)
}
