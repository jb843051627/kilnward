package handler

import (
	"github.com/jb843051627/kilnward/internal/model"
	"net/http"
)

type kilnInput struct {
	Name       string  `json:"name"`
	Location   string  `json:"location"`
	MaxTempC   float64 `json:"max_temp_c"`
	ProbeCount int     `json:"probe_count"`
}

func (h *Handler) kilns(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		state := model.KilnState(r.URL.Query().Get("state"))
		withJSON(w, func() (any, error) { return h.app.ListKilns(r.Context(), state) })
	case http.MethodPost:
		var input kilnInput
		if err := decodeJSON(r, &input); err != nil {
			writeAppError(w, err)
			return
		}
		withJSON(w, func() (any, error) {
			return h.app.RegisterKiln(r.Context(), model.Kiln{Name: input.Name, Location: input.Location, MaxTempC: input.MaxTempC, ProbeCount: input.ProbeCount})
		})
	default:
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (h *Handler) kilnDetail(w http.ResponseWriter, r *http.Request) {
	parts := splitPath(r.URL.Path, "/api/kilns/")
	if len(parts) == 0 {
		http.NotFound(w, r)
		return
	}
	id := parts[0]
	if len(parts) == 1 && r.Method == http.MethodGet {
		withJSON(w, func() (any, error) { return h.app.GetKiln(r.Context(), id) })
		return
	}
	if len(parts) == 2 && parts[1] == "state" && r.Method == http.MethodPost {
		var input struct {
			State model.KilnState `json:"state"`
		}
		if err := decodeJSON(r, &input); err != nil {
			writeAppError(w, err)
			return
		}
		withJSON(w, func() (any, error) { return h.app.ChangeKilnState(r.Context(), id, input.State) })
		return
	}
	if len(parts) == 2 && parts[1] == "quarantine" && r.Method == http.MethodPost {
		var input struct {
			Reason string `json:"reason"`
		}
		if err := decodeJSON(r, &input); err != nil {
			writeAppError(w, err)
			return
		}
		withJSON(w, func() (any, error) { return h.app.QuarantineKiln(r.Context(), id, input.Reason) })
		return
	}
	http.NotFound(w, r)
}
