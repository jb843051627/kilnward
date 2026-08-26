package service

import (
	"context"
	"github.com/jb843051627/kilnward/internal/model"
	"github.com/jb843051627/kilnward/internal/policy"
	"time"
)

func (a *App) PlanMaintenance(ctx context.Context, item model.Maintenance) (model.Maintenance, error) {
	if err := a.requireContext(ctx); err != nil {
		return item, err
	}
	if item.ID == "" {
		item.ID = a.newID("maint")
	}
	if item.OpenedAt.IsZero() {
		item.OpenedAt = a.now()
	}
	if item.Status == "" {
		item.Status = model.MaintenancePlanned
	}
	if err := item.Validate(); err != nil {
		return item, err
	}
	if err := a.repo.CreateMaintenance(ctx, item); err != nil {
		return item, err
	}
	_ = a.audit(ctx, "maintenance", item.ID, "planned", item.Kind)
	return item, nil
}

func (a *App) GetMaintenance(ctx context.Context, id string) (model.Maintenance, error) {
	if err := a.requireContext(ctx); err != nil {
		return model.Maintenance{}, err
	}
	return a.repo.GetMaintenance(ctx, id)
}

func (a *App) ListMaintenance(ctx context.Context, kilnID string, activeOnly bool) ([]model.Maintenance, error) {
	if err := a.requireContext(ctx); err != nil {
		return nil, err
	}
	return a.repo.ListMaintenance(ctx, kilnID, activeOnly)
}

func (a *App) StartMaintenance(ctx context.Context, id string) (model.Maintenance, error) {
	item, err := a.repo.GetMaintenance(ctx, id)
	if err != nil {
		return item, err
	}
	kiln, err := a.repo.GetKiln(ctx, item.KilnID)
	if err != nil {
		return item, err
	}
	if !item.CanStart() || !model.KilnTransition(kiln.State, model.KilnMaintenance) {
		return item, model.ErrInvalidState
	}
	item.Status = model.MaintenanceActive
	if err := a.repo.UpdateMaintenance(ctx, item); err != nil {
		return item, err
	}
	kiln.State, kiln.UpdatedAt = policy.NextState(kiln.State, item.Status), a.now()
	if err := a.repo.UpdateKiln(ctx, kiln); err != nil {
		return item, err
	}
	_ = a.audit(ctx, "maintenance", id, "started", kiln.ID)
	return a.repo.GetMaintenance(ctx, id)
}

func (a *App) CompleteMaintenance(ctx context.Context, id string) (model.Maintenance, error) {
	item, err := a.repo.GetMaintenance(ctx, id)
	if err != nil {
		return item, err
	}
	if err := item.Complete(a.now()); err != nil {
		return item, err
	}
	if err := a.repo.UpdateMaintenance(ctx, item); err != nil {
		return item, err
	}
	kiln, err := a.repo.GetKiln(ctx, item.KilnID)
	if err == nil && kiln.State == model.KilnMaintenance {
		kiln.State, kiln.UpdatedAt = model.KilnReady, a.now()
		_ = a.repo.UpdateKiln(ctx, kiln)
	}
	_ = a.audit(ctx, "maintenance", id, "completed", item.Technician)
	return a.repo.GetMaintenance(ctx, id)
}

func (a *App) MaintenanceBlocks(ctx context.Context, kilnID string) (bool, error) {
	kiln, err := a.repo.GetKiln(ctx, kilnID)
	if err != nil {
		return false, err
	}
	items, err := a.repo.ListMaintenance(ctx, kilnID, true)
	if err != nil {
		return false, err
	}
	return policy.BlocksOperation(kiln, items), nil
}

func (a *App) MaintenanceDue(ctx context.Context, kilnID string) (bool, error) {
	kiln, err := a.repo.GetKiln(ctx, kilnID)
	if err != nil {
		return false, err
	}
	return policy.Due(kiln.LastService, a.now(), 30*24*time.Hour), nil
}
