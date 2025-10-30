# EPUB Parsing Service - Enhanced Version

## Overview
The EPUB service has been significantly enhanced with proper OPF parsing, panic recovery, and structured data return.

## Key Features

### 1. **OPF File Discovery**
- Automatically locates the OPF file using `META-INF/container.xml`
- Handles various EPUB directory structures
- Validates container.xml structure

### 2. **Comprehensive OPF Parsing**
The service now parses all major OPF sections:

#### Metadata
- `Title` - Book title(s)
- `Creator` - Author(s)
- `Language` - Book language(s)
- `Publisher` - Publisher information
- `Description` - Book description
- `Subject` - Tags/genres/subjects
- `Date` - Publication date
- `Identifier` - ISBN or other identifiers
- `Rights` - Copyright information

#### Manifest
- List of all files in the EPUB with their:
  - ID
  - File path (href)
  - Media type (e.g., application/xhtml+xml, image/jpeg)

#### Spine
- Reading order of content files
- References to manifest items by ID
- Determines the correct chapter sequence

### 3. **Panic Recovery**
- Wraps the parsing logic in `parseEpubSafe()` with `defer/recover`
- Invalid ZIP files or corrupted archives won't crash the application
- Returns proper error messages instead of panicking

### 4. **Structured Data Return**

The service now returns an `EpubContent` struct instead of just logging:

```go
type EpubContent struct {
    Metadata     OPFMetadata           // Parsed metadata
    Manifest     []OPFManifestItem     // File manifest
    Spine        []OPFItemRef          // Reading order
    ContentFiles map[string]ContentFile // Parsed HTML/XHTML files
    RawFiles     map[string][]byte     // All raw file data
    OPFPath      string                // Path to OPF file
}
```

Each content file includes:
- `Path` - File location in EPUB
- `RawHTML` - Original HTML/XHTML content
- `PlainText` - Extracted text content (for search/preview)
- `MediaType` - MIME type

## API Response

### Success Response (200 OK)
```json
{
  "success": true,
  "message": "EPUB file parsed successfully. Check console for detailed output.",
  "data": {
    "novel_id": "123",
    "original_language": "en",
    "metadata": {
      "title": ["The Great Novel"],
      "authors": ["John Doe"],
      "language": ["en"],
      "publisher": ["Example Press"],
      "description": ["A wonderful story..."],
      "subjects": ["Fiction", "Adventure"],
      "date": ["2024"]
    },
    "manifest_items": 45,
    "spine_items": 30,
    "content_files": 28,
    "total_files": 45
  }
}
```

## Console Output

The service provides detailed console logs:

### 1. **Processing Information**
- Original language
- File size

### 2. **Metadata Summary**
- Title, Author(s), Language
- Publisher, Description
- Subjects/Tags, Date

### 3. **Manifest Summary**
- Total items count
- Breakdown by media type (HTML, CSS, images, etc.)

### 4. **Spine Summary**
- Total chapters/sections in reading order

### 5. **Content Files Preview**
- First 3 content files with:
  - File path
  - Media type
  - 200-character preview of plain text

## Usage in Future Database Integration

The returned `EpubContent` structure provides all the data needed for database integration:

### For Novel Table:
```go
novel := Novel{
    Title:            epubContent.Metadata.Title[0],
    OriginalAuthor:   strings.Join(epubContent.Metadata.Creator, ", "),
    OriginalLanguage: epubContent.Metadata.Language[0],
    Description:      strings.Join(epubContent.Metadata.Description, " "),
}
```

### For Tags/Genres:
```go
for _, subject := range epubContent.Metadata.Subject {
    // Create or link tags
}
```

### For Chapters:
```go
// Iterate through spine in order
for order, itemRef := range epubContent.Spine {
    // Find the manifest item
    for _, item := range epubContent.Manifest {
        if item.ID == itemRef.IDRef {
            // Get the content file
            contentFile := epubContent.ContentFiles[baseDir + item.Href]
            
            chapter := Chapter{
                OrderNum: order,
                Content:  contentFile.RawHTML,
                // or use PlainText for preview
            }
        }
    }
}
```

## Error Handling

The service handles various error scenarios:

1. **Empty File**: Returns error immediately
2. **Invalid ZIP**: Caught by panic recovery
3. **Missing container.xml**: Clear error message
4. **Missing OPF file**: Descriptive error
5. **Invalid XML**: XML parsing errors are caught and returned
6. **Missing content files**: Logged but doesn't stop processing

## Technical Implementation Details

### OPF Path Resolution
1. Read `META-INF/container.xml`
2. Parse XML to find `rootfile` element
3. Extract `full-path` attribute
4. Use this path to locate the OPF file

### Relative Path Handling
- Extracts base directory from OPF path
- Resolves all manifest hrefs relative to OPF location
- Handles EPUBs with various directory structures

### Text Extraction
- Parses HTML/XHTML using `golang.org/x/net/html`
- Recursively extracts text nodes
- Trims whitespace and joins with spaces
- Useful for search functionality and previews

## Dependencies

The enhanced service uses:
- `archive/zip` - ZIP archive handling
- `encoding/xml` - XML parsing for container and OPF
- `golang.org/x/net/html` - HTML parsing for text extraction
- `bytes`, `io`, `strings` - Standard library utilities

## Next Steps

With the structured data now available, you can:

1. **Implement Database Saving**
   - Map metadata to Novel table
   - Create Volume records (if EPUB has volume structure)
   - Create Chapter records from spine order
   - Store translations

2. **Add More Parsing Features**
   - Parse TOC (table of contents) NCX file
   - Extract and store images
   - Handle CSS stylesheets
   - Parse nav.xhtml for EPUB3

3. **Add Validation**
   - Validate required metadata fields
   - Check for duplicate chapters
   - Verify image references

4. **Enhance Content Processing**
   - Clean HTML (remove unnecessary tags)
   - Convert to markdown
   - Extract chapter titles from headings
   - Split long chapters by headings

## Testing

Example test with cURL:
```bash
curl -X POST "http://localhost:8080/api/v1/novels/123/epub" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -F "epub_file=@test.epub" \
  -F "original_language=en"
```

You should see:
- Immediate JSON response with parsed metadata
- Detailed console logs showing the EPUB structure
- No crashes even with invalid files
