package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// Function to add files from a directory to the multipart request
func addFilesToMultipart(writer *multipart.Writer, dirPath string) error {
	// Walk the directory and add all files
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error walking path %s: %v", path, err)
		}

		// Skip directories, we only want to add files
		if info.IsDir() {
			return nil
		}

		// Create form file for each file
		filePart, err := writer.CreateFormFile("files[]", path)
		if err != nil {
			return fmt.Errorf("error creating form field for file %s: %v", path, err)
		}

		// Open the file to be added
		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("error opening file %s: %v", path, err)
		}
		defer file.Close()

		// Copy the file's content into the form field
		n, err := io.Copy(filePart, file)
		if err != nil {
			return fmt.Errorf("error copying content for file %s: %v", path, err)
		}
		slog.Info(fmt.Sprintf("Copied %d bytes for file: %s", n, path))

		return nil
	})

	return err
}

// Function to send a POST request with files from a directory
func sendFiles(dirPath string, url string) error {
	// Create a new buffer to hold the multipart form data
	var requestBody bytes.Buffer
	// Create a new multipart writer
	writer := multipart.NewWriter(&requestBody)

	// Add files from the specified directory to the multipart form data
	err := addFilesToMultipart(writer, dirPath)
	if err != nil {
		return fmt.Errorf("error adding files to multipart: %v", err)
	}
	// Close the multipart writer to finalize the form data
	err = writer.Close()
	if err != nil {
		return fmt.Errorf("error closing multipart writer: %v", err)
	}

	// Create the POST request
	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %v", err)
	}

	// Set the Content-Type header to the appropriate multipart boundary
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Send the POST request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	fmt.Println("Response Headers:")
	for key, value := range resp.Header {
		fmt.Printf("%s: %s\n", key, value)
	}

	// Print the response status
	fmt.Println("Response status:", resp.Status)
	if resp.Header.Get("Content-Type") == "application/zip" {
		// Open output file for writing (not reading!)
		outputFile, err := os.OpenFile("data/output.zip", os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("error creating output file: %v", err)
		}
		defer outputFile.Close()

		// Copy the response body (ZIP data) to the output file
		n, err := io.Copy(outputFile, resp.Body)
		log.Printf("%d bytes written to file", n)
		if err != nil {
			return fmt.Errorf("error saving ZIP file: %v", err)
		}

		fmt.Println("ZIP file saved as data/output.zip")
	} else {
		io.Copy(os.Stdout, resp.Body)
	}

	return nil
}

func main() {
	// Path to the directory containing files to upload
	dirPath := "./zipdir/data" // Example directory
	url := "http://localhost:8080/api/archive/information"
	if len(os.Args) >= 2 {
		dirPath = os.Args[1]
	}
	if len(os.Args) >= 3 {
		url = os.Args[2]
	}
	// Send files from the directory as a multipart/form-data request
	err := sendFiles(dirPath, url)
	if err != nil {
		log.Fatalf("Error sending files: %v", err)
	}
}
