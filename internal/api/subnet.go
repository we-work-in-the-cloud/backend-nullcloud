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

func subnetRoutes(s store.Store) func(chi.Router) {
	return func(r chi.Router) {
		r.Get("/", listSubnets(s))
		r.Post("/", createSubnet(s))
		r.Get("/{id}", getSubnet(s))
		r.Delete("/{id}", deleteSubnet(s))
	}
}

func listSubnets(s store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := tokenFromCtx(r.Context())
		subnets, err := s.ListSubnets(r.Context(), token)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal_error", err.Error())
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"subnets": subnets})
	}
}

func createSubnet(s store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := tokenFromCtx(r.Context())
		var req struct {
			Name string `json:"name"`
			VPC  struct {
				ID string `json:"id"`
			} `json:"vpc"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Name == "" || req.VPC.ID == "" {
			writeError(w, http.StatusBadRequest, "bad_request", "name and vpc.id are required")
			return
		}
		_, ok, err := s.GetVPC(r.Context(), token, req.VPC.ID)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal_error", err.Error())
			return
		}
		if !ok {
			writeError(w, http.StatusNotFound, "not_found", "VPC not found")
			return
		}
		id := uid.New("subnet")
		sub := model.Subnet{
			ID:        id,
			Name:      req.Name,
			Status:    "available",
			CRN:       fmt.Sprintf("crn:nullcloud:subnet:%s", id),
			VPCID:     req.VPC.ID,
			CIDRBlock: fmt.Sprintf("10.%d.%d.0/24", rand.Intn(256), rand.Intn(256)),
			CreatedAt: time.Now().UTC(),
		}
		if err := s.CreateSubnet(r.Context(), token, sub); err != nil {
			writeError(w, http.StatusInternalServerError, "internal_error", err.Error())
			return
		}
		writeJSON(w, http.StatusCreated, sub)
	}
}

func getSubnet(s store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := tokenFromCtx(r.Context())
		id := chi.URLParam(r, "id")
		sub, ok, err := s.GetSubnet(r.Context(), token, id)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal_error", err.Error())
			return
		}
		if !ok {
			writeError(w, http.StatusNotFound, "not_found", "Subnet not found")
			return
		}
		writeJSON(w, http.StatusOK, sub)
	}
}

func deleteSubnet(s store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := tokenFromCtx(r.Context())
		id := chi.URLParam(r, "id")
		_, ok, err := s.GetSubnet(r.Context(), token, id)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal_error", err.Error())
			return
		}
		if !ok {
			writeError(w, http.StatusNotFound, "not_found", "Subnet not found")
			return
		}
		if err := s.DeleteSubnet(r.Context(), token, id); err != nil {
			writeError(w, http.StatusInternalServerError, "internal_error", err.Error())
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
