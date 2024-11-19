package services

import (
	"io"
	"mime/multipart"
	"os"
	"zip-api/internal/core/entities"
)

type ZipService interface {
	ZipInfo(zipArchive io.Reader, zipName string) (*entities.Archive, error)
	ZipArchive(fileParts []*multipart.Part) (*os.File, error)
}
