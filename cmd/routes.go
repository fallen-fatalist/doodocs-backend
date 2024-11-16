package cmd

import (
	"net/http"
	"zip-api/internal/infrastructure/controllers"
)

func routes() *http.ServeMux {
	mux := http.NewServeMux()

	// archive routes
	mux.HandleFunc("api/archive/information", controllers.ArchiveInfo)
	mux.HandleFunc("api/archive/files", controllers.ArchiveFiles)

	// mail routes
	mux.HandleFunc("api/mail/file", controllers.ArchiveInfo)
	return mux
}
