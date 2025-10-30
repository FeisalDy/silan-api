package epub

// Types that represent OPF and EPUB content exported for consumers

// OPFMetadata contains book metadata
type OPFMetadata struct {
	Title       []string `xml:"title"`
	Creator     []string `xml:"creator"`
	Language    []string `xml:"language"`
	Publisher   []string `xml:"publisher"`
	Description []string `xml:"description"`
	Subject     []string `xml:"subject"`
	Date        []string `xml:"date"`
	Identifier  []string `xml:"identifier"`
	Rights      []string `xml:"rights"`
}

// OPFManifestItem represents an item in the manifest
type OPFManifestItem struct {
	ID        string `xml:"id,attr"`
	Href      string `xml:"href,attr"`
	MediaType string `xml:"media-type,attr"`
}

// OPFItemRef references a manifest item in reading order
type OPFItemRef struct {
	IDRef string `xml:"idref,attr"`
}

// OPFSpine defines the reading order
type OPFSpine struct {
	Toc      string       `xml:"toc,attr"`
	ItemRefs []OPFItemRef `xml:"itemref"`
}

// OPFPackage represents the parsed OPF file structure
type OPFPackage struct {
	Metadata OPFMetadata       `xml:"metadata"`
	Manifest []OPFManifestItem `xml:"manifest>item"`
	Spine    OPFSpine          `xml:"spine"`
}

// ContentFile represents a parsed content file (HTML/XHTML)
type ContentFile struct {
	Path      string
	RawHTML   string
	PlainText string
	MediaType string
}

// EpubContent represents the fully parsed EPUB with all extracted data
type EpubContent struct {
	Metadata     OPFMetadata
	Manifest     []OPFManifestItem
	Spine        []OPFItemRef
	ContentFiles map[string]ContentFile // Key is the file path
	RawFiles     map[string][]byte      // All raw files for reference
	OPFPath      string
}

// RawEpub is the minimal representation returned by EpubService.
// Transformers are responsible for parsing OPF/manifest/spine from this raw data.
type RawEpub struct {
	RawFiles map[string][]byte
	OPFPath  string
}
