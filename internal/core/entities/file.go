package entities

type File struct {
	filePath string `json:"file_path"`
	size     uint32 `json: "size"`
	mimeType string `json: "mimetype"`
}
