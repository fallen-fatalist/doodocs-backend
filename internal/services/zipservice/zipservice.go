package zipservice

import (
	"archive/zip"
	"errors"
	"io"
	"net/http"
	"os"
	"zip-api/internal/core/entities"
	"zip-api/internal/utils"
)

// Errors
var (
	ErrIncorrectMimeType = errors.New("not allowed mimetype provided in archive file")
)

type zipService struct {
}

func NewZipService() *zipService {
	return &zipService{}
}

func (s *zipService) ZipInfo(archiveReader io.Reader, zipName string) (*entities.Archive, error) {
	zipFile, err := os.CreateTemp("", "*.zip")
	if err != nil {
		return nil, err
	}
	defer zipFile.Close()

	buf := make([]byte, 2048)
	var zipSize int64
	// Stream reading of the archive into the new file
	for {
		n, err := archiveReader.Read(buf)
		if err != nil && err != io.EOF {
			return nil, err
		} else if n == 0 {
			break
		}
		// Write the data to the destination file
		_, err = zipFile.Write(buf[:n])
		if err != nil {
			return nil, err
		}
		zipSize += int64(n)
	}

	// Read the ZIP archive from the temporary file
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
			if utils.ComplySingature(sniff, utils.DocxSequence) {
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

func (s *zipService) ZipArchive(files []io.Reader) (io.Reader, error) {
	return nil, nil
}
