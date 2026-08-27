package service

import (
	"context"
	"fmt"
	"github.com/jb843051627/kilnward/internal/model"
	"github.com/jb843051627/kilnward/internal/policy"
	"time"
)

func (a *App) EvaluateGate(ctx context.Context, cycleID string, stageSeq int, operator string) (model.GateReport, error) {
	if err := a.requireContext(ctx); err != nil {
		return model.GateReport{}, err
	}
	if operator == "" {
		return model.GateReport{}, model.ErrValidation
	}
	cycle, err := a.repo.GetCycle(ctx, cycleID)
	if err != nil {
		return model.GateReport{}, err
	}
	stage, err := a.repo.GetStage(ctx, cycleID, stageSeq)
	if err != nil {
		return model.GateReport{}, err
	}
	readings, err := a.repo.ListReadings(ctx, cycleID, 100)
	if err != nil {
		return model.GateReport{}, err
	}
	checks := make([]model.Gate, 0)
	now := a.now()
	required := policy.RequiredChecks(stage)
	for _, kind := range required {
		passed, reason := a.evaluateCheck(kind, stage, readings, cycle, now)
		checks = append(checks, model.Gate{ID: a.newID("gate"), CycleID: cycleID, StageSeq: stageSeq, Name: string(kind), Kind: kind, Passed: passed, Reason: reason, CheckedAt: now, CheckedBy: operator})
	}
	for _, gate := range checks {
		if err := a.repo.SaveGate(ctx, gate); err != nil {
			return model.GateReport{}, err
		}
	}
	report := model.GateReport{CycleID: cycleID, Stage: stageSeq, Checks: checks, Passed: policy.Passed(checks)}
	_ = a.audit(ctx, "cycle", cycleID, "gate_evaluated", fmt.Sprintf("stage=%d passed=%t", stageSeq, report.Passed))
	return report, nil
}

func (a *App) evaluateCheck(kind model.GateKind, stage model.Stage, readings []model.Reading, cycle model.Cycle, now time.Time) (bool, string) {
	switch kind {
	case model.GateTemperature:
		for _, reading := range readings {
			if reading.CycleID == cycle.ID && reading.IsFresh(now, 45*time.Second) && stage.WithinTemperature(reading.Temperature) {
				return true, "recent temperature accepted"
			}
		}
		return false, "no recent temperature in tolerance"
	case model.GateSampling:
		return len(readings) >= policy.DefaultSamplingPolicy().MinPerStage, fmt.Sprintf("%d readings", len(readings))
	case model.GateCalibration:
		for _, reading := range readings {
			if reading.Quality != model.QualityGood {
				return false, "suspect reading"
			}
		}
		return len(readings) > 0, "all readings have good quality"
	case model.GateSequence:
		return cycle.StageIndex == stage.Sequence, "cycle stage matches gate"
	default:
		return false, "unknown gate"
	}
}

func (a *App) GateReport(ctx context.Context, cycleID string, stage int) (model.GateReport, error) {
	if err := a.requireContext(ctx); err != nil {
		return model.GateReport{}, err
	}
	return a.repo.GateReport(ctx, cycleID, stage)
}
