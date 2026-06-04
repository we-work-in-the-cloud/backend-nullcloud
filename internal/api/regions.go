package api

import (
	"net/http"

	"github.com/we-work-in-the-cloud/nullcloud/backend/internal/model"
)

func listRegions() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{"regions": model.Regions()})
	}
}
