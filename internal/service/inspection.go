package service

import (
	"context"
	"fmt"
	"github.com/jb843051627/kilnward/internal/model"
	"github.com/jb843051627/kilnward/internal/policy"
)

type InspectionReport struct {
	KilnID  string                  `json:"kiln_id"`
	LoadID  string                  `json:"load_id"`
	Passed  bool                    `json:"passed"`
	Summary string                  `json:"summary"`
	Issues  []model.DiagnosticIssue `json:"issues"`
}

func (a *App) InspectLoad(ctx context.Context, loadID string) (InspectionReport, error) {
	load, cycle, stages, err := a.repo.LoadForSummary(ctx, loadID)
	if err != nil {
		return InspectionReport{}, err
	}
	kiln, err := a.repo.GetKiln(ctx, load.KilnID)
	if err != nil {
		return InspectionReport{}, err
	}
	readings, err := a.repo.ListReadings(ctx, cycle.ID, 200)
	if err != nil {
		return InspectionReport{}, err
	}
	gate, err := a.repo.GateReport(ctx, cycle.ID, cycle.StageIndex)
	if err != nil {
		return InspectionReport{}, err
	}
	loadResult := policy.InspectLoad(load, kiln, cycle, readings, a.now())
	profileResult := policy.InspectionResult{Passed: gate.Passed}
	if len(stages) == 0 {
		profileResult.Passed = false
		profileResult.Issues = append(profileResult.Issues, model.Issue("NO_STAGES", model.DiagnosticFailure, cycle.ID, "热循环没有阶段", "stage_count=0"))
	}
	passed, summary := policy.InspectionSummary(loadResult, profileResult)
	issues := append(append([]model.DiagnosticIssue{}, loadResult.Issues...), profileResult.Issues...)
	_ = a.audit(ctx, "load", loadID, "inspected", summary)
	return InspectionReport{KilnID: load.KilnID, LoadID: loadID, Passed: passed, Summary: summary, Issues: issues}, nil
}

func (a *App) InspectKiln(ctx context.Context, kilnID string) (InspectionReport, error) {
	kiln, err := a.repo.GetKiln(ctx, kilnID)
	if err != nil {
		return InspectionReport{}, err
	}
	readings := make([]model.Reading, 0)
	result := policy.InspectProbe(kiln, readings, a.now())
	passed, summary := policy.InspectionSummary(result)
	return InspectionReport{KilnID: kilnID, Passed: passed, Summary: summary, Issues: result.Issues}, nil
}

func (a *App) InspectProfileForKiln(ctx context.Context, profileID, kilnID string) (InspectionReport, error) {
	profile, err := a.repo.GetProfile(ctx, profileID)
	if err != nil {
		return InspectionReport{}, err
	}
	kiln, err := a.repo.GetKiln(ctx, kilnID)
	if err != nil {
		return InspectionReport{}, err
	}
	result := policy.InspectProfile(profile, kiln)
	passed, summary := policy.InspectionSummary(result)
	return InspectionReport{KilnID: kilnID, Summary: summary, Passed: passed, Issues: result.Issues}, nil
}

func (a *App) InspectCalibration(ctx context.Context, calibrationID string) (InspectionReport, error) {
	item, err := a.repo.GetCalibration(ctx, calibrationID)
	if err != nil {
		return InspectionReport{}, err
	}
	level := model.DiagnosticOK
	message := "探头漂移在允许范围内"
	if !item.WithinTolerance() {
		level, message = model.DiagnosticFailure, "探头漂移超过允许范围"
	}
	issue := model.Issue("CALIBRATION_DRIFT", level, item.Sensor, message, fmt.Sprintf("drift=%.2f", item.Drift()))
	issues := []model.DiagnosticIssue{}
	if level != model.DiagnosticOK {
		issues = append(issues, issue)
	}
	return InspectionReport{KilnID: item.KilnID, Passed: level == model.DiagnosticOK, Summary: message, Issues: issues}, nil
}
