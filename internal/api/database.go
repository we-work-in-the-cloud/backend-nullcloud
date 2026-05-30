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

func databaseRoutes(s store.Store) func(chi.Router) {
	return func(r chi.Router) {
		r.Get("/", listDatabases(s))
		r.Post("/", createDatabase(s))
		r.Get("/{id}", getDatabase(s))
		r.Patch("/{id}", updateDatabase(s))
		r.Delete("/{id}", deleteDatabase(s))
	}
}

func listDatabases(s store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := tokenFromCtx(r.Context())
		dbs, err := s.ListDatabases(r.Context(), token)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal_error", err.Error())
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"databases": dbs})
	}
}

func createDatabase(s store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := tokenFromCtx(r.Context())
		var req struct {
			Name    string `json:"name"`
			Engine  string `json:"engine"`
			Version string `json:"version"`
			Plan    string `json:"plan"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" {
			writeError(w, http.StatusBadRequest, "bad_request", "name is required")
			return
		}
		if req.Engine != "postgres" && req.Engine != "mysql" && req.Engine != "mariadb" {
			writeError(w, http.StatusBadRequest, "bad_request", "engine must be postgres, mysql, or mariadb")
			return
		}
		if req.Version == "" {
			writeError(w, http.StatusBadRequest, "bad_request", "version is required")
			return
		}
		if req.Plan != "small" && req.Plan != "medium" && req.Plan != "large" {
			writeError(w, http.StatusBadRequest, "bad_request", "plan must be small, medium, or large")
			return
		}
		id := uid.New("db")
		db := model.Database{
			ID:        id,
			Name:      req.Name,
			Status:    "available",
			CRN:       fmt.Sprintf("crn:nullcloud:database:%s", id),
			Engine:    req.Engine,
			Version:   req.Version,
			Plan:      req.Plan,
			CreatedAt: time.Now().UTC(),
		}
		if err := s.CreateDatabase(r.Context(), token, db); err != nil {
			writeError(w, http.StatusInternalServerError, "internal_error", err.Error())
			return
		}
		writeJSON(w, http.StatusCreated, db)
	}
}

func getDatabase(s store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := tokenFromCtx(r.Context())
		id := chi.URLParam(r, "id")
		db, ok, err := s.GetDatabase(r.Context(), token, id)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal_error", err.Error())
			return
		}
		if !ok {
			writeError(w, http.StatusNotFound, "not_found", "Database not found")
			return
		}
		writeJSON(w, http.StatusOK, db)
	}
}

func updateDatabase(s store.Store) http.HandlerFunc {
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
		db, ok, err := s.GetDatabase(r.Context(), token, id)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal_error", err.Error())
			return
		}
		if !ok {
			writeError(w, http.StatusNotFound, "not_found", "Database not found")
			return
		}
		if err := s.RenameDatabase(r.Context(), token, id, req.Name); err != nil {
			writeError(w, http.StatusInternalServerError, "internal_error", err.Error())
			return
		}
		db.Name = req.Name
		writeJSON(w, http.StatusOK, db)
	}
}

func deleteDatabase(s store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := tokenFromCtx(r.Context())
		id := chi.URLParam(r, "id")
		_, ok, err := s.GetDatabase(r.Context(), token, id)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal_error", err.Error())
			return
		}
		if !ok {
			writeError(w, http.StatusNotFound, "not_found", "Database not found")
			return
		}
		if err := s.DeleteDatabase(r.Context(), token, id); err != nil {
			writeError(w, http.StatusInternalServerError, "internal_error", err.Error())
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
