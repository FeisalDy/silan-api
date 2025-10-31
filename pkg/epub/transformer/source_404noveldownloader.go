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

	if len(opfPkg.Metadata.Title) > 0 {
		data.Title = opfPkg.Metadata.Title[0]
	}
	if len(opfPkg.Metadata.Creator) > 0 {
		data.OriginalAuthor = opfPkg.Metadata.Creator[0]
	}
	if len(opfPkg.Metadata.Language) > 0 {
		data.OriginalLanguage = opfPkg.Metadata.Language[0]
	}
	if len(opfPkg.Metadata.Description) > 0 {
		data.Description = strings.Join(opfPkg.Metadata.Description, " ")
	}

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

	return data, nil
}

func (t *Source404NovelDownloaderTransformer) TransformToVolumes(ctx context.Context, content *epub.RawEpub) ([]VolumeData, error) {
	volumes := []VolumeData{}

	opfBytes, ok := content.RawFiles[content.OPFPath]
	if !ok {
		// No OPF found, create a single virtual volume
		volumes = append(volumes, VolumeData{
			Number:    1,
			Title:     "Volume 1",
			IsVirtual: true,
		})
		return volumes, nil
	}

	opfPkg, err := epub.ParseOPF(opfBytes)
	if err != nil {
		logger.Error(err, "failed to parse OPF in TransformToVolumes")
		// Create a single virtual volume as fallback
		volumes = append(volumes, VolumeData{
			Number:    1,
			Title:     "Volume 1",
			IsVirtual: true,
		})
		return volumes, nil
	}

	baseDir := getBaseDir(content.OPFPath)

	// Create manifest lookup map
	manifestMap := make(map[string]epub.OPFManifestItem)
	for _, item := range opfPkg.Manifest {
		manifestMap[item.ID] = item
	}

	// Build content files map for section files
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

	// Track volumes we've seen
	volumeMap := make(map[int]string) // volume number -> title

	// Look for Section files in spine order
	for _, itemRef := range opfPkg.Spine.ItemRefs {
		manifestItem, exists := manifestMap[itemRef.IDRef]
		if !exists {
			continue
		}

		href := strings.ToLower(manifestItem.Href)

		// Check if this is a section (volume) file
		if strings.Contains(href, "section.xhtml") {
			volumeNum := extractVolumeNumberFromFilename(manifestItem.Href)

			// Get the volume title from the section file content
			fullPath := baseDir + manifestItem.Href
			volumeTitle := fmt.Sprintf("Volume %d", volumeNum)

			if contentFile, exists := contentFiles[fullPath]; exists {
				// Try to extract title from the section file
				extractedTitle := extractChapterTitle(contentFile.RawHTML, volumeNum)
				if extractedTitle != "" && !strings.Contains(extractedTitle, "Chapter") {
					volumeTitle = extractedTitle
				}
			}

			volumeMap[volumeNum] = volumeTitle
		}
	}

	// If no sections found, check if there are any chapters to determine volumes
	if len(volumeMap) == 0 {
		// Scan for chapter files to detect volumes
		for _, itemRef := range opfPkg.Spine.ItemRefs {
			manifestItem, exists := manifestMap[itemRef.IDRef]
			if !exists {
				continue
			}

			href := strings.ToLower(manifestItem.Href)
			if strings.Contains(href, "chapter.xhtml") {
				volumeNum := extractVolumeNumberFromFilename(manifestItem.Href)
				if _, exists := volumeMap[volumeNum]; !exists {
					volumeMap[volumeNum] = fmt.Sprintf("Volume %d", volumeNum)
				}
			}
		}
	}

	// If still no volumes found, create a single virtual volume
	if len(volumeMap) == 0 {
		volumes = append(volumes, VolumeData{
			Number:    1,
			Title:     "Volume 1",
			IsVirtual: true,
		})
		logger.Info("No volumes detected, created virtual volume")
		return volumes, nil
	}

	// Convert map to sorted slice
	maxVolumeNum := 0
	for num := range volumeMap {
		if num > maxVolumeNum {
			maxVolumeNum = num
		}
	}

	// Create volumes in order
	for i := 1; i <= maxVolumeNum; i++ {
		title, exists := volumeMap[i]
		if !exists {
			title = fmt.Sprintf("Volume %d", i)
		}
		volumes = append(volumes, VolumeData{
			Number:    i,
			Title:     title,
			IsVirtual: false,
		})
	}

	logger.Info(fmt.Sprintf("Source 404: Extracted %d volumes", len(volumes)))
	return volumes, nil
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

	// Parse volumes and chapters structure
	// In 404 source: No00001Section.xhtml = Volume 1, No00001Chapter.xhtml = Chapter in Volume 1
	volumeMap := make(map[int]int) // maps volume number to volume index
	currentVolumeIndex := -1
	chapterOrder := 1

	// Iterate through spine in order
	for _, itemRef := range opfPkg.Spine.ItemRefs {
		manifestItem, exists := manifestMap[itemRef.IDRef]
		if !exists {
			logger.Error(nil, fmt.Sprintf("Manifest item not found for spine ref: %s", itemRef.IDRef))
			continue
		}

		href := strings.ToLower(manifestItem.Href)

		// Skip synopsis file
		if strings.Contains(href, "synopsis") {
			logger.Info(fmt.Sprintf("Skipping synopsis file: %s", manifestItem.Href))
			continue
		}

		// Check if this is a section (volume) file
		if strings.Contains(href, "section.xhtml") {
			// This marks a new volume - we'll handle it when we encounter chapters
			logger.Info(fmt.Sprintf("Found volume marker: %s", manifestItem.Href))
			continue
		}

		// Check if this is a chapter file
		if strings.Contains(href, "chapter.xhtml") {
			// Extract volume number from filename (e.g., No00001Chapter.xhtml -> volume 1)
			volumeNum := extractVolumeNumberFromFilename(manifestItem.Href)

			// Check if we need to track a new volume
			if _, exists := volumeMap[volumeNum]; !exists {
				currentVolumeIndex++
				volumeMap[volumeNum] = currentVolumeIndex
				logger.Info(fmt.Sprintf("Detected new volume: %d (index: %d)", volumeNum, currentVolumeIndex))
			}

			volumeIndex := volumeMap[volumeNum]

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

			// Extract chapter title from content
			chapterTitle := extractChapterTitle(contentFile.RawHTML, chapterOrder)

			chapters = append(chapters, ChapterData{
				VolumeIndex: volumeIndex,
				OrderNum:    chapterOrder,
				Title:       chapterTitle,
				Content:     contentFile.RawHTML,
				PlainText:   contentFile.PlainText,
			})

			chapterOrder++
		}
	}

	logger.Info(fmt.Sprintf("Source 404: Extracted %d chapters across %d volumes", len(chapters), len(volumeMap)))
	return chapters, nil
}
