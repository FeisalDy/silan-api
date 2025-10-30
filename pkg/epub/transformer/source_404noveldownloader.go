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

func (t *Source404NovelDownloaderTransformer) DetectSource(content *epub.RawEpub) bool {
	const targetFile = "oebps/info.txt"
	const markerText = "https://github.com/404-novel-project/novel-downloader"

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

func (t *Source404NovelDownloaderTransformer) TransformToNovelData(ctx context.Context, content *epub.RawEpub) (*NovelData, error) {
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

	// Build content files map for HTML/XHTML
	baseDir := getBaseDir(content.OPFPath)
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

	// Source A specific: Extract synopsis from synopsis.xhtml
	for path, file := range contentFiles {
		if strings.Contains(strings.ToLower(path), "synopsis.xhtml") {
			data.Synopsis = file.PlainText
			logger.Info("Extracted synopsis from synopsis.xhtml")
			break
		}
	}

	// Source A doesn't have cover image by default
	logger.Info("Source A: No cover image available")

	return data, nil
}

func (t *Source404NovelDownloaderTransformer) TransformToChapters(ctx context.Context, content *epub.RawEpub) ([]ChapterData, error) {
	chapters := []ChapterData{}

	// Parse OPF and build manifest/content map
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

	logger.Info(fmt.Sprintf("Source A: Extracted %d chapters", len(chapters)))
	return chapters, nil
}
