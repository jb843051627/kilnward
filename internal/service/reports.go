package service

import (
	"context"
	"github.com/jb843051627/kilnward/internal/model"
	"github.com/jb843051627/kilnward/internal/policy"
	"sort"
)

type OperationReport struct {
	Kiln        model.Kiln              `json:"kiln"`
	Loads       []model.Load            `json:"loads"`
	Incidents   []model.Incident        `json:"incidents"`
	Maintenance []model.Maintenance     `json:"maintenance"`
	Commands    []model.OperatorCommand `json:"commands"`
	Healthy     bool                    `json:"healthy"`
}

func (a *App) OperationReport(ctx context.Context, kilnID string, limit int) (OperationReport, error) {
	kiln, err := a.repo.GetKiln(ctx, kilnID)
	if err != nil {
		return OperationReport{}, err
	}
	loads, err := a.repo.ListLoads(ctx, kilnID, "")
	if err != nil {
		return OperationReport{}, err
	}
	incidents, err := a.repo.ListIncidents(ctx, kilnID, "")
	if err != nil {
		return OperationReport{}, err
	}
	maintenance, err := a.repo.ListMaintenance(ctx, kilnID, false)
	if err != nil {
		return OperationReport{}, err
	}
	commands, err := a.repo.ListCommands(ctx, kilnID, "")
	if err != nil {
		return OperationReport{}, err
	}
	if limit <= 0 {
		limit = 50
	}
	if len(loads) > limit {
		loads = loads[:limit]
	}
	if len(incidents) > limit {
		incidents = incidents[:limit]
	}
	if len(maintenance) > limit {
		maintenance = maintenance[:limit]
	}
	if len(commands) > limit {
		commands = commands[:limit]
	}
	sort.Slice(loads, func(i, j int) bool { return loads[i].UpdatedAt.After(loads[j].UpdatedAt) })
	incidents = policy.RankIncidents(incidents)
	return OperationReport{Kiln: kiln, Loads: loads, Incidents: incidents, Maintenance: maintenance, Commands: commands, Healthy: policy.SummaryHealthy(kiln, incidents, maintenance)}, nil
}

func (a *App) RecentEvents(ctx context.Context, kilnID string, limit int) ([]model.TimelineEntry, error) {
	return a.AuditForKiln(ctx, kilnID, limit)
}
