package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
	"greenlight.jennaluz.dev/internal/validator"
)

type envelope map[string]any

func (app *application) GetIDParam(r *http.Request) (int, error) {
	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil || id < 1 {
		return 0, errors.New("Invalid id parameter")
	}

	return id, nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	// Limit the size of the request body to 1MiB.
	r.Body = http.MaxBytesReader(w, r.Body, 1_048_576)

	// Init the json.Decoder and call DisallowUnknownFields() on it before
	// decoding.
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var maxBytesError *http.MaxBytesError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("Body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return fmt.Errorf("Body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("Body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}

			return fmt.Errorf("Body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		case errors.Is(err, io.EOF):
			return errors.New("Body must not be empty")

		case strings.HasPrefix(err.Error(), "json: unknown field"):
			// TODO: We shouldn't return user input.
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		case errors.As(err, &maxBytesError):
			// TODO: Return 413 Content Too Large
			return fmt.Errorf("body must not be larger than %d bytes", maxBytesError.Limit)

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}

	// Call decode again to check if the request body included more than one
	// JSON object. If only one JSON object was sent, then an io.EOF error will
	// be returned. If anything else is returned, we know that additionally data
	// was included.
	err = decoder.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		return errors.New("body must only contain a single JSON value")
	}

	return nil
}

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	js = append(js, '\n')

	// Add all headers (map[string][]string)
	for key, values := range headers {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

// Returns a string value from the query string or the provided default value if
// the no matching key could be found.
func (app *application) readString(qs url.Values, key string, defaultValue string) string {
	// Extract the value for a given key from the query string.
	// If no key exists, this will return an emtpy string.
	str := qs.Get(key)
	if str == "" {
		return defaultValue
	}

	return str
}

// Reads a string value from the query string and then splits it into a slice on
// the comma character.
// Returns the provided default value if no matching key could be found.
func (app *application) readCSV(qs url.Values, key string, defaultValues []string) []string {
	csv := qs.Get(key)
	if csv == "" {
		return defaultValues
	}

	return strings.Split(csv, ",")
}

// Reads a string value from the query string and converts it to an integer
// before returning it.
// Returns the provided default value if no matching key could be found.
// Returns an error message in the provided Validator instance if it cannot be
// converted to a valid integer.
func (app *application) readInt(qs url.Values, key string, defaultValue int, v *validator.Validator) int {
	str := qs.Get(key)
	if str == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(str)
	if err != nil {
		v.AddError(key, "must be an integer value")
		return defaultValue
	}

	return i
}
