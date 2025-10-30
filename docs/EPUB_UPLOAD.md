# EPUB Upload Feature

## Overview
This feature allows users to upload EPUB files to novels. The system will parse the EPUB file and print its contents to the console for review before implementing the database saving logic.

## Endpoint

```
POST {{base_url}}/novels/:id/epub
```

### Authentication
Requires JWT authentication and `novel:update` permission.

### Request Parameters

#### Path Parameters
- `id` (string, required): The ID of the novel

#### Form Data
- `epub_file` (file, required): The EPUB file to upload
- `original_language` (text, required): The original language of the novel (e.g., "en", "id", "jp")

### Example Request

Using cURL:
```bash
curl -X POST "http://localhost:8080/api/v1/novels/123/epub" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -F "epub_file=@/path/to/your/novel.epub" \
  -F "original_language=en"
```

Using Postman:
1. Set method to POST
2. URL: `{{base_url}}/novels/:id/epub`
3. Headers: Add `Authorization: Bearer YOUR_JWT_TOKEN`
4. Body: Select `form-data`
   - Key: `epub_file`, Type: File, Value: Select your EPUB file
   - Key: `original_language`, Type: Text, Value: `en` (or your language code)

### Response

#### Success Response (200 OK)
```json
{
  "success": true,
  "message": "EPUB file uploaded successfully for novel 123. Check console for parsed content.",
  "data": null
}
```

#### Error Responses

**Missing EPUB file (400 Bad Request)**
```json
{
  "success": false,
  "message": "Missing epub_file in form data"
}
```

**Missing original language (400 Bad Request)**
```json
{
  "success": false,
  "message": "Missing original_language in form data"
}
```

**Unauthorized (401 Unauthorized)**
```json
{
  "success": false,
  "message": "Unauthorized"
}
```

**Processing Error (500 Internal Server Error)**
```json
{
  "success": false,
  "message": "Failed to process epub file: [error details]"
}
```

## Console Output

When an EPUB file is uploaded, the service will print detailed information to the console:

1. **File Information**: File size and original language
2. **File List**: All files contained in the EPUB archive
3. **Metadata**: Content from content.opf file (book metadata)
4. **Content Files**: HTML/XHTML files with truncated previews
5. **Table of Contents**: Content from toc.ncx file

Example console output:
```
[INFO] Processing EPUB file
[INFO] Original Language: en
[INFO] File Size: 524288 bytes
[INFO] ========== EPUB CONTENTS ==========
[INFO] Total files found: 25
[INFO] 
[INFO] ========== FILE LIST ==========
[INFO]   - mimetype
[INFO]   - META-INF/container.xml
[INFO]   - OEBPS/content.opf
[INFO]   - OEBPS/toc.ncx
[INFO]   - OEBPS/Text/chapter1.xhtml
[INFO]   ...
```

## Implementation Details

### Files Modified/Created

1. **`internal/service/epub_service.go`** (NEW)
   - `EpubService`: Main service for handling EPUB files
   - `UploadAndParseEpub()`: Processes the EPUB and prints contents
   - `parseEpub()`: Extracts files from EPUB (ZIP) format
   - `printEpubContents()`: Formats and prints EPUB structure

2. **`internal/handler/novel_handler.go`** (MODIFIED)
   - Added `epubService` field to `NovelHandler`
   - Added `UploadEpub()` handler method

3. **`internal/app/app.go`** (MODIFIED)
   - Initialize `EpubService`
   - Pass `epubService` to `NovelHandler`

4. **`internal/server/gin/gin_server.go`** (MODIFIED)
   - Added route: `POST /:id/epub`

### Technical Notes

- EPUB files are parsed as ZIP archives (EPUB is based on ZIP format)
- The service extracts and reads all files in the EPUB
- Content files (HTML/XHTML) are truncated to 500 characters in console output
- No database operations are performed yet - this is just for reviewing the parsed content

## Next Steps

After reviewing the console output and deciding on the parsing strategy, you can:

1. Implement proper EPUB parsing to extract:
   - Novel metadata (title, author, description)
   - Chapter structure
   - Chapter content
   - Images and other media

2. Map the parsed data to your database schema:
   - Create/update Novel record
   - Create Volume records
   - Create Chapter records with translations
   - Store content and metadata

3. Add error handling for malformed EPUB files
4. Implement progress tracking for large files
5. Add validation for EPUB structure and content
