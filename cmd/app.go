package cmd

import (
	"log"
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

	err = http.ListenAndServe(":"+config.Port, mux)
	log.Fatal(err)
}
