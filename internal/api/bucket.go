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

func bucketRoutes(s store.Store) func(chi.Router) {
	return func(r chi.Router) {
		r.Get("/", listBuckets(s))
		r.Post("/", createBucket(s))
		r.Get("/{id}", getBucket(s))
		r.Patch("/{id}", updateBucket(s))
		r.Delete("/{id}", deleteBucket(s))
	}
}

func listBuckets(s store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := tokenFromCtx(r.Context())
		buckets, err := s.ListBuckets(r.Context(), token)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal_error", err.Error())
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"buckets": buckets})
	}
}

func createBucket(s store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := tokenFromCtx(r.Context())
		var req struct {
			Name   string `json:"name"`
			Region string `json:"region"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
			writeError(w, http.StatusBadRequest, "bad_request", "name is required")
			return
		}
		if req.Region == "" {
			req.Region = "us-east-1"
		}
		id := uid.New("bkt")
		b := model.Bucket{
			ID:        id,
			Name:      req.Name,
			Status:    "available",
			CRN:       fmt.Sprintf("crn:nullcloud:bucket:%s", id),
			Region:    req.Region,
			CreatedAt: time.Now().UTC(),
		}
		if err := s.CreateBucket(r.Context(), token, b); err != nil {
			writeError(w, http.StatusInternalServerError, "internal_error", err.Error())
			return
		}
		writeJSON(w, http.StatusCreated, b)
	}
}

func getBucket(s store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := tokenFromCtx(r.Context())
		id := chi.URLParam(r, "id")
		b, ok, err := s.GetBucket(r.Context(), token, id)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal_error", err.Error())
			return
		}
		if !ok {
			writeError(w, http.StatusNotFound, "not_found", "Bucket not found")
			return
		}
		writeJSON(w, http.StatusOK, b)
	}
}

func updateBucket(s store.Store) http.HandlerFunc {
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
		b, ok, err := s.GetBucket(r.Context(), token, id)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal_error", err.Error())
			return
		}
		if !ok {
			writeError(w, http.StatusNotFound, "not_found", "Bucket not found")
			return
		}
		if err := s.RenameBucket(r.Context(), token, id, req.Name); err != nil {
			writeError(w, http.StatusInternalServerError, "internal_error", err.Error())
			return
		}
		b.Name = req.Name
		writeJSON(w, http.StatusOK, b)
	}
}

func deleteBucket(s store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := tokenFromCtx(r.Context())
		id := chi.URLParam(r, "id")
		_, ok, err := s.GetBucket(r.Context(), token, id)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal_error", err.Error())
			return
		}
		if !ok {
			writeError(w, http.StatusNotFound, "not_found", "Bucket not found")
			return
		}
		if err := s.DeleteBucket(r.Context(), token, id); err != nil {
			writeError(w, http.StatusInternalServerError, "internal_error", err.Error())
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
