package media

type UploadAndSaveDTO struct {
	// One of the following sources must be provided
	FileBytes []byte
	Base64    string
	URL       string

	// Optional metadata
	Name        string
	UploaderID  string
	Description *string
	MimeType    *string
	FileSize    *int64
	TTL         uint64
}
