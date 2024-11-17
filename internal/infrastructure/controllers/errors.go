package controllers

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type errorEnveloper struct {
	Message string `json:"error"`
}

func jsonErrorRespond(w http.ResponseWriter, message string, statusCode int) {
	errJSON := errorEnveloper{message}
	if statusCode == 0 {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(statusCode)
	}
	jsonError, err := json.MarshalIndent(errJSON, "", "   ")
	if err != nil {
		slog.Error(err.Error())
	}
	w.Write(append(jsonError, '\n'))
}
