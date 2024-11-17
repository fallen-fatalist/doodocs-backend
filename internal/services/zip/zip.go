package zip

import "zip-api/internal/core/entities"

type zipService struct {
}

func NewZipService() *zipService {
	return &zipService{}
}

func (s *zipService) ZipInfo(zipArchive []byte) (*entities.Archive, error) {
	return &entities.Archive{}, nil
}

func (s *zipService) ZipArchive(files []byte) ([]byte, error) {
	return []byte{}, nil
}
