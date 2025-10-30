package transformer

import (
	"context"
	"fmt"
	"simple-go/pkg/epub"
	"simple-go/pkg/logger"
	"strings"
)

type Source404NovelDownloaderTransformer struct{}

func NewSource404NovelDownloaderTransformer() *Source404NovelDownloaderTransformer {
	return &Source404NovelDownloaderTransformer{}
}

func (t *Source404NovelDownloaderTransformer) DetectSource(content *epub.EpubContent) bool {
	const targetFile = "oebps/info.txt"
	const markerText = "https://github.com/404-novel-project/novel-downloader"

	// Normalize file names to lowercase for consistent matching
	for path, data := range content.RawFiles {
		if strings.ToLower(path) == targetFile {
			text := string(data)
			if strings.Contains(text, markerText) {
				logger.Info("Detected Source A format (info.txt contains novel-downloader marker)")
				return true
			}
			logger.Info("Found info.txt but no marker text inside")
			return false
		}
	}

	return false
}

func (t *Source404NovelDownloaderTransformer) GetSourceType() EpubSourceType {
	return EpubSource404NovelDownloader
}

func (t *Source404NovelDownloaderTransformer) TransformToNovelData(ctx context.Context, content *epub.EpubContent) (*NovelData, error) {
	data := &NovelData{
		Tags: []string{},
	}

	// Extract basic metadata
	if len(content.Metadata.Title) > 0 {
		data.Title = content.Metadata.Title[0]
	}
	if len(content.Metadata.Creator) > 0 {
		data.OriginalAuthor = strings.Join(content.Metadata.Creator, ", ")
	}
	if len(content.Metadata.Language) > 0 {
		data.OriginalLanguage = content.Metadata.Language[0]
	}
	if len(content.Metadata.Publisher) > 0 {
		data.Publisher = content.Metadata.Publisher[0]
	}
	if len(content.Metadata.Description) > 0 {
		data.Description = strings.Join(content.Metadata.Description, " ")
	}

	// Extract tags from subjects
	data.Tags = content.Metadata.Subject

	// Source A specific: Extract synopsis from synopsis.xhtml
	for path, file := range content.ContentFiles {
		if strings.Contains(strings.ToLower(path), "synopsis.xhtml") {
			data.Synopsis = file.PlainText
			logger.Info("Extracted synopsis from synopsis.xhtml")
			break
		}
	}

	// Source A doesn't have cover image
	logger.Info("Source A: No cover image available")

	return data, nil
}

func (t *Source404NovelDownloaderTransformer) TransformToChapters(ctx context.Context, content *epub.EpubContent) ([]ChapterData, error) {
	chapters := []ChapterData{}
	baseDir := getBaseDir(content.OPFPath)

	// Create manifest lookup map
	manifestMap := make(map[string]epub.OPFManifestItem)
	for _, item := range content.Manifest {
		manifestMap[item.ID] = item
	}

	// Iterate through spine in order
	for order, itemRef := range content.Spine {
		manifestItem, exists := manifestMap[itemRef.IDRef]
		if !exists {
			logger.Error(nil, fmt.Sprintf("Manifest item not found for spine ref: %s", itemRef.IDRef))
			continue
		}

		// Skip synopsis file - it's not a chapter
		if strings.Contains(strings.ToLower(manifestItem.Href), "synopsis") {
			logger.Info(fmt.Sprintf("Skipping synopsis file in chapters: %s", manifestItem.Href))
			continue
		}

		// Only process HTML/XHTML files
		if !strings.Contains(manifestItem.MediaType, "html") {
			continue
		}

		fullPath := baseDir + manifestItem.Href
		contentFile, exists := content.ContentFiles[fullPath]
		if !exists {
			logger.Error(nil, fmt.Sprintf("Content file not found: %s", fullPath))
			continue
		}

		// Extract chapter title from content or use order number
		chapterTitle := extractChapterTitle(contentFile.RawHTML, order+1)

		chapters = append(chapters, ChapterData{
			OrderNum:  order + 1,
			Title:     chapterTitle,
			Content:   contentFile.RawHTML,
			PlainText: contentFile.PlainText,
		})
	}

	logger.Info(fmt.Sprintf("Source A: Extracted %d chapters", len(chapters)))
	return chapters, nil
}
