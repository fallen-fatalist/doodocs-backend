package zipservice

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
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

	buf := make([]byte, 4096)
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
			if utils.ComplySignature(sniff, utils.DocxSequence) {
				mimeType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
			}
		}

		// Extract file metadata
		file := entities.FileMetadata{
			FilePath: file.Name,
			Size:     uint32(file.UncompressedSize64),
			MimeType: mimeType,
		}

		// Append the file metadata to the archive
		archive.Files = append(archive.Files, file)
	}

	return archive, nil

}

func (s *zipService) ZipArchive(fileParts []*multipart.Part) (*os.File, error) {
	return nil, nil
	// Create a temporary file to store the ZIP archive
	tmpArchive, err := os.CreateTemp("", "*.zip")
	if err != nil {
		return nil, fmt.Errorf("unable to create temporary file: %v", err)
	}

	// Create a zip.Writer that writes to the temporary ZIP file
	zipWriter := zip.NewWriter(tmpArchive)

	// Buffer for stream reading
	buf := make([]byte, 4096)

	for _, filePart := range fileParts {
		fileHeader := &zip.FileHeader{
			Name:   filepath.Base(filePart.FileName()),
			Method: zip.Deflate,
		}
		// Create a new entry in the ZIP file using the file header
		fileWriter, err := zipWriter.CreateHeader(fileHeader)
		if err != nil {
			return nil, fmt.Errorf("error creating zip header for file %s: %v", filePart.FileName(), err)
		}
		total := 0

		for {
			n, err := filePart.Read(buf)
			total += n
			if err != nil && err != io.EOF {
				return nil, err
			} else if n == 0 && err == io.EOF {
				slog.Info(fmt.Sprintf("total %d bytes written into temporary archive with file: %s", total, filePart.FileName()))
				break
			}
			// Write the data to the destination file
			_, err = fileWriter.Write(buf[:n])
			if err != nil {
				return nil, err
			}
		}

	}
	zipWriter.Close()
	_, err = tmpArchive.Seek(0, io.SeekStart)
	if err != nil {
		return nil, fmt.Errorf("error seeking to the start of the temporary file: %v", err)
	}

	// Return the temporary file as an io.Reader to allow streaming it out
	return tmpArchive, nil
}
