package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/tomasen/realip"
	"golang.org/x/time/rate"
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

func (app *application) rateLimit(next http.Handler) http.Handler {
	if !app.config.limiter.enabled {
		return next
	}

	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	// Declare a mutex and a map to hold client IP addresses and rate limiters.
	var (
		clientsMu sync.Mutex
		clients   = make(map[string]*client)
	)

	// Launch a background goroutine to remove old entries from the clients
	// map once every minute.
	go func() {
		for {
			time.Sleep(time.Minute)

			clientsMu.Lock()

			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}

			clientsMu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := realip.FromRequest(r)

		clientsMu.Lock()

		if _, found := clients[ip]; !found {
			// Init a new rate limiter which allows an average of two requests per
			// second and a maximum of 4 requests in a single burst.
			clients[ip] = &client{
				limiter: rate.NewLimiter(rate.Limit(app.config.limiter.rps), app.config.limiter.burst),
			}
		}

		clients[ip].lastSeen = time.Now()

		if !clients[ip].limiter.Allow() {
			clientsMu.Unlock()
			app.writeTooManyRequestsResponse(w, r)
			return
		}

		// Defer is not used to unlock the mutex, as that would mean it would
		// not be unlocked until all of the downstream handlers have been
		// returned.
		clientsMu.Unlock()

		next.ServeHTTP(w, r)
	})
}
