package main

import (
	"fmt"
	"io"
	"os"
)

type Archive struct {
	FileName   string `json:"filename"`
	Size       uint32 `json:"archive_size"`
	TotalSize  uint32 `json:"total_size"`
	TotalFiles uint32 `json:"total_files"`
	Files      []File `json:"files"`
}

type File struct {
	FilePath string `json:"file_path"`
	Size     uint32 `json:"size"`
	MimeType string `json:"mimetype"`
}

var (
	AllowedMimeTypes = []string{
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document",
		"application/xml",
		"image/jpeg",
		"image/png",
	}

	docxSequence = []byte{80, 75, 3, 4}
)

func main() {
	// Open the source zip file for reading
	zipFile, err := os.Open("data/dummy.zip")
	if err != nil {
		panic(err)
	}
	defer zipFile.Close()

	// Create the destination zip file for writing
	fileToWrite, err := os.Create("data/dummy_dummy.zip")
	if err != nil {
		panic(err)
	}
	defer fileToWrite.Close()

	// Create a buffer for reading the file in chunks
	buf := make([]byte, 2048)
	var zipSize int64

	// Stream reading of the archive into the new file
	for {
		n, err := zipFile.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n == 0 {
			break
		}
		// Write the data to the destination file
		_, err = fileToWrite.Write(buf[:n])
		if err != nil {
			panic(err)
		}
		zipSize += int64(n)
	}

	// Print the size of the written zip file
	fmt.Println("Written zip file with size:", zipSize)
}
