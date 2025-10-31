package transformer

import (
	"context"

	"simple-go/pkg/epub"
)

type EpubSourceType string

const (
	EpubSource404NovelDownloader      EpubSourceType = "404_novel_downloader"
	EpubSourceDipubdLightnovelCrawler EpubSourceType = "dipubd_lightnovel_crawler"
	EpubSourceGeneric                 EpubSourceType = "generic"
)

type EpubTransformer interface {
	DetectSource(content *epub.RawEpub) bool
	TransformToNovelData(ctx context.Context, content *epub.RawEpub) (*NovelData, error)
	TransformToVolumes(ctx context.Context, content *epub.RawEpub) ([]VolumeData, error)
	TransformToChapters(ctx context.Context, content *epub.RawEpub) ([]ChapterData, error)
	GetSourceType() EpubSourceType
}

type NovelData struct {
	Title            string
	OriginalAuthor   string
	Description      string
	Publisher        string
	OriginalLanguage string
	Tags             []string
	CoverImage       []byte
}

type VolumeData struct {
	Number    int    // Volume number (1, 2, 3, etc.)
	Title     string // Volume title (optional)
	IsVirtual bool   // True if this is a virtual volume (when source has no volumes)
}

type ChapterData struct {
	VolumeIndex int // Index into the Volumes array (which volume this chapter belongs to)
	OrderNum    int // Chapter order within the volume
	Title       string
	Content     string
	PlainText   string
}

type EpubProcessResult struct {
	RawContent    *epub.RawEpub
	NovelData     *NovelData
	Volumes       []VolumeData  // List of volumes (at least one, even if virtual)
	Chapters      []ChapterData // Chapters with volume references
	SourceType    EpubSourceType
	TotalChapters int
	TotalVolumes  int
}
