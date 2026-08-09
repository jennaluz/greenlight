package main

import (
	"net/http"
)

func (app *application) handleHealthCheckGet(w http.ResponseWriter, r *http.Request) {
	env := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": app.config.env,
			"version":     version,
		},
	}

	err := app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		app.writeServerErrorResponse(w, r, err)
		return
	}
}
