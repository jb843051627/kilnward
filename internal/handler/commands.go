package handler

import (
	"github.com/jb843051627/kilnward/internal/model"
	"net/http"
)

func (h *Handler) commands(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		withJSON(w, func() (any, error) {
			return h.app.ListCommands(r.Context(), r.URL.Query().Get("kiln_id"), r.URL.Query().Get("status"))
		})
		return
	}
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var input model.OperatorCommand
	if err := decodeJSON(r, &input); err != nil {
		writeAppError(w, err)
		return
	}
	withJSON(w, func() (any, error) { return h.app.IssueCommand(r.Context(), input) })
}

func (h *Handler) commandDetail(w http.ResponseWriter, r *http.Request) {
	parts := splitPath(r.URL.Path, "/api/commands/")
	if len(parts) == 0 {
		http.NotFound(w, r)
		return
	}
	id := parts[0]
	if len(parts) == 1 && r.Method == http.MethodGet {
		withJSON(w, func() (any, error) { return h.app.GetCommand(r.Context(), id) })
		return
	}
	if len(parts) != 2 || r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}
	switch parts[1] {
	case "run":
		withJSON(w, func() (any, error) { return h.app.RunCommand(r.Context(), id) })
	case "cancel":
		withJSON(w, func() (any, error) { return h.app.CancelCommand(r.Context(), id) })
	case "retry":
		withJSON(w, func() (any, error) { return h.app.RetryCommand(r.Context(), id) })
	default:
		http.NotFound(w, r)
	}
}
