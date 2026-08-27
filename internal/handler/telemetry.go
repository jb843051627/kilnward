package handler

import (
	"github.com/jb843051627/kilnward/internal/model"
	"net/http"
)

func (h *Handler) telemetry(w http.ResponseWriter, r *http.Request) {
	parts := splitPath(r.URL.Path, "/api/telemetry/")
	if len(parts) != 1 {
		http.NotFound(w, r)
		return
	}
	cycleID := parts[0]
	if r.Method == http.MethodGet {
		withJSON(w, func() (any, error) { return h.app.ListTelemetry(r.Context(), cycleID, queryInt(r, "limit", 50)) })
		return
	}
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var frame model.TelemetryFrame
	if err := decodeJSON(r, &frame); err != nil {
		writeAppError(w, err)
		return
	}
	frame.CycleID = cycleID
	withJSON(w, func() (any, error) { return h.app.IngestTelemetry(r.Context(), frame) })
}
