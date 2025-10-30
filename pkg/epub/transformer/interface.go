package transformer

import (
	"context"
	"simple-go/pkg/epub"
)

// EpubSourceType represents different EPUB sources with different formats
type EpubSourceType string

const (
	EpubSource404NovelDownloader      EpubSourceType = "source_a"
	EpubSourceDipubdLightnovelCrawler EpubSourceType = "source_b"
	EpubSourceGeneric                 EpubSourceType = "generic"
)

// EpubTransformer defines the interface for transforming EPUB content to database models
type EpubTransformer interface {
	DetectSource(content *epub.EpubContent) bool
	TransformToNovelData(ctx context.Context, content *epub.EpubContent) (*NovelData, error)
	TransformToChapters(ctx context.Context, content *epub.EpubContent) ([]ChapterData, error)
	GetSourceType() EpubSourceType
}

// NovelData represents extracted novel information ready for database
type NovelData struct {
	Title            string
	OriginalAuthor   string
	Description      string
	Publisher        string
	OriginalLanguage string
	Tags             []string
	CoverImage       []byte
	Synopsis         string
}

// ChapterData represents extracted chapter information
type ChapterData struct {
	OrderNum  int
	Title     string
	Content   string
	PlainText string
}

// EpubProcessResult contains all processed EPUB data ready for database insertion
type EpubProcessResult struct {
	RawContent    *epub.EpubContent
	NovelData     *NovelData
	Chapters      []ChapterData
	SourceType    EpubSourceType
	TotalChapters int
}
