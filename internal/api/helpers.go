package api

import (
	"encoding/json"
	"net/http"
)

type apiError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type apiErrors struct {
	Errors []apiError `json:"errors"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, apiErrors{Errors: []apiError{{Code: code, Message: message}}})
}
