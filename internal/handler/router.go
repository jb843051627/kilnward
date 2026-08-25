package handler

import (
	"github.com/jb843051627/kilnward/internal/service"
	"net/http"
	"strings"
)

type Handler struct {
	app *service.App
	mux *http.ServeMux
}

func New(app *service.App) http.Handler {
	h := &Handler{app: app, mux: http.NewServeMux()}
	h.mux.HandleFunc("/", h.home)
	h.mux.HandleFunc("/health", h.health)
	h.mux.HandleFunc("/api/kilns", h.kilns)
	h.mux.HandleFunc("/api/kilns/", h.kilnDetail)
	h.mux.HandleFunc("/api/loads", h.loads)
	h.mux.HandleFunc("/api/loads/", h.loadDetail)
	h.mux.HandleFunc("/api/cycles/", h.cycleDetail)
	h.mux.HandleFunc("/api/incidents", h.incidents)
	h.mux.HandleFunc("/api/incidents/", h.incidentDetail)
	h.mux.HandleFunc("/api/maintenance", h.maintenance)
	h.mux.HandleFunc("/api/maintenance/", h.maintenanceDetail)
	h.mux.HandleFunc("/api/summary", h.summary)
	h.mux.HandleFunc("/api/profiles", h.profiles)
	h.mux.HandleFunc("/api/profiles/", h.profileDetail)
	h.mux.HandleFunc("/api/telemetry/", h.telemetry)
	h.mux.HandleFunc("/api/commands", h.commands)
	h.mux.HandleFunc("/api/commands/", h.commandDetail)
	h.mux.HandleFunc("/api/reports/", h.reportDetail)
	h.mux.HandleFunc("/api/calibrations", h.calibrations)
	h.mux.HandleFunc("/api/calibrations/", h.calibrationDetail)
	h.mux.HandleFunc("/api/diagnostics/", h.diagnosticDetail)
	return logging(h.mux)
}

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Kilnward-Trace", "active")
		next.ServeHTTP(w, r)
	})
}

func splitPath(path, prefix string) []string {
	value := strings.TrimPrefix(path, prefix)
	value = strings.Trim(value, "/")
	if value == "" {
		return nil
	}
	return strings.Split(value, "/")
}

func method(w http.ResponseWriter, r *http.Request, allowed ...string) bool {
	for _, item := range allowed {
		if r.Method == item {
			return true
		}
	}
	writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	return false
}
