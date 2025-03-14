package main

import (
	"net/http"
)

func (app *application) healthCheckHandler(w http.ResponseWriter, _ *http.Request) {
	data := map[string]string{
		"status":  "ok",
		"env":     app.config.env,
		"version": version,
	}

	if err := writeJSON(w, http.StatusOK, data); err != nil {
		err := respondWithErrorJSON(w, http.StatusInternalServerError, err.Error())
		if err != nil {
			return
		}
	}
}
