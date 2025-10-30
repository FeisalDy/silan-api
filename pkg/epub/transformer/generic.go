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

func (t *GenericTransformer) DetectSource(content *epub.RawEpub) bool {
	// Generic is the fallback - always returns true
	logger.Info("Using Generic EPUB transformer")
	// return true
	return false
}

func (t *GenericTransformer) GetSourceType() EpubSourceType {
	return EpubSourceGeneric
}

func (t *GenericTransformer) TransformToNovelData(ctx context.Context, content *epub.RawEpub) (*NovelData, error) {
	data := &NovelData{
		Tags: []string{},
	}

	// Parse OPF from raw files
	opfBytes, ok := content.RawFiles[content.OPFPath]
	if !ok {
		logger.Info("OPF not found in raw epub; returning empty metadata")
		return data, nil
	}

	opfPkg, err := epub.ParseOPF(opfBytes)
	if err != nil {
		logger.Error(err, "failed to parse OPF in generic transformer")
		return data, nil
	}

	// Extract basic metadata
	if len(opfPkg.Metadata.Title) > 0 {
		data.Title = opfPkg.Metadata.Title[0]
	}
	if len(opfPkg.Metadata.Creator) > 0 {
		data.OriginalAuthor = strings.Join(opfPkg.Metadata.Creator, ", ")
	}
	if len(opfPkg.Metadata.Language) > 0 {
		data.OriginalLanguage = opfPkg.Metadata.Language[0]
	}
	if len(opfPkg.Metadata.Publisher) > 0 {
		data.Publisher = opfPkg.Metadata.Publisher[0]
	}
	if len(opfPkg.Metadata.Description) > 0 {
		data.Description = strings.Join(opfPkg.Metadata.Description, " ")
		data.Synopsis = data.Description
	}

	// Extract tags from subjects
	data.Tags = opfPkg.Metadata.Subject

	logger.Info("Generic transformer: Using standard EPUB metadata")
	return data, nil
}

func (t *GenericTransformer) TransformToChapters(ctx context.Context, content *epub.RawEpub) ([]ChapterData, error) {
	chapters := []ChapterData{}

	// Parse OPF from raw files
	opfBytes, ok := content.RawFiles[content.OPFPath]
	if !ok {
		return chapters, nil
	}

	opfPkg, err := epub.ParseOPF(opfBytes)
	if err != nil {
		logger.Error(err, "failed to parse OPF in TransformToChapters")
		return chapters, nil
	}

	baseDir := getBaseDir(content.OPFPath)

	// Create manifest lookup map
	manifestMap := make(map[string]epub.OPFManifestItem)
	for _, item := range opfPkg.Manifest {
		manifestMap[item.ID] = item
	}

	// Build content files map
	contentFiles := make(map[string]epub.ContentFile)
	for _, item := range opfPkg.Manifest {
		if strings.Contains(item.MediaType, "html") || strings.Contains(item.MediaType, "xhtml") {
			fullPath := baseDir + item.Href
			if raw, exists := content.RawFiles[fullPath]; exists {
				contentFiles[fullPath] = epub.ContentFile{
					Path:      fullPath,
					RawHTML:   string(raw),
					PlainText: epub.ExtractText(raw),
					MediaType: item.MediaType,
				}
			}
		}
	}

	// Iterate through spine in order
	for order, itemRef := range opfPkg.Spine.ItemRefs {
		manifestItem, exists := manifestMap[itemRef.IDRef]
		if !exists {
			continue
		}

		// Only process HTML/XHTML files
		if !strings.Contains(manifestItem.MediaType, "html") {
			continue
		}

		fullPath := baseDir + manifestItem.Href
		contentFile, exists := contentFiles[fullPath]
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
