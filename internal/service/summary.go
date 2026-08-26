package service

import (
	"context"
	"github.com/jb843051627/kilnward/internal/model"
	"github.com/jb843051627/kilnward/internal/store"
)

func (a *App) Summary(ctx context.Context) (store.SummaryCounts, error) {
	if err := a.requireContext(ctx); err != nil {
		return store.SummaryCounts{}, err
	}
	return a.repo.SummaryCounts(ctx)
}

func (a *App) LoadSummary(ctx context.Context, loadID string) (model.LoadSummary, error) {
	load, cycle, stages, err := a.repo.LoadForSummary(ctx, loadID)
	if err != nil {
		return model.LoadSummary{}, err
	}
	readings, err := a.repo.ListReadings(ctx, cycle.ID, 200)
	if err != nil {
		return model.LoadSummary{}, err
	}
	incidents, err := a.repo.ListIncidents(ctx, load.KilnID, "")
	if err != nil {
		return model.LoadSummary{}, err
	}
	report, err := a.repo.GateReport(ctx, cycle.ID, cycle.StageIndex)
	if err != nil {
		return model.LoadSummary{}, err
	}
	snapshots := make([]model.StageSnapshot, 0, len(stages))
	for _, stage := range stages {
		snapshots = append(snapshots, stageSnapshot(stage, readings))
	}
	healthy := report.Passed
	for _, incident := range incidents {
		if incident.LoadID == load.ID && incident.Status != model.IncidentResolved {
			healthy = false
		}
	}
	return model.LoadSummary{Load: load, Cycle: cycle, Stages: snapshots, Incidents: incidents, GateReport: report, Healthy: healthy}, nil
}

func stageSnapshot(stage model.Stage, readings []model.Reading) model.StageSnapshot {
	snapshot := model.StageSnapshot{Name: stage.Name, Status: string(stage.Status), TargetTempC: stage.TargetTempC}
	for _, reading := range readings {
		if reading.RecordedAt.After(snapshot.LastRecordedAt) {
			snapshot.LastRecordedAt, snapshot.LastTempC = reading.RecordedAt, reading.Temperature
		}
		if stage.WithinTemperature(reading.Temperature) {
			snapshot.ReadingCount++
		}
	}
	return snapshot
}

func (a *App) MetricsSnapshot() map[string]int64 {
	out := a.metrics.Snapshot()
	out["exported_at"] = a.now().Unix()
	return out
}
