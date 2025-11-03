package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"simple-go/internal/domain/job"
	"simple-go/internal/middleware"
	"simple-go/internal/service"
	"simple-go/pkg/response"

	"github.com/gin-gonic/gin"
)

type TranslationJobHandler struct {
	jobService *service.TranslationJobService
}

func NewTranslationJobHandler(jobService *service.TranslationJobService) *TranslationJobHandler {
	return &TranslationJobHandler{
		jobService: jobService,
	}
}

func (h *TranslationJobHandler) CreateTranslationJob(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	var req job.CreateTranslationJobDTO

	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid request body", response.MapValidationErrors(err, job.CreateTranslationJobDTO{}))
		return
	}

	createdJob, err := h.jobService.CreateTranslationJob(c.Request.Context(), userID, req)
	if err != nil {
		response.Error(c, http.StatusBadRequest, fmt.Sprintf("Failed to create translation job: %v", err))
		return
	}

	response.Success(c, http.StatusCreated, "Translation job created successfully", createdJob)
}

// GetJobByID retrieves a translation job by ID with all subtasks
func (h *TranslationJobHandler) GetJobByID(c *gin.Context) {
	id := c.Param("id")

	result, err := h.jobService.GetJobByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, fmt.Sprintf("Translation job not found: %v", err))
		return
	}

	response.Success(c, http.StatusOK, "Translation job retrieved successfully", result)
}

// GetAllJobs retrieves all translation jobs with pagination and optional status filter
func (h *TranslationJobHandler) GetAllJobs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.DefaultQuery("status", "")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	jobs, total, err := h.jobService.GetAllJobs(c.Request.Context(), limit, offset, status)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, fmt.Sprintf("Failed to retrieve translation jobs: %v", err))
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

	response.PaginatedSuccess(c, http.StatusOK, "Translation jobs retrieved successfully", jobs, pagination)
}

// GetJobsByNovelID retrieves all translation jobs for a specific novel
func (h *TranslationJobHandler) GetJobsByNovelID(c *gin.Context) {
	novelID := c.Param("novel_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	jobs, err := h.jobService.GetJobsByNovelID(c.Request.Context(), novelID, limit, offset)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, fmt.Sprintf("Failed to retrieve translation jobs: %v", err))
		return
	}

	response.Success(c, http.StatusOK, "Translation jobs retrieved successfully", jobs)
}

// CancelJob cancels a pending or in-progress translation job
func (h *TranslationJobHandler) CancelJob(c *gin.Context) {
	id := c.Param("id")

	err := h.jobService.CancelJob(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusBadRequest, fmt.Sprintf("Failed to cancel translation job: %v", err))
		return
	}

	response.Success(c, http.StatusOK, "Translation job cancelled successfully", nil)
}
