package transformer

import (
	"context"
	"fmt"
	"simple-go/pkg/epub"
	"simple-go/pkg/logger"
	"strings"
)

// GenericTransformer handles standard EPUB format (fallback)
type GenericTransformer struct{}

func NewGenericTransformer() *GenericTransformer {
	return &GenericTransformer{}
}

func (t *GenericTransformer) DetectSource(content *epub.EpubContent) bool {
	// Generic is the fallback - always returns true
	logger.Info("Using Generic EPUB transformer")
	return true
}

func (t *GenericTransformer) GetSourceType() EpubSourceType {
	return EpubSourceGeneric
}

func (t *GenericTransformer) TransformToNovelData(ctx context.Context, content *epub.EpubContent) (*NovelData, error) {
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
		data.Synopsis = data.Description
	}

	// Extract tags from subjects
	data.Tags = content.Metadata.Subject

	logger.Info("Generic transformer: Using standard EPUB metadata")
	return data, nil
}

func (t *GenericTransformer) TransformToChapters(ctx context.Context, content *epub.EpubContent) ([]ChapterData, error) {
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
			continue
		}

		// Only process HTML/XHTML files
		if !strings.Contains(manifestItem.MediaType, "html") {
			continue
		}

		fullPath := baseDir + manifestItem.Href
		contentFile, exists := content.ContentFiles[fullPath]
		if !exists {
			continue
		}

		chapters = append(chapters, ChapterData{
			OrderNum:  order + 1,
			Title:     fmt.Sprintf("Chapter %d", order+1),
			Content:   contentFile.RawHTML,
			PlainText: contentFile.PlainText,
		})
	}

	logger.Info(fmt.Sprintf("Generic transformer: Extracted %d chapters", len(chapters)))
	return chapters, nil
}
