package services

import (
	"zip-api/internal/core/entities"
)

type ZipService interface {
	ZipInfo(zipArchive []byte) (*entities.Archive, error)
	ZipArchive(files []byte) ([]byte, error)
}
