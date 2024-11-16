package cmd

import (
	"log"
	"net/http"
	"zip-api/internal/infrastructure/config"
)

func Run() {
	config.Init()

	mux := routes()

	err := http.ListenAndServe(":"+config.Port, mux)
	log.Fatal(err)
}
