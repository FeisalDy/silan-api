package transformer

import (
	"context"
	"fmt"
	"simple-go/pkg/epub"
	"simple-go/pkg/logger"
	"strings"
)

type SourceDipubdLightnovelCrawlerTransformer struct{}

func NewSourceDipubdLightnovelCrawlerTransformer() *SourceDipubdLightnovelCrawlerTransformer {
	return &SourceDipubdLightnovelCrawlerTransformer{}
}

func (t *SourceDipubdLightnovelCrawlerTransformer) DetectSource(content *epub.RawEpub) bool {
	for path, file := range content.RawFiles {
		if strings.EqualFold(path, "EPUB/intro.xhtml") {
			data := string(file)
			if strings.Contains(data, "https://github.com/dipu-bd/lightnovel-crawler") {
				logger.Info("Detected Source B format (Lightnovel Crawler)")
				return true
			}
		}
	}

	return false
}

func (t *SourceDipubdLightnovelCrawlerTransformer) GetSourceType() EpubSourceType {
	return EpubSourceDipubdLightnovelCrawler
}

func (t *SourceDipubdLightnovelCrawlerTransformer) TransformToNovelData(ctx context.Context, content *epub.RawEpub) (*NovelData, error) {
	data := &NovelData{
		Tags: []string{},
	}

	// Parse OPF from raw files
	opfBytes, ok := content.RawFiles[content.OPFPath]
	if !ok {
		logger.Info("OPF not found in raw epub; returning best-effort metadata")
		return data, nil
	}

	opfPkg, err := epub.ParseOPF(opfBytes)
	if err != nil {
		logger.Error(err, "failed to parse OPF in transformer")
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
	}

	// Extract tags from subjects
	data.Tags = opfPkg.Metadata.Subject

	// Source B specific: Extract cover image
	for path, fileBytes := range content.RawFiles {
		lowerPath := strings.ToLower(path)
		if strings.Contains(lowerPath, "cover.jpg") ||
			strings.Contains(lowerPath, "cover.jpeg") ||
			strings.Contains(lowerPath, "cover.png") {
			data.CoverImage = fileBytes
			logger.Info(fmt.Sprintf("Extracted cover image from: %s", path))
			break
		}
	}

	// If no direct cover.jpg, try to find from manifest
	if data.CoverImage == nil {
		baseDir := getBaseDir(content.OPFPath)
		for _, item := range opfPkg.Manifest {
			if strings.Contains(strings.ToLower(item.ID), "cover") &&
				strings.Contains(item.MediaType, "image") {
				coverPath := baseDir + item.Href
				if coverBytes, exists := content.RawFiles[coverPath]; exists {
					data.CoverImage = coverBytes
					logger.Info(fmt.Sprintf("Extracted cover image from manifest: %s", coverPath))
					break
				}
			}
		}
	}

	// Source B doesn't have separate synopsis file
	logger.Info("Source B: Using description as synopsis")
	data.Synopsis = data.Description

	return data, nil
}

func (t *SourceDipubdLightnovelCrawlerTransformer) TransformToChapters(ctx context.Context, content *epub.RawEpub) ([]ChapterData, error) {
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
			logger.Error(nil, fmt.Sprintf("Manifest item not found for spine ref: %s", itemRef.IDRef))
			continue
		}

		// Skip cover page
		if strings.Contains(strings.ToLower(manifestItem.Href), "cover") {
			logger.Info(fmt.Sprintf("Skipping cover page in chapters: %s", manifestItem.Href))
			continue
		}

		// Only process HTML/XHTML files
		if !strings.Contains(manifestItem.MediaType, "html") {
			continue
		}

		fullPath := baseDir + manifestItem.Href
		contentFile, exists := contentFiles[fullPath]
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

	logger.Info(fmt.Sprintf("Source B: Extracted %d chapters", len(chapters)))
	return chapters, nil
}
