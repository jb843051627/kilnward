package service

import (
	"context"
	"fmt"
	"github.com/jb843051627/kilnward/internal/model"
	"github.com/jb843051627/kilnward/internal/policy"
)

func (a *App) PlanCalibration(ctx context.Context, item model.Calibration) (model.Calibration, error) {
	if err := a.requireContext(ctx); err != nil {
		return item, err
	}
	if item.ID == "" {
		item.ID = a.newID("calibration")
	}
	if item.Status == "" {
		item.Status = model.CalibrationPending
	}
	if item.CreatedAt.IsZero() {
		item.CreatedAt = a.now()
	}
	kiln, err := a.repo.GetKiln(ctx, item.KilnID)
	if err != nil {
		return item, err
	}
	decision := policy.StartCalibration(kiln, item.Sensor)
	if !decision.Allowed {
		return item, fmt.Errorf("calibration blocked: %s", decision.Reason)
	}
	if err := item.Validate(); err != nil {
		return item, err
	}
	if err := a.repo.CreateCalibration(ctx, item); err != nil {
		return item, err
	}
	_ = a.audit(ctx, "calibration", item.ID, "planned", item.Sensor)
	return item, nil
}

func (a *App) GetCalibration(ctx context.Context, id string) (model.Calibration, error) {
	if err := a.requireContext(ctx); err != nil {
		return model.Calibration{}, err
	}
	return a.repo.GetCalibration(ctx, id)
}
func (a *App) ListCalibrations(ctx context.Context, kilnID, sensor string) ([]model.Calibration, error) {
	if err := a.requireContext(ctx); err != nil {
		return nil, err
	}
	return a.repo.ListCalibrations(ctx, kilnID, sensor)
}

func (a *App) CompleteCalibration(ctx context.Context, id string, observed float64) (model.Calibration, error) {
	item, err := a.repo.GetCalibration(ctx, id)
	if err != nil {
		return item, err
	}
	item.ObservedC, item.Status = observed, model.CalibrationRunning
	if err := item.Complete(a.now()); err != nil {
		return item, err
	}
	if err := a.repo.UpdateCalibration(ctx, item); err != nil {
		return item, err
	}
	if item.Status == model.CalibrationFailed {
		_, _ = a.OpenIncident(ctx, model.Incident{KilnID: item.KilnID, Code: "probe-drift", Severity: policy.DriftSeverity(item), Detail: "probe calibration drift exceeded tolerance", Status: model.IncidentOpen})
	}
	_ = a.audit(ctx, "calibration", id, "completed", string(item.Status))
	return a.repo.GetCalibration(ctx, id)
}

func (a *App) CalibrationRequired(ctx context.Context, kilnID, sensor string) (bool, error) {
	kiln, err := a.repo.GetKiln(ctx, kilnID)
	if err != nil {
		return false, err
	}
	items, err := a.repo.ListCalibrations(ctx, kilnID, sensor)
	if err != nil {
		return false, err
	}
	stage := model.Stage{TargetTempC: kiln.MaxTempC}
	return !policy.RequiredForStage(stage, items, a.now()), nil
}
