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

func clusterRoutes(s store.Store) func(chi.Router) {
	return func(r chi.Router) {
		r.Get("/", listClusters(s))
		r.Post("/", createCluster(s))
		r.Get("/{id}", getCluster(s))
		r.Patch("/{id}", updateCluster(s))
		r.Delete("/{id}", deleteCluster(s))
	}
}

func listClusters(s store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := tokenFromCtx(r.Context())
		clusters, err := s.ListKubernetesClusters(r.Context(), token)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal_error", err.Error())
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"clusters": clusters})
	}
}

func createCluster(s store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := tokenFromCtx(r.Context())
		var req struct {
			Name      string `json:"name"`
			Version   string `json:"version"`
			NodeCount int    `json:"node_count"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
			writeError(w, http.StatusBadRequest, "bad_request", "name is required")
			return
		}
		if req.Version == "" {
			writeError(w, http.StatusBadRequest, "bad_request", "version is required")
			return
		}
		if req.NodeCount < 1 {
			writeError(w, http.StatusBadRequest, "bad_request", "node_count must be at least 1")
			return
		}
		id := uid.New("k8s")
		c := model.KubernetesCluster{
			ID:        id,
			Name:      req.Name,
			Status:    "running",
			CRN:       fmt.Sprintf("crn:nullcloud:cluster:%s", id),
			Version:   req.Version,
			NodeCount: req.NodeCount,
			CreatedAt: time.Now().UTC(),
		}
		if err := s.CreateKubernetesCluster(r.Context(), token, c); err != nil {
			writeError(w, http.StatusInternalServerError, "internal_error", err.Error())
			return
		}
		writeJSON(w, http.StatusCreated, c)
	}
}

func getCluster(s store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := tokenFromCtx(r.Context())
		id := chi.URLParam(r, "id")
		c, ok, err := s.GetKubernetesCluster(r.Context(), token, id)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal_error", err.Error())
			return
		}
		if !ok {
			writeError(w, http.StatusNotFound, "not_found", "Kubernetes cluster not found")
			return
		}
		writeJSON(w, http.StatusOK, c)
	}
}

func updateCluster(s store.Store) http.HandlerFunc {
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
		c, ok, err := s.GetKubernetesCluster(r.Context(), token, id)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal_error", err.Error())
			return
		}
		if !ok {
			writeError(w, http.StatusNotFound, "not_found", "Kubernetes cluster not found")
			return
		}
		if err := s.RenameKubernetesCluster(r.Context(), token, id, req.Name); err != nil {
			writeError(w, http.StatusInternalServerError, "internal_error", err.Error())
			return
		}
		c.Name = req.Name
		writeJSON(w, http.StatusOK, c)
	}
}

func deleteCluster(s store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := tokenFromCtx(r.Context())
		id := chi.URLParam(r, "id")
		_, ok, err := s.GetKubernetesCluster(r.Context(), token, id)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal_error", err.Error())
			return
		}
		if !ok {
			writeError(w, http.StatusNotFound, "not_found", "Kubernetes cluster not found")
			return
		}
		if err := s.DeleteKubernetesCluster(r.Context(), token, id); err != nil {
			writeError(w, http.StatusInternalServerError, "internal_error", err.Error())
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
