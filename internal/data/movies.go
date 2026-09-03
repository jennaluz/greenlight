package data

import (
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
	"greenlight.jennaluz.dev/internal/validator"
)

type Movie struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Year      int       `json:"year,omitzero"`
	Runtime   Runtime   `json:"runtime,omitzero"`
	Genres    []string  `json:"genres,omitzero"`
	Version   int       `json:"version"`
}

type MovieModel struct {
	DB *sql.DB
}

func ValidateMovie(v *validator.Validator, movie Movie) {
	v.Check(movie.Title != "", "title", "must be provided")
	v.Check(len(movie.Title) <= 500, "title", "must not be more than 500 bytes long")

	v.Check(movie.Year != 0, "year", "must be provided")
	v.Check(movie.Year >= 1888, "year", "must be greater than 1888")
	v.Check(movie.Year <= time.Now().Year(), "year", "must not be in the future")

	v.Check(movie.Runtime != 0, "runtime", "must be provided")
	v.Check(movie.Runtime > 0, "runtime", "must be a positive integer")

	v.Check(movie.Genres != nil, "genres", "must be provided")
	v.Check(len(movie.Genres) != 0, "genres", "must contain at least 1 genre")
	v.Check(len(movie.Genres) <= 5, "genres", "must not contain more than 5 genres")
	v.Check(validator.IsUnique(movie.Genres), "genres", "must not contain duplicate values")
}

func (m MovieModel) Add(movie Movie) (Movie, error) {
	query := `
    	INSERT INTO movies(
    	    title,
    	    year,
    	    runtime,
    	    genres)
    	VALUES(
    	    $1,
    	    $2,
    	    $3,
    	    $4)
    	RETURNING
    	    id,
    	    created_at,
    	    version;
	`

	args := []any{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres)}

	err := m.DB.QueryRow(query, args...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)

	return movie, err
}

func (m MovieModel) Get(id int) (Movie, error) {
	if id < 1 {
		return Movie{}, ErrRecordNotFound
	}

	query := `
    	SELECT
			m.id,
			m.created_at,
			m.title,
			m.year,
			m.runtime,
			m.genres,
			m.version
    	FROM movies m
    	WHERE m.id = $1;
	`

	var movie Movie

	err := m.DB.QueryRow(query, id).Scan(
		&movie.ID,
		&movie.CreatedAt,
		&movie.Title,
		&movie.Year,
		&movie.Runtime,
		pq.Array(&movie.Genres),
		&movie.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return Movie{}, ErrRecordNotFound
		default:
			return Movie{}, err
		}
	}

	return movie, nil
}

func (m MovieModel) Update(movie Movie) (Movie, error) {
	query := `
		UPDATE movies
		SET title = $1,
		    year  = $2,
			runtime = $3,
			genres  = $4,
			version = version + 1
		WHERE id = $5
		RETURNING version;
	`

	// Create an args slice containing values for placeholder parameters.
	args := []any{
		movie.Title,
		movie.Year,
		movie.Runtime,
		pq.Array(movie.Genres),
		movie.ID,
	}

	err := m.DB.QueryRow(query, args...).Scan(&movie.Version)
	return movie, err
}

func (m MovieModel) Delete(id int) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM movies
		WHERE id = $1;
	`

	result, err := m.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
