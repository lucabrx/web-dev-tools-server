package main

import (
	"errors"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/wdt/internal/data"
	validator "github.com/wdt/internal/validators"
	"golang.org/x/time/rate"
)

func (app *application) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", app.config.ClientAddress)
			w.Header().Set("Access-Control-Allow-Credentials", "true")

			if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
				w.Header().Set("Access-Control-Allow-Methods", "PATCH, DELETE, GET, POST")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
				w.WriteHeader(http.StatusOK)
				return
			}
		}

		w.Header().Add("Vary", "Origin")
		w.Header().Add("Vary", "Access-Control-Request-Method")

		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuthenticatedUser(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)
		if user.IsAnonymous() {
			app.authenticationRequiredResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")
		cookie, err := r.Cookie("session")

		if errors.Is(err, http.ErrNoCookie) {
			r = app.contextSetUser(r, data.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		} else if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		token := cookie.Value
		v := validator.New()
		if data.ValidateTokenPlaintext(v, token); !v.Valid() {
			app.invalidAuthenticationTokenResponse(w, r)
			cookie := app.sessionCookie("", time.Unix(0, 0))
			http.SetCookie(w, cookie)
			return
		}

		user, err := app.models.Users.GetForToken(data.ScopeAuthentication, token)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				cookie := app.sessionCookie("", time.Unix(0, 0))
				http.SetCookie(w, cookie)
				app.invalidAuthenticationTokenResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}

		r = app.contextSetUser(r, user)
		next.ServeHTTP(w, r)
	})

}

func (app *application) adminPermission(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)
		if user.IsAnonymous() {
			app.authenticationRequiredResponse(w, r)
			return
		}

		if user.Role != "admin" {
			app.notPermittedResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	}
}

func (app *application) rateLimit(next http.Handler) http.Handler {

	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	//  goroutine which removes old entries from the clients map once every minute.
	go func() {
		for {
			time.Sleep(time.Minute)
			// Lock the mutex to prevent any rate limiter checks from happening while the cleanup is taking place.
			mu.Lock()

			// Loop through all clients. If they haven't been seen within the last three minutes, delete the corresponding entry from the map.
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}

			mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		mu.Lock()

		if _, found := clients[ip]; !found {
			// Create and add a new client struct to the map if it doesn't already exist.
			clients[ip] = &client{limiter: rate.NewLimiter(2, 4)}
		}

		// Update the last seen time for the client.
		clients[ip].lastSeen = time.Now()

		if !clients[ip].limiter.Allow() {
			mu.Unlock()
			app.rateLimitExceededResponse(w, r)
			return
		}

		mu.Unlock()

		next.ServeHTTP(w, r)
	})
}
