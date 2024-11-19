package controllers

import (
	"archive/zip"
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"regexp"
	"zip-api/internal/infrastructure/config"
	"zip-api/internal/services"
	"zip-api/internal/services/zipservice"
)

// Filename regexp for Linux filesystem rules
var filenameRegexp = regexp.MustCompile(`^[^/\0][^/\0]*[^/\0 ]$`)

func ArchiveInfo(w http.ResponseWriter, r *http.Request) {
	slog.Info(fmt.Sprintf("Got the %s request with URL: %s", r.Method, r.URL.Path))
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Header().Set("Allow", "POST")
		return
	} else {
		// Sets boundary for request body
		r.Body = http.MaxBytesReader(w, r.Body, config.BodyLimitInBytes)

		// reader for streaming read
		reader, err := r.MultipartReader()
		if err != nil {
			message, code := "Cannot read multipart form data, incorrect format of Request", http.StatusBadRequest
			switch err {
			case multipart.ErrMessageTooLarge:
				message = fmt.Sprintf("Request body too large it must not exceed: %d bytes", config.BodyLimitInBytes)
			}
			jsonErrorRespond(w, message, code)
			return
		}

		// Read archive part
		archivePart, err := reader.NextPart()

		if err != nil {
			slog.Error(fmt.Sprintf("Error while reading the Request body part: %s", err))
			jsonErrorRespond(w, "Error while reading the request body", http.StatusInternalServerError)
			return
		} else if archivePart.FormName() != "file" {
			jsonErrorRespond(w, "Missing archive file", http.StatusBadRequest)
			return
		} else if !filenameRegexp.Match([]byte(archivePart.FileName())) {
			jsonErrorRespond(w, "Archive name does not comply the linux filename rules", http.StatusBadRequest)
			return
		}

		buf := bufio.NewReader(archivePart)
		sniff, _ := buf.Peek(512)
		contentType := http.DetectContentType(sniff)
		if contentType != "application/zip" {
			jsonErrorRespond(w, "Incorrect filetype zip file must be provided", http.StatusBadRequest)
			return
		}

		archiveMetadata, err := services.ZipServiceInstance.ZipInfo(buf, archivePart.FileName())
		if err != nil {
			statusCode, _ := http.StatusBadRequest, ""
			switch err {

			case zipservice.ErrIncorrectMimeType:
				jsonErrorRespond(w, "Archive contains not allowed mime type", statusCode)
				return
			default:
				slog.Error(fmt.Sprintf("Error while reading the zip archive: %s", err))
				jsonErrorRespond(w, "Error while unzipping the archive, incorrect zip format", statusCode)
				return
			}
		}
		jsonPayload, err := json.MarshalIndent(archiveMetadata, "", "   ")
		if err != nil {
			slog.Error(fmt.Sprintf("Error while marshalling JSON of archive struct: %s", err))
			jsonErrorRespond(w, "Error while returning the archive", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonPayload)
		return
	}
}

func ArchiveFiles(w http.ResponseWriter, r *http.Request) {
	slog.Info(fmt.Sprintf("Got the %s request with URL: %s", r.Method, r.URL.Path))
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Header().Set("Allow", "POST")
		return
	} else {
		// Sets boundary for request body
		r.Body = http.MaxBytesReader(w, r.Body, config.BodyLimitInBytes)
		defer r.Body.Close()

		// Reader for validation
		reader, err := r.MultipartReader()
		if err != nil {
			jsonErrorRespond(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Assuming the zipWriter is initialized to write to the response body
		zipWriter := zip.NewWriter(w)

		// Initialize the buffer for reading file parts
		buf := make([]byte, 1024) // Buffer to read file parts in chunks

		// Process each file part
		for filePart, err := reader.NextPart(); err != io.EOF; filePart, err = reader.NextPart() {
			defer filePart.Close()

			if err != nil {
				slog.Error(fmt.Sprintf("Error while reading the Request body part: %s", err.Error()))
				jsonErrorRespond(w, "Error while reading the request body", http.StatusInternalServerError)
				return
			}
			// Ensure this matches your filename rule and MIME type checks
			if filePart.FormName() != "files[]" {
				jsonErrorRespond(w, "Missing files[] form name", http.StatusBadRequest)
				return
			}

			// Create a new file entry in the ZIP archive
			fileHeader := &zip.FileHeader{
				Name:   filepath.Base(filePart.FileName()),
				Method: zip.Store, // Optional compression method (Deflate is commonly used)
			}

			// Create an entry in the ZIP archive for the current file part
			fileWriter, err := zipWriter.CreateHeader(fileHeader)
			if err != nil {
				jsonErrorRespond(w, err.Error(), http.StatusInternalServerError)
				return
			}

			total := 0
			for {
				n, err := filePart.Read(buf)
				total += n
				if err != nil && err != io.EOF {
					jsonErrorRespond(w, err.Error(), http.StatusInternalServerError)
					return
				} else if n == 0 {
					slog.Info(fmt.Sprintf("total %d bytes written into temporary archive with file: %s", total, filePart.FileName()))
					break
				}

				// Write the data chunk to the ZIP archive
				_, err = fileWriter.Write(buf[:n])
				if err != nil {
					jsonErrorRespond(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}

		}
		zipWriter.Close()

		// Content-Type header for ZIP file
		w.Header().Set("Content-Type", "application/zip")
	}

}

func MailFile(w http.ResponseWriter, r *http.Request) {
	slog.Info(fmt.Sprintf("Got the %s request with URL: %s", r.Method, r.URL.Path))
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Header().Set("Allow", "POST")
		return
	} else {
		return
	}

}
