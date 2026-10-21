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
	router.HandlerFunc(http.MethodGet, "/v1/movies", app.hanldeMoviesGet)
	router.HandlerFunc(http.MethodPost, "/v1/movies", app.handleMoviesPost)
	router.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.handleMovieGet)
	router.HandlerFunc(http.MethodPut, "/v1/movies/:id", app.handleMoviePut)
	router.HandlerFunc(http.MethodDelete, "/v1/movies/:id", app.handleMovieDelete)

	// Wraper the router with rateLimit() middleware.
	// The rate limit should come "after" the recover panic so that any panics
	// in the rate limit can be recovered. Otherwise, the rate limit should be
	// used as early as possible.
	return app.recoverPanic(app.rateLimit(router))
}
