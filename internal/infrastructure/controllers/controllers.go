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
	"os"
	"path/filepath"
	"regexp"
	"zip-api/internal/infrastructure/config"
	"zip-api/internal/services"
	"zip-api/internal/services/zipservice"
	"zip-api/internal/utils"
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

		// Reader for validation
		reader, err := r.MultipartReader()
		if err != nil {
			jsonErrorRespond(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpArchive, err := os.CreateTemp("", "*.zip")
		if err != nil {
			jsonErrorRespond(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Create a zip.Writer that writes to the temporary ZIP file
		zipWriter := zip.NewWriter(tmpArchive)

		// Buffer for stream reading
		buf := make([]byte, 4096)

		// File parts validation and reading
		for filePart, err := reader.NextPart(); err != io.EOF; filePart, err = reader.NextPart() {
			if err != nil {
				slog.Error(fmt.Sprintf("Error while reading the Request body part: %s", err.Error()))
				jsonErrorRespond(w, "Error while reading the request body", http.StatusInternalServerError)
				return
			} else if filePart.FormName() != "files[]" {
				jsonErrorRespond(w, "Missing files[] form name", http.StatusBadRequest)
				return
			} else if !filenameRegexp.Match([]byte(filePart.FileName())) {
				jsonErrorRespond(w, fmt.Sprintf("File name: %s does not comply the linux filename rules", filePart.FileName()), http.StatusBadRequest)
				return
			}

			fileReader := bufio.NewReader(filePart)
			sniff, err := fileReader.Peek(512)
			if err != nil {
				slog.Error(fmt.Sprintf("error while Sniffing the file: %s", err.Error()))
				jsonErrorRespond(w, "Error while reading the file", http.StatusInternalServerError)
				return
			}

			contentType := http.DetectContentType(sniff)
			if contentType == "text/xml; charset=utf-8" {
				contentType = "application/xml"
			} else if contentType == "application/zip" {
				if utils.ComplySignature(sniff, utils.DocxSequence) {
					contentType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
				}
			}
			if !utils.In(contentType, utils.AllowedMimeTypes) {
				slog.Error(fmt.Sprintf("entered %s file with not allowed mimetype: %s", filePart.FileName(), contentType))
				jsonErrorRespond(w, fmt.Sprintf("entered %s file with not allowed mimetype: %s", filePart.FileName(), contentType), http.StatusBadRequest)
				return
			}

			// Archive writing
			fileHeader := &zip.FileHeader{
				Name:   filepath.Base(filePart.FileName()),
				Method: zip.Deflate,
			}
			// Create a new entry in the ZIP file using the file header
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
				} else if n == 0 && err == io.EOF {
					slog.Info(fmt.Sprintf("total %d bytes written into temporary archive with file: %s", total, filePart.FileName()))
					break
				}
				// Write the data to the destination file
				_, err = fileWriter.Write(buf[:n])
				if err != nil {
					jsonErrorRespond(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
		}

		zipWriter.Close()
		_, err = tmpArchive.Seek(0, io.SeekStart)

		if err != nil {
			slog.Error(fmt.Sprintf("Error while archiving the files: %s", err))
			jsonErrorRespond(w, "Error while archiving the files", http.StatusInternalServerError)
			return
		}
		defer tmpArchive.Close()

		buf = make([]byte, 4096)
		w.Header().Set("Content-Type", "application/zip")

		total := 0
		for {
			n, err := tmpArchive.Read(buf)
			if err != nil && err != io.EOF {
				slog.Error(fmt.Sprintf("Error while reading the archive reader: %s", err))
				jsonErrorRespond(w, "Error while archiving the files", http.StatusInternalServerError)
				return
			} else if n == 0 {
				slog.Info(fmt.Sprintf("Archived %d bytes", total))
				break
			}
			// Write the data to the destination file
			_, err = w.Write(buf[:n])
			total += n
			if err != nil {
				slog.Error(fmt.Sprintf("Error while reading the archive reader: %s", err))
				jsonErrorRespond(w, "Error while archiving the files", http.StatusInternalServerError)
				return
			}
		}
		return
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
