package controllers

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type errorEnveloper struct {
	message string `json: "error"`
}

func jsonErrorRespond(w http.ResponseWriter, err error, statusCode int) {
	errJSON := errorEnveloper{err.Error()}
	if statusCode == 0 {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(statusCode)
	}
	jsonError, err := json.MarshalIndent(errJSON, "", "   ")
	if err != nil {
		slog.Error(err.Error())
	}
	w.Write(jsonError)
}
