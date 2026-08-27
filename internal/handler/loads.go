package handler

import (
	"github.com/jb843051627/kilnward/internal/model"
	"net/http"
)

type loadInput struct {
	KilnID      string         `json:"kiln_id"`
	Label       string         `json:"label"`
	Material    model.Material `json:"material"`
	Profile     string         `json:"profile"`
	TargetTempC float64        `json:"target_temp_c"`
}

func (h *Handler) loads(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet, http.MethodPost) {
		return
	}
	if r.Method == http.MethodGet {
		kilnID := r.URL.Query().Get("kiln_id")
		state := model.LoadState(r.URL.Query().Get("state"))
		withJSON(w, func() (any, error) { return h.app.ListLoads(r.Context(), kilnID, state) })
		return
	}
	var input loadInput
	if err := decodeJSON(r, &input); err != nil {
		writeAppError(w, err)
		return
	}
	withJSON(w, func() (any, error) {
		return h.app.CreateLoad(r.Context(), model.Load{KilnID: input.KilnID, Label: input.Label, Material: input.Material, Profile: input.Profile, TargetTempC: input.TargetTempC})
	})
}

func (h *Handler) loadDetail(w http.ResponseWriter, r *http.Request) {
	parts := splitPath(r.URL.Path, "/api/loads/")
	if len(parts) == 0 {
		http.NotFound(w, r)
		return
	}
	id := parts[0]
	if len(parts) == 1 && r.Method == http.MethodGet {
		withJSON(w, func() (any, error) { return h.app.GetLoad(r.Context(), id) })
		return
	}
	if len(parts) != 2 || r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}
	switch parts[1] {
	case "attach":
		withJSON(w, func() (any, error) { return h.app.AttachLoad(r.Context(), id) })
	case "start":
		withJSON(w, func() (any, error) { return h.app.StartLoad(r.Context(), id) })
	case "reject":
		var input struct {
			Reason string `json:"reason"`
		}
		if err := decodeJSON(r, &input); err != nil {
			writeAppError(w, err)
			return
		}
		withJSON(w, func() (any, error) { return h.app.RejectLoad(r.Context(), id, input.Reason) })
	case "summary":
		withJSON(w, func() (any, error) { return h.app.LoadSummary(r.Context(), id) })
	default:
		http.NotFound(w, r)
	}
}
