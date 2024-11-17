package controllers

import (
	"bufio"
	"fmt"
	"log/slog"
	"mime/multipart"
	"net/http"
	"regexp"
	"zip-api/internal/infrastructure/config"
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
		w.Header().Set("Content-Type", "application/json")
		// Sets boundary for request body
		r.Body = http.MaxBytesReader(w, r.Body, config.BodyLimitInBytes)

		// reader for streaming read
		reader, err := r.MultipartReader()
		if err != nil {
			message, code := "", http.StatusBadRequest
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
			jsonErrorRespond(w, "Incorrect Content-Type, it must be application/zip", http.StatusBadRequest)
			return
		}
		


		return
	}
}

func ArchiveFiles(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Header().Set("Allow", "POST")
		return
	} else {
		w.Header().Set("Content-Type", "application/zip")
		return
	}

}

func MailFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Header().Set("Allow", "POST")
		return
	} else {
		return
	}

}
