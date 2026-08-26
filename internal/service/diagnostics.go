package service

import (
	"context"
	"github.com/jb843051627/kilnward/internal/model"
	"github.com/jb843051627/kilnward/internal/policy"
	"time"
)

func (a *App) DiagnoseKiln(ctx context.Context, kilnID string) (model.DiagnosticReport, error) {
	started := time.Now()
	kiln, err := a.repo.GetKiln(ctx, kilnID)
	if err != nil {
		return model.DiagnosticReport{}, err
	}
	incidents, err := a.repo.ListIncidents(ctx, kilnID, "")
	if err != nil {
		return model.DiagnosticReport{}, err
	}
	maintenance, err := a.repo.ListMaintenance(ctx, kilnID, false)
	if err != nil {
		return model.DiagnosticReport{}, err
	}
	report := model.DiagnosticReport{ID: a.newID("diagnostic"), KilnID: kilnID, GeneratedAt: a.now(), Healthy: true}
	for _, issue := range policy.DiagnoseKiln(kiln, incidents, maintenance, a.now()) {
		report = report.Add(issue)
	}
	report.DurationMS = time.Since(started).Milliseconds()
	_ = a.audit(ctx, "kiln", kilnID, "diagnosed", "operation health check")
	return report, nil
}

func (a *App) DiagnoseLoad(ctx context.Context, loadID string) (model.DiagnosticReport, error) {
	started := time.Now()
	load, cycle, stages, err := a.repo.LoadForSummary(ctx, loadID)
	if err != nil {
		return model.DiagnosticReport{}, err
	}
	readings, err := a.repo.ListReadings(ctx, cycle.ID, 200)
	if err != nil {
		return model.DiagnosticReport{}, err
	}
	gates, err := a.repo.GateReport(ctx, cycle.ID, cycle.StageIndex)
	if err != nil {
		return model.DiagnosticReport{}, err
	}
	report := model.DiagnosticReport{ID: a.newID("diagnostic"), KilnID: load.KilnID, GeneratedAt: a.now(), Healthy: true}
	for _, issue := range policy.DiagnoseLoad(load, cycle, stages, readings, gates, a.now()) {
		report = report.Add(issue)
	}
	report.DurationMS = time.Since(started).Milliseconds()
	_ = a.audit(ctx, "load", loadID, "diagnosed", "cycle health check")
	return report, nil
}

func (a *App) DiagnoseAll(ctx context.Context) ([]model.DiagnosticReport, error) {
	kilns, err := a.repo.ListKilns(ctx, "")
	if err != nil {
		return nil, err
	}
	reports := make([]model.DiagnosticReport, 0, len(kilns))
	for _, kiln := range kilns {
		report, err := a.DiagnoseKiln(ctx, kiln.ID)
		if err != nil {
			return nil, err
		}
		reports = append(reports, report)
	}
	return reports, nil
}

func (a *App) HasBlockingDiagnostics(ctx context.Context, kilnID string) (bool, error) {
	report, err := a.DiagnoseKiln(ctx, kilnID)
	if err != nil {
		return false, err
	}
	return !report.Healthy, nil
}
