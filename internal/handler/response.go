package handler

import (
	"encoding/json"
	"errors"
	"github.com/jb843051627/kilnward/internal/model"
	"net/http"
)

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]any{"error": message})
}

func writeAppError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	switch {
	case errors.Is(err, model.ErrNotFound):
		status = http.StatusNotFound
	case errors.Is(err, model.ErrValidation):
		status = http.StatusBadRequest
	case errors.Is(err, model.ErrInvalidState), errors.Is(err, model.ErrConflict):
		status = http.StatusConflict
	case errors.Is(err, model.ErrMaintenance):
		status = http.StatusLocked
	case errors.Is(err, model.ErrMissingReading):
		status = http.StatusPreconditionFailed
	}
	writeError(w, status, err.Error())
}

func withJSON(w http.ResponseWriter, fn func() (any, error)) {
	value, err := fn()
	if err != nil {
		writeAppError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, value)
}
