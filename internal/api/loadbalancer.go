package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/we-work-in-the-cloud/nullcloud/backend/internal/model"
	"github.com/we-work-in-the-cloud/nullcloud/backend/internal/store"
	"github.com/we-work-in-the-cloud/nullcloud/backend/internal/uid"
)

func loadBalancerRoutes(s store.Store) func(chi.Router) {
	return func(r chi.Router) {
		r.Get("/", listLoadBalancers(s))
		r.Post("/", createLoadBalancer(s))
		r.Get("/{id}", getLoadBalancer(s))
		r.Patch("/{id}", updateLoadBalancer(s))
		r.Delete("/{id}", deleteLoadBalancer(s))
	}
}

func listLoadBalancers(s store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := tokenFromCtx(r.Context())
		lbs, err := s.ListLoadBalancers(r.Context(), token)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal_error", err.Error())
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"load_balancers": lbs})
	}
}

func createLoadBalancer(s store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := tokenFromCtx(r.Context())
		var req struct {
			Name     string `json:"name"`
			Protocol string `json:"protocol"`
			Port     int    `json:"port"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
			writeError(w, http.StatusBadRequest, "bad_request", "name is required")
			return
		}
		if req.Protocol != "tcp" && req.Protocol != "http" && req.Protocol != "https" {
			writeError(w, http.StatusBadRequest, "bad_request", "protocol must be tcp, http, or https")
			return
		}
		if req.Port < 1 || req.Port > 65535 {
			writeError(w, http.StatusBadRequest, "bad_request", "port must be between 1 and 65535")
			return
		}
		id := uid.New("lb")
		lb := model.LoadBalancer{
			ID:        id,
			Name:      req.Name,
			Status:    "active",
			CRN:       fmt.Sprintf("crn:nullcloud:loadbalancer:%s", id),
			Protocol:  req.Protocol,
			Port:      req.Port,
			CreatedAt: time.Now().UTC(),
		}
		if err := s.CreateLoadBalancer(r.Context(), token, lb); err != nil {
			writeError(w, http.StatusInternalServerError, "internal_error", err.Error())
			return
		}
		writeJSON(w, http.StatusCreated, lb)
	}
}

func getLoadBalancer(s store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := tokenFromCtx(r.Context())
		id := chi.URLParam(r, "id")
		lb, ok, err := s.GetLoadBalancer(r.Context(), token, id)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal_error", err.Error())
			return
		}
		if !ok {
			writeError(w, http.StatusNotFound, "not_found", "Load balancer not found")
			return
		}
		writeJSON(w, http.StatusOK, lb)
	}
}

func updateLoadBalancer(s store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := tokenFromCtx(r.Context())
		id := chi.URLParam(r, "id")
		var req struct {
			Name string `json:"name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
			writeError(w, http.StatusBadRequest, "bad_request", "name is required")
			return
		}
		lb, ok, err := s.GetLoadBalancer(r.Context(), token, id)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal_error", err.Error())
			return
		}
		if !ok {
			writeError(w, http.StatusNotFound, "not_found", "Load balancer not found")
			return
		}
		if err := s.RenameLoadBalancer(r.Context(), token, id, req.Name); err != nil {
			writeError(w, http.StatusInternalServerError, "internal_error", err.Error())
			return
		}
		lb.Name = req.Name
		writeJSON(w, http.StatusOK, lb)
	}
}

func deleteLoadBalancer(s store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := tokenFromCtx(r.Context())
		id := chi.URLParam(r, "id")
		_, ok, err := s.GetLoadBalancer(r.Context(), token, id)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal_error", err.Error())
			return
		}
		if !ok {
			writeError(w, http.StatusNotFound, "not_found", "Load balancer not found")
			return
		}
		if err := s.DeleteLoadBalancer(r.Context(), token, id); err != nil {
			writeError(w, http.StatusInternalServerError, "internal_error", err.Error())
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
