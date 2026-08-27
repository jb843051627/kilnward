package handler

import (
	"github.com/jb843051627/kilnward/internal/model"
	"net/http"
)

func (h *Handler) profiles(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		withJSON(w, func() (any, error) {
			return h.app.ListProfiles(r.Context(), model.ProfileStatus(r.URL.Query().Get("status")))
		})
		return
	}
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}
	var input model.Profile
	if err := decodeJSON(r, &input); err != nil {
		writeAppError(w, err)
		return
	}
	withJSON(w, func() (any, error) { return h.app.CreateProfile(r.Context(), input) })
}

func (h *Handler) profileDetail(w http.ResponseWriter, r *http.Request) {
	parts := splitPath(r.URL.Path, "/api/profiles/")
	if len(parts) == 0 {
		http.NotFound(w, r)
		return
	}
	id := parts[0]
	if len(parts) == 1 && r.Method == http.MethodGet {
		withJSON(w, func() (any, error) { return h.app.GetProfile(r.Context(), id) })
		return
	}
	if len(parts) != 2 || r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}
	switch parts[1] {
	case "review":
		withJSON(w, func() (any, error) { return h.app.ReviewProfile(r.Context(), id) })
	case "publish":
		withJSON(w, func() (any, error) { return h.app.PublishProfile(r.Context(), id) })
	case "retire":
		withJSON(w, func() (any, error) { return h.app.RetireProfile(r.Context(), id) })
	default:
		http.NotFound(w, r)
	}
}
