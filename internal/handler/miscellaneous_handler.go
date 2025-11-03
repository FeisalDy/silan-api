package handler

import (
	"net/http"
	"simple-go/pkg/miscellaneous"
	"simple-go/pkg/response"

	"github.com/gin-gonic/gin"
)

type MiscellaneousHandler struct{}

func NewMiscellaneousHandler() *MiscellaneousHandler {
	return &MiscellaneousHandler{}
}

// GetLanguages retrieves languages by code or search
// Query parameters:
//   - code: ISO 639-1 code to get specific language
//   - search: search query to match in name or native fields
//
// If code is provided, search is ignored
func (h *MiscellaneousHandler) GetLanguages(c *gin.Context) {
	code := c.DefaultQuery("code", "")
	search := c.DefaultQuery("search", "")

	// If code is provided, get specific language by code
	if code != "" {
		language := miscellaneous.GetLanguageByCode(code)
		if language == nil {
			response.Error(c, http.StatusNotFound, "Language not found")
			return
		}
		response.Success(c, http.StatusOK, "Language retrieved successfully", language)
		return
	}

	// If search is provided, search in name and native
	if search != "" {
		results := miscellaneous.SearchLanguages(search)
		response.Success(c, http.StatusOK, "Languages search completed", results)
		return
	}

	// If neither provided, return all languages
	languages := miscellaneous.GetAllLanguages()
	response.Success(c, http.StatusOK, "Languages retrieved successfully", languages)
}
