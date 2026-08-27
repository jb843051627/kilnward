package handler

import (
	"github.com/jb843051627/kilnward/internal/model"
	"net/http"
)

func (h *Handler) incidents(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		withJSON(w, func() (any, error) {
			return h.app.ListIncidents(r.Context(), r.URL.Query().Get("kiln_id"), r.URL.Query().Get("status"))
		})
		return
	}
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var input model.Incident
	if err := decodeJSON(r, &input); err != nil {
		writeAppError(w, err)
		return
	}
	withJSON(w, func() (any, error) { return h.app.OpenIncident(r.Context(), input) })
}

func (h *Handler) incidentDetail(w http.ResponseWriter, r *http.Request) {
	parts := splitPath(r.URL.Path, "/api/incidents/")
	if len(parts) == 0 {
		http.NotFound(w, r)
		return
	}
	id := parts[0]
	if len(parts) == 1 && r.Method == http.MethodGet {
		withJSON(w, func() (any, error) { return h.app.GetIncident(r.Context(), id) })
		return
	}
	if len(parts) != 2 || r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}
	var input struct {
		Owner string `json:"owner"`
	}
	if err := decodeJSON(r, &input); err != nil {
		writeAppError(w, err)
		return
	}
	if parts[1] == "acknowledge" {
		withJSON(w, func() (any, error) { return h.app.AcknowledgeIncident(r.Context(), id, input.Owner) })
		return
	}
	if parts[1] == "resolve" {
		withJSON(w, func() (any, error) { return h.app.ResolveIncident(r.Context(), id, input.Owner) })
		return
	}
	http.NotFound(w, r)
}
