package main

import (
	"fmt"
	"net/http"
	"time"

	"greenlight.jennaluz.dev/internal/data"
)

func (app *application) handleMoviesPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Created a new movie")
}

func (app *application) handleMovieGet(w http.ResponseWriter, r *http.Request) {
	id, err := app.GetIDParam(r)
	if err != nil {
		app.writeNotFoundResponse(w, r)
		return
	}

	movie := data.Movie{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Casablanca",
		Runtime:   102,
		Genres:    []string{"drama", "romance", "war"},
		Version:   1,
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.writeServerErrorResponse(w, r, err)
	}
}
