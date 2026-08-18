package main

import (
	"net/http"
	"time"

	"greenlight.jennaluz.dev/internal/data"
	"greenlight.jennaluz.dev/internal/validator"
)

func (app *application) handleMoviesPost(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   string       `json:"title"`
		Year    int          `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.writeBadRequestResponse(w, r, err)
		return
	}

	movie := data.Movie{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres,
	}

	v := validator.New()
	if data.ValidateMovie(v, movie); !v.IsValid() {
		app.writeUnprocessableEntityResponse(w, r, v.Errors)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": input}, nil)
	if err != nil {
		app.writeErrorResponse(w, r, http.StatusBadRequest, err.Error())
	}
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
