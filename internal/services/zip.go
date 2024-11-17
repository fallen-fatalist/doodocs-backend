package services

import (
	"io"
	"zip-api/internal/core/entities"
)

type ZipService interface {
	ZipInfo(zipArchive io.Reader, zipName string) (*entities.Archive, error)
	ZipArchive(files []io.Reader) ([]byte, error)
}
