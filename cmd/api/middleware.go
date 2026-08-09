package main

import (
	"fmt"
	"net/http"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			panicValue := recover()
			if panicValue != nil {
				// If there was a panice, set a "Connection: close" header on
				// the response. This acts as a trigger to make Go's HTTP server
				// automatically close the current connection after the response
				// has been sent.
				w.Header().Set("Connection", "close")

				app.writeServerErrorResponse(w, r, fmt.Errorf("%v", panicValue))
			}
		}()

		next.ServeHTTP(w, r)
	})
}
