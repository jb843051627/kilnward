package service

import (
	"context"
	"fmt"
	"github.com/jb843051627/kilnward/internal/model"
	"github.com/jb843051627/kilnward/internal/policy"
)

func (a *App) CreateLoad(ctx context.Context, load model.Load) (model.Load, error) {
	if err := a.requireContext(ctx); err != nil {
		return load, err
	}
	kiln, err := a.repo.GetKiln(ctx, load.KilnID)
	if err != nil {
		return load, err
	}
	if !kiln.CanAcceptLoad() {
		return load, model.ErrInvalidState
	}
	now := a.now()
	if load.ID == "" {
		load.ID = a.newID("load")
	}
	if load.State == "" {
		load.State = model.LoadDraft
	}
	load.CreatedAt, load.UpdatedAt, load.Version = now, now, 1
	if err := load.Validate(); err != nil {
		return load, err
	}
	if err := a.repo.CreateLoad(ctx, load); err != nil {
		return load, fmt.Errorf("create load: %w", err)
	}
	cycle, err := a.planDefaultCycle(ctx, load)
	if err != nil {
		return load, err
	}
	if err := a.audit(ctx, "load", load.ID, "created", cycle.ID); err != nil {
		return load, err
	}
	a.metrics.Add("load.created", 1)
	return load, nil
}

func (a *App) planDefaultCycle(ctx context.Context, load model.Load) (model.Cycle, error) {
	now := a.now()
	cycle := model.Cycle{ID: a.newID("cycle"), LoadID: load.ID, Profile: load.Profile, Status: model.CyclePlanned, StageIndex: 0, CreatedAt: now, UpdatedAt: now, Version: 1}
	if err := a.repo.CreateCycle(ctx, cycle); err != nil {
		return cycle, err
	}
	for i, stage := range defaultStages(cycle.ID, load.TargetTempC) {
		if err := a.repo.CreateStage(ctx, stageForID(stage, a.newID(fmt.Sprintf("stage-%d", i)))); err != nil {
			return cycle, err
		}
	}
	return cycle, nil
}

func defaultStages(cycleID string, target float64) []model.Stage {
	return []model.Stage{
		{CycleID: cycleID, Sequence: 0, Name: "预热", TargetTempC: target * 0.35, ToleranceC: 18, MinHoldSeconds: 1, MaxHoldSeconds: 600, Status: model.StageWaiting},
		{CycleID: cycleID, Sequence: 1, Name: "升温", TargetTempC: target * 0.75, ToleranceC: 12, MinHoldSeconds: 1, MaxHoldSeconds: 900, Status: model.StageWaiting},
		{CycleID: cycleID, Sequence: 2, Name: "保温", TargetTempC: target, ToleranceC: 8, MinHoldSeconds: 1, MaxHoldSeconds: 1200, Status: model.StageWaiting},
		{CycleID: cycleID, Sequence: 3, Name: "冷却", TargetTempC: target * 0.25, ToleranceC: 25, MinHoldSeconds: 1, MaxHoldSeconds: 1500, Status: model.StageWaiting},
	}
}

func stageForID(stage model.Stage, id string) model.Stage { stage.ID = id; return stage }

func (a *App) GetLoad(ctx context.Context, id string) (model.Load, error) {
	if err := a.requireContext(ctx); err != nil {
		return model.Load{}, err
	}
	return a.repo.GetLoad(ctx, id)
}

func (a *App) ListLoads(ctx context.Context, kilnID string, state model.LoadState) ([]model.Load, error) {
	if err := a.requireContext(ctx); err != nil {
		return nil, err
	}
	return a.repo.ListLoads(ctx, kilnID, state)
}

func (a *App) AttachLoad(ctx context.Context, id string) (model.Load, error) {
	load, err := a.repo.GetLoad(ctx, id)
	if err != nil {
		return load, err
	}
	kiln, err := a.repo.GetKiln(ctx, load.KilnID)
	if err != nil {
		return load, err
	}
	if err := policy.AttachAllowed(kiln, load); err != nil {
		return load, err
	}
	load.State, load.UpdatedAt = model.LoadAttached, a.now()
	if err := a.repo.UpdateLoad(ctx, load); err != nil {
		return load, err
	}
	if err := a.audit(ctx, "load", id, "attached", kiln.ID); err != nil {
		return load, err
	}
	return a.repo.GetLoad(ctx, id)
}

func (a *App) RejectLoad(ctx context.Context, id, reason string) (model.Load, error) {
	load, err := a.repo.GetLoad(ctx, id)
	if err != nil {
		return load, err
	}
	if !load.CanTransition(model.LoadRejected) {
		return load, model.ErrInvalidState
	}
	load.State, load.LastError, load.UpdatedAt = model.LoadRejected, reason, a.now()
	if err := a.repo.UpdateLoad(ctx, load); err != nil {
		return load, err
	}
	_ = a.audit(ctx, "load", id, "rejected", reason)
	return a.repo.GetLoad(ctx, id)
}
