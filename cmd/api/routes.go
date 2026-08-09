package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.writeNotFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.writeMethodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.handleHealthCheckGet)
	router.HandlerFunc(http.MethodPost, "/v1/movies", app.handleMoviesPost)
	router.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.handleMovieGet)

	return app.recoverPanic(router)
}
