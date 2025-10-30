package transformer

import (
	"context"

	"simple-go/pkg/epub"
)

// EpubSourceType represents different EPUB sources with different formats
type EpubSourceType string

const (
	EpubSource404NovelDownloader      EpubSourceType = "404_novel_downloader"
	EpubSourceDipubdLightnovelCrawler EpubSourceType = "dipubd_lightnovel_crawler"
	EpubSourceGeneric                 EpubSourceType = "generic"
)

// EpubTransformer defines the interface for transforming EPUB content to database models
type EpubTransformer interface {
	// DetectSource should look at raw EPUB files and determine whether this transformer
	// can handle the source. Transformers are responsible for parsing OPF/manifest/spine.
	DetectSource(content *epub.RawEpub) bool
	TransformToNovelData(ctx context.Context, content *epub.RawEpub) (*NovelData, error)
	TransformToChapters(ctx context.Context, content *epub.RawEpub) ([]ChapterData, error)
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
	// RawContent contains the original raw files from the uploaded EPUB. Transformers
	// will parse OPF/manifest/spine and produce NovelData/Chapters.
	RawContent    *epub.RawEpub
	NovelData     *NovelData
	Chapters      []ChapterData
	SourceType    EpubSourceType
	TotalChapters int
}
