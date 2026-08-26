package handler

import (
	"github.com/jb843051627/kilnward/internal/model"
	"net/http"
)

func (h *Handler) readingBatch(w http.ResponseWriter, r *http.Request, cycleID string) {
	var readings []model.Reading
	if err := decodeJSON(r, &readings); err != nil {
		writeAppError(w, err)
		return
	}
	for i := range readings {
		readings[i].CycleID = cycleID
	}
	if err := h.app.RecordReadings(r.Context(), readings); err != nil {
		writeAppError(w, err)
		return
	}
	writeJSON(w, http.StatusAccepted, map[string]any{"accepted": len(readings)})
}

func (h *Handler) readings(w http.ResponseWriter, r *http.Request) {
	parts := splitPath(r.URL.Path, "/api/cycles/")
	if len(parts) != 2 || parts[1] != "readings" {
		http.NotFound(w, r)
		return
	}
	if r.Method == http.MethodGet {
		withJSON(w, func() (any, error) { return h.app.ListReadings(r.Context(), parts[0], queryInt(r, "limit", 100)) })
		return
	}
	if r.Method == http.MethodPost {
		h.readingBatch(w, r, parts[0])
		return
	}
	writeError(w, http.StatusMethodNotAllowed, "method not allowed")
}
