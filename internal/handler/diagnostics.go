package handler

import "net/http"

func (h *Handler) diagnosticDetail(w http.ResponseWriter, r *http.Request) {
	parts := splitPath(r.URL.Path, "/api/diagnostics/")
	if len(parts) != 2 || r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}
	switch parts[0] {
	case "kilns":
		withJSON(w, func() (any, error) { return h.app.DiagnoseKiln(r.Context(), parts[1]) })
	case "loads":
		withJSON(w, func() (any, error) { return h.app.DiagnoseLoad(r.Context(), parts[1]) })
	default:
		http.NotFound(w, r)
	}
}
