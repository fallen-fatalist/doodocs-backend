package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"path/filepath"
)

func main() {
	dirPath := "zipdir/data"
	if len(os.Args) == 2 {
		dirPath = os.Args[1]
	}

	err := createZip(dirPath)
	if err != nil {
		log.Fatal(err)
	}
}

type FileContent struct {
	FilePath string
	Reader   io.Reader
}

func collectFiles(dirPath string) ([]FileContent, error) {
	files := make([]FileContent, 0)
	// Walk the directory and add all files
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error walking path %s: %v", path, err)
		}

		// Skip directories, we only want to add files
		if info.IsDir() {
			return nil
		}

		// Open the file to be added
		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("error opening file %s: %v", path, err)
		}

		// Copy the file's content into the form field
		files = append(files, FileContent{
			FilePath: path,
			Reader:   file,
		})
		if err != nil {
			return fmt.Errorf("error copying content for file %s: %v", path, err)
		}

		return nil
	})

	return files, err
}

func createZip(dirPath string) error {

	files, err := collectFiles(dirPath)
	if err != nil {
		log.Fatal(err)
	}
	// Create a temporary file to store the ZIP archive
	archiveFile, err := os.OpenFile("data/sample.zip", os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		return fmt.Errorf("unable to create temporary file: %v", err)
	}
	defer archiveFile.Close()

	// Create a zip.Writer that writes to the temporary ZIP file
	zipWriter := zip.NewWriter(archiveFile)

	// Buffer for stream reading
	buf := make([]byte, 4096)

	for _, fileContent := range files {
		fileHeader := &zip.FileHeader{
			Name:   filepath.Base(fileContent.FilePath),
			Method: zip.Deflate,
		}

		// Create a new entry in the ZIP file using the file header
		writer, err := zipWriter.CreateHeader(fileHeader)
		if err != nil {
			return fmt.Errorf("error creating zip header for file %s: %v", fileContent.FilePath, err)
		}
		total := 0
		for {
			n, err := fileContent.Reader.Read(buf)
			total += n
			if err != nil && err != io.EOF {
				return err
			} else if n == 0 {
				slog.Info(fmt.Sprintf("total %d bytes written into temporary archive with file: %s", total, fileContent.FilePath))
				break
			}
			// Write the data to the destination file
			_, err = writer.Write(buf[:n])
			if err != nil {
				return err
			}
		}
	}
	zipWriter.Flush()

	// Return the temporary file as an io.Reader to allow streaming it out
	return nil
}
