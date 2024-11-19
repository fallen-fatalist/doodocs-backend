package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	zipPath := "data/dummy.zip"
	url := "http://localhost:8080/api/archive/information"
	if len(os.Args) >= 2 {
		zipPath = os.Args[1]
	}
	if len(os.Args) >= 3 {
		url = os.Args[2]
	}

	getInfoZip(zipPath, url)
}

func getInfoZip(fileName, url string) {
	zipFile, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer zipFile.Close()

	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)
	part, err := writer.CreateFormFile("file", filepath.Base(zipFile.Name()))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating form file: %v\n", err)
		return
	}

	size, err := io.Copy(part, zipFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error copying zip file data: %v\n", err)
		return
	}
	writer.Close()
	fmt.Fprintf(os.Stdout, "Copied %v bytes for uploading...\n", size)

	response, err := http.Post(url, writer.FormDataContentType(), &buffer)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error making POST request to get zip info: %v\n", err)
		return
	}
	defer response.Body.Close()

	fmt.Fprintf(os.Stdout, "Successfully got the response : %v\n", response.StatusCode)

	// Print the Response Headers
	fmt.Println("Response Headers:")
	for key, value := range response.Header {
		fmt.Printf("%s: %s\n", key, value)
	}

	// // Read the Response Body
	// body, err := io.ReadAll(response.Body)
	// if err != nil {
	// 	log.Fatalf("Error reading response body: %s", err)
	// }

	// Print the Response Body
	fmt.Println("\nResponse Body:")
	// fmt.Println(string(body))
	io.Copy(os.Stdout, response.Body)

	// Optionally, print the status code
	fmt.Printf("\nResponse Status: %s\n", response.Status)
}
