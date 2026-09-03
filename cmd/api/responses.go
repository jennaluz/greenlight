package main

import (
	"fmt"
	"net/http"
)

type DefaultResponse struct {
	Status int    `json:"status"`
	Title  string `json:"title"`
	Detail string `json:"detail"`
}

func (app *application) logError(r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	app.logger.Error(err.Error(), "method", method, "uri", uri)
}

func (app *application) writeErrorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := envelope{"error": message}

	app.logger.Info("Writing error response", "env", env)

	err := app.writeJSON(w, status, env, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (app *application) writeServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)

	message := DefaultResponse{
		Status: http.StatusInternalServerError,
		Title:  "Internal Server Error",
		Detail: "Unable to process your request",
	}

	app.writeErrorResponse(w, r, http.StatusInternalServerError, message)
}

func (app *application) writeBadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.writeErrorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (app *application) writeNotFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := DefaultResponse{
		Status: http.StatusNotFound,
		Title:  "Not Found",
		Detail: "Unable to find your resource",
	}

	app.writeErrorResponse(w, r, http.StatusNotFound, message)
}

func (app *application) writeMethodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := DefaultResponse{
		Status: http.StatusMethodNotAllowed,
		Title:  "Method Not Allowed",
		Detail: fmt.Sprintf("The %s method is not surrported for this resource", r.Method),
	}

	app.writeErrorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func (app *application) writeUnprocessableEntityResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.writeErrorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

func (app *application) writeSuccessResponse(w http.ResponseWriter, r *http.Request, detail string) {
	message := DefaultResponse{
		Status: http.StatusOK,
		Title:  "Success",
		Detail: detail,
	}

	err := app.writeJSON(w, http.StatusOK, envelope{"success": message}, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
