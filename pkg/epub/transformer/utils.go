package transformer

import (
	"fmt"
	"strings"
)

// getBaseDir extracts the base directory from OPF path
func getBaseDir(opfPath string) string {
	if idx := strings.LastIndex(opfPath, "/"); idx >= 0 {
		return opfPath[:idx+1]
	}
	return ""
}

// extractChapterTitle tries to extract chapter title from HTML content
func extractChapterTitle(htmlContent string, defaultNum int) string {
	// Try to extract title from h1 tag
	if strings.Contains(htmlContent, "<h1>") {
		start := strings.Index(htmlContent, "<h1>") + 4
		end := strings.Index(htmlContent[start:], "</h1>")
		if end > 0 {
			return strings.TrimSpace(htmlContent[start : start+end])
		}
	}

	// Try to extract title from h2 tag
	if strings.Contains(htmlContent, "<h2>") {
		start := strings.Index(htmlContent, "<h2>") + 4
		end := strings.Index(htmlContent[start:], "</h2>")
		if end > 0 {
			return strings.TrimSpace(htmlContent[start : start+end])
		}
	}

	// Fallback to default chapter number
	return fmt.Sprintf("Chapter %d", defaultNum)
}
