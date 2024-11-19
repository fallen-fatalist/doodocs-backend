package entities

import "io"

type FileMetadata struct {
	FilePath string `json:"file_path"`
	Size     uint32 `json:"size"`
	MimeType string `json:"mimetype"`
}

type FileContent struct {
	FilePath string
	Reader   io.Reader
}
