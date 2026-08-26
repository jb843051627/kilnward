package policy

import (
	"github.com/jb843051627/kilnward/internal/model"
	"sort"
)

type ReportFilter struct {
	KilnID    string
	LoadState model.LoadState
	Limit     int
}

func NormalizeFilter(filter ReportFilter) ReportFilter {
	if filter.Limit <= 0 {
		filter.Limit = 50
	}
	if filter.Limit > 500 {
		filter.Limit = 500
	}
	return filter
}

func RankIncidents(items []model.Incident) []model.Incident {
	out := append([]model.Incident(nil), items...)
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].Severity != out[j].Severity {
			return severityRank(out[i].Severity) > severityRank(out[j].Severity)
		}
		return out[i].OpenedAt.Before(out[j].OpenedAt)
	})
	return out
}

func severityRank(value model.IncidentSeverity) int {
	switch value {
	case model.SeverityCritical:
		return 3
	case model.SeverityWarning:
		return 2
	default:
		return 1
	}
}

func SummaryHealthy(kiln model.Kiln, incidents []model.Incident, maintenance []model.Maintenance) bool {
	if kiln.State == model.KilnQuarantined || kiln.State == model.KilnMaintenance {
		return false
	}
	for _, incident := range incidents {
		if incident.Status != model.IncidentResolved && incident.Severity == model.SeverityCritical {
			return false
		}
	}
	for _, item := range maintenance {
		if item.Status == model.MaintenanceActive {
			return false
		}
	}
	return true
}
