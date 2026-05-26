package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/we-work-in-the-cloud/nullcloud/backend/internal/store"
)

func NewServer(s store.Store) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(authMiddleware)

	r.Route("/v1", func(r chi.Router) {
		r.Route("/vpcs", vpcRoutes(s))
		r.Route("/subnets", subnetRoutes(s))
		r.Route("/instances", vsiRoutes(s))
	})

	return r
}
