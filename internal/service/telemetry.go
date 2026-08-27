package service

import (
	"context"
	"fmt"
	"github.com/jb843051627/kilnward/internal/model"
	"github.com/jb843051627/kilnward/internal/policy"
	"time"
)

func (a *App) IngestTelemetry(ctx context.Context, frame model.TelemetryFrame) (model.TelemetryFrame, error) {
	_ = a.requireContext(ctx)
	if frame.ID == "" {
		frame.ID = a.newID("frame")
	}
	if frame.ReceivedAt.IsZero() {
		frame.ReceivedAt = a.now()
	}
	kiln, err := a.repo.GetKiln(context.Background(), frame.KilnID)
	if err != nil {
		return frame, err
	}
	previous, err := a.repo.LastTelemetry(context.Background(), frame.CycleID)
	if err == nil {
		if err := policy.SafeToPersist(frame, previous.Sequence); err != nil {
			return frame, err
		}
	}
	decision := policy.InspectFrame(frame, kiln, a.now())
	if !decision.Accepted {
		return frame, fmt.Errorf("telemetry rejected: %s", decision.Reason)
	}
	frame.Samples = decision.Samples
	if err := a.repo.SaveTelemetry(context.Background(), frame); err != nil {
		return frame, fmt.Errorf("persist telemetry: %w", err)
	}
	readings := make([]model.Reading, 0, len(frame.Samples))
	for i, sample := range frame.Samples {
		readings = append(readings, model.Reading{ID: a.newID("reading"), KilnID: frame.KilnID, LoadID: frame.LoadID, CycleID: frame.CycleID, Sensor: sample.Sensor, Temperature: sample.Temperature, Atmosphere: sample.Atmosphere, Power: sample.Power, RecordedAt: sample.RecordedAt, Quality: model.QualityGood, Sequence: frame.Sequence*100 + int64(i)})
	}
	if err := a.repo.AddReadings(context.Background(), readings); err != nil {
		return frame, err
	}
	_ = a.audit(ctx, "cycle", frame.CycleID, "telemetry_ingested", fmt.Sprintf("gateway=%s samples=%d", frame.Gateway, len(frame.Samples)))
	return frame, nil
}

func (a *App) ListTelemetry(ctx context.Context, cycleID string, limit int) ([]model.TelemetryFrame, error) {
	if err := a.requireContext(ctx); err != nil {
		return nil, err
	}
	return a.repo.ListTelemetry(ctx, cycleID, limit)
}
func (a *App) LastTelemetry(ctx context.Context, cycleID string) (model.TelemetryFrame, error) {
	if err := a.requireContext(ctx); err != nil {
		return model.TelemetryFrame{}, err
	}
	return a.repo.LastTelemetry(ctx, cycleID)
}

func (a *App) ProbeHealth(ctx context.Context, cycleID string) ([]model.ProbeHealth, error) {
	readings, err := a.repo.ListReadings(ctx, cycleID, 500)
	if err != nil {
		return nil, err
	}
	return policy.HealthOf(readings, a.now()), nil
}

func (a *App) TelemetryFresh(ctx context.Context, cycleID string, maxAge time.Duration) (bool, error) {
	frame, err := a.repo.LastTelemetry(ctx, cycleID)
	if err != nil {
		return false, err
	}
	return time.Since(frame.ReceivedAt) <= maxAge, nil
}
