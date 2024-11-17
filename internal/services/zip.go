package services

import (
	"io"
	"zip-api/internal/core/entities"
)

type ZipService interface {
	ZipInfo(zipArchive io.Reader) (*entities.Archive, error)
	ZipArchive(files []io.Reader) ([]byte, error)
}
