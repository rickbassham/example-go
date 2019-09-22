// Package router is used to define the routes to our request handlers.
package router

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth"
	newrelic "github.com/newrelic/go-agent"
	"go.uber.org/zap"

	"github.com/rickbassham/example-go/chiapi/middleware"
)

// Handler exposes the functions for handling web requests.
type Handler interface {
	Health(w http.ResponseWriter, r *http.Request)
	Cached(w http.ResponseWriter, r *http.Request)
	Protected(w http.ResponseWriter, r *http.Request)
	NotFound(w http.ResponseWriter, r *http.Request)
	Unauthorized(w http.ResponseWriter, r *http.Request)
}

// NewRouter creates a new CORS enabled router for our API. All requests will be logged and
// instrumented with New Relic.
func NewRouter(h Handler, log *zap.Logger, nr newrelic.Application, tokenAuth *jwtauth.JWTAuth, version, corsOrigin string) http.Handler {
	r := chi.NewRouter()

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{corsOrigin},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	r.Use(middleware.Version(version))
	r.Use(middleware.TraceID)
	r.Use(middleware.Logger(log))
	r.Use(middleware.NewRelicChiRouter(nr))
	r.Use(cors.Handler)
	r.NotFound(h.NotFound)

	r.Get("/health", h.Health)
	r.Get("/cached", h.Cached)

	r.Route("/protected", func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(middleware.Authenticator(http.HandlerFunc(h.Unauthorized)))
		r.Use(middleware.User)

		r.Get("/{id:[0-9]+}", h.Protected)
	})

	return r
}
