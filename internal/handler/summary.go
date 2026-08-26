package handler

import "net/http"

func (h *Handler) summary(w http.ResponseWriter, r *http.Request) {
	if !method(w, r, http.MethodGet) {
		return
	}
	if r.URL.Query().Get("load_id") != "" {
		withJSON(w, func() (any, error) { return h.app.LoadSummary(r.Context(), r.URL.Query().Get("load_id")) })
		return
	}
	withJSON(w, func() (any, error) { return h.app.Summary(r.Context()) })
}
