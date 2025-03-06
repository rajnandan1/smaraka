package models

type FileUpload struct {
	FileName      string `json:"file_name"`
	FileSize      int64  `json:"file_size"`
	FileType      string `json:"file_type"`
	ImportType    string `json:"import_type"`
	FileExtension string `json:"file_extension"`
	FileContent   string `json:"file_content"`
}
type FileUploadResponse struct {
	Name    string `json:"name"`
	Icon    string `json:"icon"`
	URL     string `json:"url"`
	AddedOn string `json:"added_on"`
}
