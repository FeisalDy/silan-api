package service

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"simple-go/pkg/epub"
	"simple-go/pkg/logger"
	"strings"
)

type EpubService struct{}

func NewEpubService() *EpubService { return &EpubService{} }

func (s *EpubService) UploadAndExtractRawEpub(ctx context.Context, fileBytes []byte) (*epub.RawEpub, error) {
	if len(fileBytes) == 0 {
		return nil, errors.New("epub file is empty")
	}

	epubContent, err := s.parseEpubSafe(fileBytes)
	if err != nil {
		logger.Error(err, "failed to parse EPUB file")
		return nil, err
	}

	return &epub.RawEpub{RawFiles: epubContent.RawFiles, OPFPath: epubContent.OPFPath}, nil
}

func (s *EpubService) parseEpubSafe(fileBytes []byte) (content *epub.EpubContent, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic during EPUB parsing: %v", r)
			logger.Error(err, "Recovered from panic")
		}
	}()

	return s.parseEpub(fileBytes)
}

// parseEpub extracts and parses all EPUB content
func (s *EpubService) parseEpub(fileBytes []byte) (*epub.EpubContent, error) {
	reader, err := zip.NewReader(bytes.NewReader(fileBytes), int64(len(fileBytes)))
	if err != nil {
		return nil, errors.New("failed to open epub as zip")
	}

	rawFiles := make(map[string][]byte)
	for _, file := range reader.File {
		if file.FileInfo().IsDir() {
			continue
		}

		rc, err := file.Open()
		if err != nil {
			logger.Error(err, fmt.Sprintf("Failed to open file: %s", file.Name))
			continue
		}

		content, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			logger.Error(err, fmt.Sprintf("Failed to read file: %s", file.Name))
			continue
		}

		rawFiles[file.Name] = content
	}

	if len(rawFiles) == 0 {
		return nil, errors.New("no files found in epub")
	}

	opfPath, err := epub.FindOpfPath(rawFiles)
	if err != nil {
		return nil, errors.New("failed to find OPF file")
	}

	opfContent, ok := rawFiles[opfPath]
	if !ok {
		return nil, errors.New("OPF file not found")
	}

	opfPackage, err := epub.ParseOPF(opfContent)
	if err != nil {
		return nil, fmt.Errorf("failed to parse OPF: %w", err)
	}

	baseDir := ""
	if idx := strings.LastIndex(opfPath, "/"); idx >= 0 {
		baseDir = opfPath[:idx+1]
	}

	// Step 6: Parse content files based on manifest
	contentFiles := make(map[string]epub.ContentFile)
	for _, item := range opfPackage.Manifest {
		// Only process HTML/XHTML content files
		if strings.Contains(item.MediaType, "html") || strings.Contains(item.MediaType, "xhtml") {
			fullPath := baseDir + item.Href

			if rawContent, exists := rawFiles[fullPath]; exists {
				plainText := epub.ExtractText(rawContent)

				contentFiles[fullPath] = epub.ContentFile{
					Path:      fullPath,
					RawHTML:   string(rawContent),
					PlainText: plainText,
					MediaType: item.MediaType,
				}
			} else {
				logger.Error(nil, fmt.Sprintf("Content file not found: %s", fullPath))
			}
		}
	}

	return &epub.EpubContent{
		Metadata:     opfPackage.Metadata,
		Manifest:     opfPackage.Manifest,
		Spine:        opfPackage.Spine.ItemRefs,
		ContentFiles: contentFiles,
		RawFiles:     rawFiles,
		OPFPath:      opfPath,
	}, nil
}
