package handler

import "net/http"

func (h *Handler) cycleDetail(w http.ResponseWriter, r *http.Request) {
	parts := splitPath(r.URL.Path, "/api/cycles/")
	if len(parts) == 0 {
		http.NotFound(w, r)
		return
	}
	id := parts[0]
	if len(parts) == 1 && r.Method == http.MethodGet {
		withJSON(w, func() (any, error) { return h.app.GetCycle(r.Context(), id) })
		return
	}
	if len(parts) == 2 && parts[1] == "readings" {
		if r.Method == http.MethodGet {
			withJSON(w, func() (any, error) { return h.app.ListReadings(r.Context(), id, queryInt(r, "limit", 100)) })
			return
		}
		if r.Method == http.MethodPost {
			h.readingBatch(w, r, id)
			return
		}
	}
	if len(parts) != 2 || r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}
	switch parts[1] {
	case "advance":
		withJSON(w, func() (any, error) { return h.app.AdvanceCycle(r.Context(), id) })
	case "hold":
		withJSON(w, func() (any, error) { return h.app.HoldCycle(r.Context(), id) })
	case "gates":
		withJSON(w, func() (any, error) { return h.app.EvaluateGate(r.Context(), id, 0, "operator") })
	default:
		http.NotFound(w, r)
	}
}
