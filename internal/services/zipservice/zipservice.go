package zipservice

import (
	"archive/zip"
	"errors"
	"io"
	"net/http"
	"os"
	"zip-api/internal/core/entities"
)

// Errors
var (
	ErrIncorrectMimeType = errors.New("not allowed mimetype provided in archive file")
)

var (
	AllowedMimeTypes = []string{
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		"application/xml",
		"text/xml",
		"image/jpeg",
		"image/png",
	}

	docxSequence = []byte{80, 75, 3, 4}
)

type zipService struct {
}

func NewZipService() *zipService {
	return &zipService{}
}

func (s *zipService) ZipInfo(zipArchiveBinaryReader io.Reader, zipName string) (*entities.Archive, error) {
	zipFile, err := os.CreateTemp("", "*.zip")
	if err != nil {
		return nil, err
	}

	zipSize, err := io.Copy(zipFile, zipArchiveBinaryReader)
	if err != nil {
		return nil, err
	}
	defer zipFile.Close()

	// Read the ZIP archive from the io.Reader
	zipReader, err := zip.NewReader(zipFile, zipSize)
	if err != nil {
		return nil, err
	}

	// Prepare the Archive struct
	archive := &entities.Archive{
		FileName:   zipName,
		Size:       uint32(zipSize),             // size of the ZIP file in bytes
		TotalSize:  0,                           // Total uncompressed size of files in the archive
		TotalFiles: uint32(len(zipReader.File)), // Total number of files
	}

	sniff := make([]byte, 512)

	// Iterate through each file in the ZIP archive
	for _, file := range zipReader.File {
		fileReader, err := file.Open()
		if err != nil {
			return nil, err
		}
		defer fileReader.Close()

		n, err := fileReader.Read(sniff)
		if err != nil && err != io.EOF {
			return nil, err
		}

		// Update the total uncompressed size
		archive.TotalSize += uint32(file.UncompressedSize64)

		mimeType := http.DetectContentType(sniff[:n])
		if mimeType == "text/xml; charset=utf-8" {
			mimeType = "application/xml"
		} else if mimeType == "application/zip" {
			if complySingature(sniff, docxSequence) {
				mimeType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
			}
		}

		// Extract file metadata
		file := entities.File{
			FilePath: file.Name,
			Size:     uint32(file.UncompressedSize64),
			MimeType: mimeType,
		}

		// Append the file metadata to the archive
		archive.Files = append(archive.Files, file)
	}

	return archive, nil

}

func (s *zipService) ZipArchive(files []io.Reader) ([]byte, error) {
	return []byte{}, nil
}

func In(s string, strs []string) bool {
	for _, str := range strs {
		if str == s {
			return true
		}
	}
	return false
}

func complySingature(sniff []byte, signature []byte) bool {
	for idx, _ := range signature {
		if idx >= len(sniff) {
			return false
		} else if sniff[idx] != signature[idx] {
			return false
		}
	}
	return true
}
