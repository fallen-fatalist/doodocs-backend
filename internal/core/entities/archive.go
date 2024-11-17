package entities

type Archive struct {
	FileName   string `json: "filename"`
	Size       uint32 `json: "archive_size"`
	TotalSize  uint32 `json: "total_size"`
	TotalFiles uint32 `json: "total_files"`
	Files      []file `json: "files"`
}
