package handler

import "net/http"

func (h *Handler) reportDetail(w http.ResponseWriter, r *http.Request) {
	parts := splitPath(r.URL.Path, "/api/reports/")
	if len(parts) != 2 || r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}
	if parts[0] == "kilns" && parts[1] != "" {
		withJSON(w, func() (any, error) { return h.app.OperationReport(r.Context(), parts[1], queryInt(r, "limit", 50)) })
		return
	}
	http.NotFound(w, r)
}
