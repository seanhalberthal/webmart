package main

import (
	"net/http"
)

// HealthCheck godoc
//
//	@Summary		Health check
//	@Description	Returns API status, environment, and version
//	@Tags			ops
//	@Produce		json
//	@Success		200	{object}	map[string]string
//	@Failure		500	{object}	error
//	@Router			/healthcheck [get]
func (app *application) healthCheckHandler(w http.ResponseWriter, _ *http.Request) {
	data := map[string]string{
		"status":  "ok",
		"env":     app.config.env,
		"version": version,
	}

	if err := writeJSONResponse(w, http.StatusOK, data); err != nil {
		handleError(w, http.StatusInternalServerError, err)
		return
	}
}
