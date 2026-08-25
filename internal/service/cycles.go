package service

import (
	"context"
	"fmt"
	"github.com/jb843051627/kilnward/internal/model"
	"github.com/jb843051627/kilnward/internal/policy"
	"time"
)

func (a *App) GetCycle(ctx context.Context, id string) (model.Cycle, error) {
	if err := a.requireContext(ctx); err != nil {
		return model.Cycle{}, err
	}
	return a.repo.GetCycle(ctx, id)
}

func (a *App) CycleForLoad(ctx context.Context, loadID string) (model.Cycle, error) {
	if err := a.requireContext(ctx); err != nil {
		return model.Cycle{}, err
	}
	return a.repo.GetCycleByLoad(ctx, loadID)
}

func (a *App) StartLoad(ctx context.Context, loadID string) (model.Load, error) {
	load, err := a.repo.GetLoad(ctx, loadID)
	if err != nil {
		return load, err
	}
	kiln, err := a.repo.GetKiln(ctx, load.KilnID)
	if err != nil {
		return load, err
	}
	cycle, err := a.repo.GetCycleByLoad(ctx, loadID)
	if err != nil {
		return load, err
	}
	report, err := a.EvaluateGate(ctx, cycle.ID, 0, "operator")
	if err != nil {
		return load, err
	}
	if err := policy.StartAllowed(kiln, load, report); err != nil {
		return load, err
	}
	now := a.now()
	load.State, load.StartedAt, load.UpdatedAt = model.LoadRunning, &now, now
	cycle.Status, cycle.StartedAt, cycle.UpdatedAt = model.CycleActive, &now, now
	if err := a.repo.UpdateLoad(ctx, load); err != nil {
		return load, err
	}
	if err := a.repo.UpdateCycle(ctx, cycle); err != nil {
		return load, err
	}
	if err := a.repo.StartStage(ctx, cycle.ID, 0, now.Format(time.RFC3339Nano)); err != nil {
		return load, err
	}
	if err := a.audit(ctx, "load", loadID, "started", cycle.ID); err != nil {
		return load, err
	}
	a.metrics.Add("load.started", 1)
	return a.repo.GetLoad(ctx, loadID)
}

func (a *App) AdvanceCycle(ctx context.Context, cycleID string) (model.Cycle, error) {
	cycle, err := a.repo.GetCycle(ctx, cycleID)
	if err != nil {
		return cycle, err
	}
	stages, err := a.repo.ListStages(ctx, cycleID)
	if err != nil {
		return cycle, err
	}
	if cycle.StageIndex < 0 || cycle.StageIndex >= len(stages) {
		return cycle, model.ErrInvalidState
	}
	if !cycle.CanAdvance(len(stages)) {
		return cycle, model.ErrInvalidState
	}
	current := stages[cycle.StageIndex]
	if err := policy.CompleteStage(&current, a.now()); err != nil {
		return cycle, err
	}
	if err := a.repo.UpdateStage(ctx, current); err != nil {
		return cycle, err
	}
	cycle = cycle.NextStage()
	if err := a.repo.UpdateCycle(ctx, cycle); err != nil {
		return cycle, err
	}
	if err := a.repo.StartStage(ctx, cycle.ID, cycle.StageIndex, a.now().Format(time.RFC3339Nano)); err != nil {
		return cycle, err
	}
	_ = a.audit(ctx, "cycle", cycle.ID, "advanced", fmt.Sprintf("stage=%d", cycle.StageIndex))
	return a.repo.GetCycle(ctx, cycle.ID)
}

func (a *App) HoldCycle(ctx context.Context, cycleID string) (model.Cycle, error) {
	cycle, err := a.repo.GetCycle(ctx, cycleID)
	if err != nil {
		return cycle, err
	}
	if cycle.Status != model.CycleActive {
		return cycle, model.ErrInvalidState
	}
	stage, err := a.repo.GetStage(ctx, cycleID, cycle.StageIndex)
	if err != nil {
		return cycle, err
	}
	if err := policy.HoldStage(&stage); err != nil {
		return cycle, err
	}
	if err := a.repo.UpdateStage(ctx, stage); err != nil {
		return cycle, err
	}
	cycle.Status, cycle.UpdatedAt = model.CyclePaused, a.now()
	if err := a.repo.UpdateCycle(ctx, cycle); err != nil {
		return cycle, err
	}
	_ = a.audit(ctx, "cycle", cycle.ID, "held", stage.Name)
	return a.repo.GetCycle(ctx, cycle.ID)
}
