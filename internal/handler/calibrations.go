package handler

import (
	"github.com/jb843051627/kilnward/internal/model"
	"net/http"
)

func (h *Handler) calibrations(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		withJSON(w, func() (any, error) {
			return h.app.ListCalibrations(r.Context(), r.URL.Query().Get("kiln_id"), r.URL.Query().Get("sensor"))
		})
		return
	}
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var input model.Calibration
	if err := decodeJSON(r, &input); err != nil {
		writeAppError(w, err)
		return
	}
	withJSON(w, func() (any, error) { return h.app.PlanCalibration(r.Context(), input) })
}

func (h *Handler) calibrationDetail(w http.ResponseWriter, r *http.Request) {
	parts := splitPath(r.URL.Path, "/api/calibrations/")
	if len(parts) == 0 {
		http.NotFound(w, r)
		return
	}
	id := parts[0]
	if len(parts) == 1 && r.Method == http.MethodGet {
		withJSON(w, func() (any, error) { return h.app.GetCalibration(r.Context(), id) })
		return
	}
	if len(parts) != 2 || r.Method != http.MethodPost || parts[1] != "complete" {
		http.NotFound(w, r)
		return
	}
	var input struct {
		ObservedC float64 `json:"observed_c"`
	}
	if err := decodeJSON(r, &input); err != nil {
		writeAppError(w, err)
		return
	}
	withJSON(w, func() (any, error) { return h.app.CompleteCalibration(r.Context(), id, input.ObservedC) })
}
