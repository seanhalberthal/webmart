package main

import (
	"encoding/json"
	"net/http"
)

func writeJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1_048_576 // 1MB response limit
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))
	return json.NewDecoder(r.Body).Decode(data)
}

func respondWithErrorJSON(w http.ResponseWriter, status int, message string) error {
	type envelope struct {
		Error string `json:"error"`
	}

	return writeJSON(w, status, &envelope{Error: message})
}

func handleError(w http.ResponseWriter, statusCode int, err error) {
	err = respondWithErrorJSON(w, statusCode, err.Error())
	if err != nil {
		return
	}
}

func writeJSONResponse(w http.ResponseWriter, status int, data interface{}) error {
	type envelope struct {
		Data any `json:"data"`
	}

	return writeJSON(w, status, &envelope{Data: data})
}
