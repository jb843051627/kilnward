package service

import (
	"context"
	"fmt"
	"github.com/jb843051627/kilnward/internal/model"
	"github.com/jb843051627/kilnward/internal/policy"
	"github.com/jb843051627/kilnward/internal/validation"
)

func (a *App) RecordReading(ctx context.Context, reading model.Reading) (model.Reading, error) {
	if err := a.requireContext(ctx); err != nil {
		return reading, err
	}
	if reading.ID == "" {
		reading.ID = a.newID("reading")
	}
	if reading.RecordedAt.IsZero() {
		reading.RecordedAt = a.now()
	}
	if reading.Quality == "" {
		reading.Quality = model.QualityGood
	}
	if err := reading.Validate(); err != nil {
		return reading, err
	}
	if err := validation.Temperature(reading.Temperature); err != nil {
		return reading, err
	}
	if err := validation.Percentage(reading.Atmosphere); err != nil {
		return reading, err
	}
	cycle, err := a.repo.GetCycle(ctx, reading.CycleID)
	if err != nil {
		return reading, err
	}
	if cycle.LoadID != reading.LoadID {
		return reading, model.ErrConflict
	}
	if err := a.repo.AddReading(ctx, reading); err != nil {
		return reading, fmt.Errorf("save reading: %w", err)
	}
	a.metrics.Add("reading.received", 1)
	if err := a.audit(ctx, "cycle", reading.CycleID, "reading_received", reading.Sensor); err != nil {
		return reading, err
	}
	return reading, nil
}

func (a *App) RecordReadings(ctx context.Context, readings []model.Reading) error {
	if err := a.requireContext(ctx); err != nil {
		return err
	}
	if len(readings) == 0 {
		return model.ErrValidation
	}
	for i := range readings {
		if readings[i].ID == "" {
			readings[i].ID = a.newID("reading")
		}
		if readings[i].RecordedAt.IsZero() {
			readings[i].RecordedAt = a.now()
		}
		if readings[i].Quality == "" {
			readings[i].Quality = model.QualityGood
		}
	}
	if err := a.repo.AddReadings(ctx, readings); err != nil {
		return err
	}
	a.metrics.Add("reading.batch", int64(len(readings)))
	return nil
}

func (a *App) ListReadings(ctx context.Context, cycleID string, limit int) ([]model.Reading, error) {
	if err := a.requireContext(ctx); err != nil {
		return nil, err
	}
	return a.repo.ListReadings(ctx, cycleID, limit)
}

func (a *App) ReadingWindow(ctx context.Context, cycleID string) ([]model.Reading, error) {
	readings, err := a.repo.ListReadings(ctx, cycleID, 50)
	if err != nil {
		return nil, err
	}
	return readings, nil
}

func (a *App) SamplingStable(ctx context.Context, cycleID string) (bool, error) {
	cycle, err := a.repo.GetCycle(ctx, cycleID)
	if err != nil {
		return false, err
	}
	stage, err := a.repo.GetStage(ctx, cycleID, cycle.StageIndex)
	if err != nil {
		return false, err
	}
	readings, err := a.repo.LatestReadings(ctx, cycleID, nil)
	if err != nil {
		return false, err
	}
	if err := policy.DefaultSamplingPolicy().ValidateWindow(readings, a.now()); err != nil {
		return false, nil
	}
	for _, reading := range readings {
		if !reading.IsSafeFor(stage) {
			return false, nil
		}
	}
	return policy.DefaultSamplingPolicy().Stable(readings), nil
}
