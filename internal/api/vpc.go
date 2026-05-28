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

func vpcRoutes(s store.Store) func(chi.Router) {
	return func(r chi.Router) {
		r.Get("/", listVPCs(s))
		r.Post("/", createVPC(s))
		r.Get("/{id}", getVPC(s))
		r.Patch("/{id}", updateVPC(s))
		r.Delete("/{id}", deleteVPC(s))
	}
}

func listVPCs(s store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := tokenFromCtx(r.Context())
		vpcs, err := s.ListVPCs(r.Context(), token)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal_error", err.Error())
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"vpcs": vpcs})
	}
}

func createVPC(s store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := tokenFromCtx(r.Context())
		var req struct {
			Name string `json:"name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
			writeError(w, http.StatusBadRequest, "bad_request", "name is required")
			return
		}
		id := uid.New("vpc")
		vpc := model.VPC{
			ID:        id,
			Name:      req.Name,
			Status:    "available",
			CRN:       fmt.Sprintf("crn:nullcloud:vpc:%s", id),
			CreatedAt: time.Now().UTC(),
		}
		if err := s.CreateVPC(r.Context(), token, vpc); err != nil {
			writeError(w, http.StatusInternalServerError, "internal_error", err.Error())
			return
		}
		writeJSON(w, http.StatusCreated, vpc)
	}
}

func getVPC(s store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := tokenFromCtx(r.Context())
		id := chi.URLParam(r, "id")
		vpc, ok, err := s.GetVPC(r.Context(), token, id)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal_error", err.Error())
			return
		}
		if !ok {
			writeError(w, http.StatusNotFound, "not_found", "VPC not found")
			return
		}
		writeJSON(w, http.StatusOK, vpc)
	}
}

func updateVPC(s store.Store) http.HandlerFunc {
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
		vpc, ok, err := s.GetVPC(r.Context(), token, id)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal_error", err.Error())
			return
		}
		if !ok {
			writeError(w, http.StatusNotFound, "not_found", "VPC not found")
			return
		}
		if err := s.RenameVPC(r.Context(), token, id, req.Name); err != nil {
			writeError(w, http.StatusInternalServerError, "internal_error", err.Error())
			return
		}
		vpc.Name = req.Name
		writeJSON(w, http.StatusOK, vpc)
	}
}

func deleteVPC(s store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := tokenFromCtx(r.Context())
		id := chi.URLParam(r, "id")
		_, ok, err := s.GetVPC(r.Context(), token, id)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal_error", err.Error())
			return
		}
		if !ok {
			writeError(w, http.StatusNotFound, "not_found", "VPC not found")
			return
		}
		if err := s.DeleteVPC(r.Context(), token, id); err != nil {
			writeError(w, http.StatusInternalServerError, "internal_error", err.Error())
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
