package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/we-work-in-the-cloud/nullcloud/backend/internal/store"
)

// NewServer serves both the UI at /ui and the API at /v1 on a single handler.
func NewServer(s store.Store, allowedTokens []string) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Mount("/ui", NewUIHandler())
	r.Mount("/", NewAPIHandler(s, allowedTokens))
	return r
}

// NewAPIHandler serves only the /v1 API routes with token authentication.
func NewAPIHandler(s store.Store, allowedTokens []string) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Group(func(r chi.Router) {
		r.Use(tokenMiddleware(allowedTokens))
		r.Route("/v1", func(r chi.Router) {
			r.Route("/vpcs", vpcRoutes(s))
			r.Route("/subnets", subnetRoutes(s))
			r.Route("/instances", vsiRoutes(s))
		})
	})
	return r
}

// NewUIHandler serves the UI at / plus its static assets (style.css, app.js).
func NewUIHandler() http.Handler {
	r := chi.NewRouter()
	r.Get("/", uiHandler())
	r.Get("/style.css", uiCSSHandler())
	r.Get("/app.js", uiJSHandler())
	return r
}
