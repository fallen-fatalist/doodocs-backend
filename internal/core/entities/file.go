package entities

type file struct {
	FilePath string `json:"file_path"`
	Size     uint32 `json: "size"`
	MimeType string `json: "mimetype"`
}
