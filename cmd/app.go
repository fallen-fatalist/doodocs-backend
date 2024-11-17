package cmd

import (
	"log"
	"log/slog"
	"net/http"
	"os"
	"zip-api/internal/infrastructure/config"
)

func Run() {
	err := config.Init()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	mux := routes()

	slog.Info("Running server on port: " + config.Port)
	err = http.ListenAndServe(":"+config.Port, mux)
	log.Fatal(err)
}
