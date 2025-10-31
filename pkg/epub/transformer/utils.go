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

// extractVolumeNumberFromFilename extracts volume number from 404 source filenames
// e.g., "No00001Chapter.xhtml" -> 1, "No00042Section.xhtml" -> 42
func extractVolumeNumberFromFilename(filename string) int {
	// Extract the numeric part after "No" and before "Chapter" or "Section"
	lower := strings.ToLower(filename)

	// Find "no" prefix
	noIndex := strings.Index(lower, "no")
	if noIndex == -1 {
		return 1 // Default to volume 1
	}

	// Start after "no"
	numStart := noIndex + 2
	numEnd := numStart

	// Find the end of the numeric sequence
	for numEnd < len(filename) && filename[numEnd] >= '0' && filename[numEnd] <= '9' {
		numEnd++
	}

	if numEnd == numStart {
		return 1 // No number found, default to volume 1
	}

	// Parse the number
	var volumeNum int
	fmt.Sscanf(filename[numStart:numEnd], "%d", &volumeNum)

	if volumeNum == 0 {
		return 1 // Invalid number, default to volume 1
	}

	return volumeNum
}
