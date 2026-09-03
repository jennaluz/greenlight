package main

import (
	"errors"
	"fmt"
	"net/http"

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

	movie, err = app.models.Movies.Add(movie)
	if err != nil {
		app.writeServerErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
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

	movie, err := app.models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.writeNotFoundResponse(w, r)
		default:
			app.writeServerErrorResponse(w, r, err)
		}

		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.writeServerErrorResponse(w, r, err)
	}
}

func (app *application) handleMoviePut(w http.ResponseWriter, r *http.Request) {
	id, err := app.GetIDParam(r)
	if err != nil {
		app.writeNotFoundResponse(w, r)
		return
	}

	movie, err := app.models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.writeNotFoundResponse(w, r)
		default:
			app.writeServerErrorResponse(w, r, err)
		}
	}

	var input struct {
		Title   string       `json:"title"`
		Year    int          `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.writeBadRequestResponse(w, r, err)
		return
	}

	movie.Title = input.Title
	movie.Year = input.Year
	movie.Runtime = input.Runtime
	movie.Genres = input.Genres

	v := validator.New()
	if data.ValidateMovie(v, movie); !v.IsValid() {
		app.writeUnprocessableEntityResponse(w, r, v.Errors)
	}

	movie, err = app.models.Movies.Update(movie)
	if err != nil {
		app.writeServerErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.writeServerErrorResponse(w, r, err)
	}
}

func (app *application) handleMovieDelete(w http.ResponseWriter, r *http.Request) {
	id, err := app.GetIDParam(r)
	if err != nil {
		app.writeNotFoundResponse(w, r)
		return
	}

	err = app.models.Movies.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.writeNotFoundResponse(w, r)
		default:
			app.writeServerErrorResponse(w, r, err)
		}
	}

	app.writeSuccessResponse(w, r, "movie successfully deleted")
}
