package handler

import (
	"github.com/jb843051627/kilnward/internal/model"
	"net/http"
)

func (h *Handler) maintenance(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		withJSON(w, func() (any, error) {
			return h.app.ListMaintenance(r.Context(), r.URL.Query().Get("kiln_id"), queryBool(r, "active"))
		})
		return
	}
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var input model.Maintenance
	if err := decodeJSON(r, &input); err != nil {
		writeAppError(w, err)
		return
	}
	withJSON(w, func() (any, error) { return h.app.PlanMaintenance(r.Context(), input) })
}

func (h *Handler) maintenanceDetail(w http.ResponseWriter, r *http.Request) {
	parts := splitPath(r.URL.Path, "/api/maintenance/")
	if len(parts) == 0 {
		http.NotFound(w, r)
		return
	}
	id := parts[0]
	if len(parts) == 1 && r.Method == http.MethodGet {
		withJSON(w, func() (any, error) { return h.app.GetMaintenance(r.Context(), id) })
		return
	}
	if len(parts) != 2 || r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}
	switch parts[1] {
	case "start":
		withJSON(w, func() (any, error) { return h.app.StartMaintenance(r.Context(), id) })
	case "complete":
		withJSON(w, func() (any, error) { return h.app.CompleteMaintenance(r.Context(), id) })
	default:
		http.NotFound(w, r)
	}
}
