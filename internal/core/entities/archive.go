package entities

type Archive struct {
	fileName   string `json: "filename"`
	size       uint32 `json: "archive_size"`
	totalSize  uint32 `json: "total_size"`
	totalFiles uint32 `json: "total_files"`
	files      []File `json: "files"`
}
