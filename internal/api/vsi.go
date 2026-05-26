package api

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/we-work-in-the-cloud/nullcloud/backend/internal/model"
	"github.com/we-work-in-the-cloud/nullcloud/backend/internal/store"
	"github.com/we-work-in-the-cloud/nullcloud/backend/internal/uid"
)

func vsiRoutes(s store.Store) func(chi.Router) {
	return func(r chi.Router) {
		r.Get("/", listVSIs(s))
		r.Post("/", createVSI(s))
		r.Get("/{id}", getVSI(s))
		r.Delete("/{id}", deleteVSI(s))
	}
}

func listVSIs(s store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := tokenFromCtx(r.Context())
		vsis, err := s.ListVSIs(r.Context(), token)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal_error", err.Error())
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"instances": vsis})
	}
}

func createVSI(s store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := tokenFromCtx(r.Context())
		var req struct {
			Name   string `json:"name"`
			Subnet struct {
				ID string `json:"id"`
			} `json:"subnet"`
			Profile struct {
				Name string `json:"name"`
			} `json:"profile"`
			Image struct {
				ID string `json:"id"`
			} `json:"image"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" || req.Subnet.ID == "" {
			writeError(w, http.StatusBadRequest, "bad_request", "name and subnet.id are required")
			return
		}
		sub, ok, err := s.GetSubnet(r.Context(), token, req.Subnet.ID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal_error", err.Error())
			return
		}
		if !ok {
			writeError(w, http.StatusNotFound, "not_found", "Subnet not found")
			return
		}
		id := uid.New("vsi")
		vsi := model.VSI{
			ID:        id,
			Name:      req.Name,
			Status:    "running",
			SubnetID:  sub.ID,
			VPCID:     sub.VPCID,
			Profile:   req.Profile.Name,
			Image:     req.Image.ID,
			PrimaryIP: fmt.Sprintf("10.%d.%d.%d", rand.Intn(256), rand.Intn(256), rand.Intn(254)+1),
			CreatedAt: time.Now().UTC(),
		}
		if err := s.CreateVSI(r.Context(), token, vsi); err != nil {
			writeError(w, http.StatusInternalServerError, "internal_error", err.Error())
			return
		}
		writeJSON(w, http.StatusCreated, vsi)
	}
}

func getVSI(s store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := tokenFromCtx(r.Context())
		id := chi.URLParam(r, "id")
		vsi, ok, err := s.GetVSI(r.Context(), token, id)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal_error", err.Error())
			return
		}
		if !ok {
			writeError(w, http.StatusNotFound, "not_found", "Instance not found")
			return
		}
		writeJSON(w, http.StatusOK, vsi)
	}
}

func deleteVSI(s store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := tokenFromCtx(r.Context())
		id := chi.URLParam(r, "id")
		_, ok, err := s.GetVSI(r.Context(), token, id)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal_error", err.Error())
			return
		}
		if !ok {
			writeError(w, http.StatusNotFound, "not_found", "Instance not found")
			return
		}
		if err := s.DeleteVSI(r.Context(), token, id); err != nil {
			writeError(w, http.StatusInternalServerError, "internal_error", err.Error())
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
