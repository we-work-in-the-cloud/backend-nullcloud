package api

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
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
		r.Patch("/{id}", updateVSI(s))
		r.Delete("/{id}", deleteVSI(s))
		r.Post("/{id}/actions", vsiAction(s))
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
		existing, err := s.ListVSIs(r.Context(), token)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal_error", err.Error())
			return
		}
		usedIPs := make(map[string]bool, len(existing))
		for _, v := range existing {
			if v.SubnetID == sub.ID {
				usedIPs[v.PrimaryIP] = true
			}
		}
		primaryIP, err := allocateSubnetIP(sub.CIDRBlock, usedIPs)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal_error", err.Error())
			return
		}
		id := uid.New("vsi")
		vsi := model.VSI{
			ID:        id,
			Name:      req.Name,
			Status:    "running",
			CRN:       fmt.Sprintf("crn:nullcloud:instance:%s", id),
			SubnetID:  sub.ID,
			VPCID:     sub.VPCID,
			Profile:   req.Profile.Name,
			Image:     req.Image.ID,
			PrimaryIP: primaryIP,
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

func vsiAction(s store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := tokenFromCtx(r.Context())
		id := chi.URLParam(r, "id")
		var req struct {
			Type string `json:"type"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "bad_request", "invalid request body")
			return
		}
		var newStatus string
		switch req.Type {
		case "start", "restart":
			newStatus = "running"
		case "stop":
			newStatus = "stopped"
		default:
			writeError(w, http.StatusBadRequest, "bad_request", "type must be start, stop, or restart")
			return
		}
		vsi, ok, err := s.GetVSI(r.Context(), token, id)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal_error", err.Error())
			return
		}
		if !ok {
			writeError(w, http.StatusNotFound, "not_found", "Instance not found")
			return
		}
		if err := s.UpdateVSIStatus(r.Context(), token, id, newStatus); err != nil {
			writeError(w, http.StatusInternalServerError, "internal_error", err.Error())
			return
		}
		vsi.Status = newStatus
		writeJSON(w, http.StatusOK, vsi)
	}
}

func updateVSI(s store.Store) http.HandlerFunc {
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
		vsi, ok, err := s.GetVSI(r.Context(), token, id)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "internal_error", err.Error())
			return
		}
		if !ok {
			writeError(w, http.StatusNotFound, "not_found", "Instance not found")
			return
		}
		if err := s.RenameVSI(r.Context(), token, id, req.Name); err != nil {
			writeError(w, http.StatusInternalServerError, "internal_error", err.Error())
			return
		}
		vsi.Name = req.Name
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

// allocateSubnetIP picks a random unused host address within cidr.
func allocateSubnetIP(cidr string, usedIPs map[string]bool) (string, error) {
	_, network, err := net.ParseCIDR(cidr)
	if err != nil {
		return "", fmt.Errorf("invalid subnet CIDR %s: %w", cidr, err)
	}
	ip4 := network.IP.To4()
	if ip4 == nil {
		return "", fmt.Errorf("only IPv4 CIDRs are supported")
	}
	ones, _ := network.Mask.Size()
	hostBits := 32 - ones
	if hostBits < 2 {
		return "", fmt.Errorf("subnet %s is too small to allocate an IP", cidr)
	}
	totalHosts := (1 << hostBits) - 2 // exclude network address and broadcast
	base := uint32(ip4[0])<<24 | uint32(ip4[1])<<16 | uint32(ip4[2])<<8 | uint32(ip4[3])

	var available []string
	for i := uint32(1); i <= uint32(totalHosts); i++ {
		n := base + i
		ip := fmt.Sprintf("%d.%d.%d.%d", n>>24, (n>>16)&0xff, (n>>8)&0xff, n&0xff)
		if !usedIPs[ip] {
			available = append(available, ip)
		}
	}
	if len(available) == 0 {
		return "", fmt.Errorf("no available IPs in %s", cidr)
	}
	return available[rand.Intn(len(available))], nil
}
