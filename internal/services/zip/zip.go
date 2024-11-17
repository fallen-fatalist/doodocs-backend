package zip

import (
	"io"
	"os"
	"zip-api/internal/core/entities"
)

type zipService struct {
}

func NewZipService() *zipService {
	return &zipService{}
}

func (s *zipService) ZipInfo(zipArchive io.Reader) (*entities.Archive, error) {
	tempFile, err := os.CreateTemp("", "")
	if err != nil {
		return nil, err
	}
	defer tempFile.Close()

	
	return &entities.Archive{}, nil
}

func (s *zipService) ZipArchive(files []io.Reader) ([]byte, error) {
	return []byte{}, nil
}
