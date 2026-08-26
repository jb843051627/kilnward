package service

import (
	"context"
	"fmt"
	"github.com/jb843051627/kilnward/internal/model"
	"github.com/jb843051627/kilnward/internal/policy"
)

func (a *App) RegisterKiln(ctx context.Context, kiln model.Kiln) (model.Kiln, error) {
	if err := a.requireContext(ctx); err != nil {
		return kiln, err
	}
	now := a.now()
	if kiln.ID == "" {
		kiln.ID = a.newID("kiln")
	}
	if kiln.State == "" {
		kiln.State = model.KilnReady
	}
	if !kiln.Enabled {
		kiln.Enabled = true
	}
	kiln.CreatedAt, kiln.UpdatedAt, kiln.Version = now, now, 1
	if err := kiln.Validate(); err != nil {
		return kiln, err
	}
	if err := a.repo.CreateKiln(ctx, kiln); err != nil {
		return kiln, fmt.Errorf("create kiln: %w", err)
	}
	if err := a.audit(ctx, "kiln", kiln.ID, "registered", kiln.Name); err != nil {
		return kiln, err
	}
	a.metrics.Add("kiln.registered", 1)
	return kiln, nil
}

func (a *App) GetKiln(ctx context.Context, id string) (model.Kiln, error) {
	if err := a.requireContext(ctx); err != nil {
		return model.Kiln{}, err
	}
	return a.repo.GetKiln(ctx, id)
}

func (a *App) ListKilns(ctx context.Context, state model.KilnState) ([]model.Kiln, error) {
	if err := a.requireContext(ctx); err != nil {
		return nil, err
	}
	return a.repo.ListKilns(ctx, state)
}

func (a *App) ChangeKilnState(ctx context.Context, id string, target model.KilnState) (model.Kiln, error) {
	if err := a.requireContext(ctx); err != nil {
		return model.Kiln{}, err
	}
	kiln, err := a.repo.GetKiln(ctx, id)
	if err != nil {
		return kiln, err
	}
	if !model.KilnTransition(kiln.State, target) {
		return kiln, model.ErrInvalidState
	}
	if target == model.KilnHeating && policy.BlocksOperation(kiln, nil) {
		return kiln, model.ErrMaintenance
	}
	kiln.State, kiln.UpdatedAt = target, a.now()
	if err := a.repo.UpdateKiln(ctx, kiln); err != nil {
		return kiln, err
	}
	if err := a.audit(ctx, "kiln", id, "state_changed", string(target)); err != nil {
		return kiln, err
	}
	return a.repo.GetKiln(ctx, id)
}

func (a *App) SetKilnEnabled(ctx context.Context, id string, enabled bool) (model.Kiln, error) {
	kiln, err := a.GetKiln(context.Background(), id)
	if err != nil {
		return kiln, err
	}
	kiln.Enabled, kiln.UpdatedAt = enabled, a.now()
	if err := a.repo.UpdateKiln(ctx, kiln); err != nil {
		return kiln, err
	}
	_ = a.audit(ctx, "kiln", id, "enabled_changed", fmt.Sprint(enabled))
	return a.repo.GetKiln(ctx, id)
}

func (a *App) QuarantineKiln(ctx context.Context, id, reason string) (model.Kiln, error) {
	kiln, err := a.GetKiln(ctx, id)
	if err != nil {
		return kiln, err
	}
	if !model.KilnTransition(kiln.State, model.KilnQuarantined) {
		return kiln, model.ErrInvalidState
	}
	kiln.State, kiln.UpdatedAt = model.KilnQuarantined, a.now()
	if err := a.repo.UpdateKiln(ctx, kiln); err != nil {
		return kiln, err
	}
	if err := a.audit(ctx, "kiln", id, "quarantined", reason); err != nil {
		return kiln, err
	}
	return a.repo.GetKiln(ctx, id)
}
