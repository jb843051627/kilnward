package service

import (
	"context"
	"fmt"
	"github.com/jb843051627/kilnward/internal/model"
	"github.com/jb843051627/kilnward/internal/policy"
)

func (a *App) IssueCommand(ctx context.Context, command model.OperatorCommand) (model.OperatorCommand, error) {
	if err := a.requireContext(ctx); err != nil {
		return command, err
	}
	if command.ID == "" {
		command.ID = a.newID("command")
	}
	if command.Status == "" {
		command.Status = model.CommandQueued
	}
	if command.CreatedAt.IsZero() {
		command.CreatedAt = a.now()
	}
	kiln, err := a.repo.GetKiln(ctx, command.KilnID)
	if err != nil {
		return command, err
	}
	var load *model.Load
	if command.LoadID != "" {
		item, loadErr := a.repo.GetLoad(ctx, command.LoadID)
		if loadErr != nil {
			return command, loadErr
		}
		load = &item
	}
	if err := policy.CommandAllowed(command, kiln, load); err != nil {
		return command, err
	}
	if err := a.repo.CreateCommand(ctx, command); err != nil {
		return command, err
	}
	_ = a.audit(ctx, "command", command.ID, "queued", string(command.Action))
	return command, nil
}

func (a *App) GetCommand(ctx context.Context, id string) (model.OperatorCommand, error) {
	if err := a.requireContext(ctx); err != nil {
		return model.OperatorCommand{}, err
	}
	return a.repo.GetCommand(ctx, id)
}
func (a *App) ListCommands(ctx context.Context, kilnID, status string) ([]model.OperatorCommand, error) {
	if err := a.requireContext(ctx); err != nil {
		return nil, err
	}
	return a.repo.ListCommands(ctx, kilnID, status)
}

func (a *App) RunCommand(ctx context.Context, id string) (model.OperatorCommand, error) {
	command, err := a.repo.GetCommand(ctx, id)
	if err != nil {
		return command, err
	}
	if err := command.Begin(a.now()); err != nil {
		return command, err
	}
	if err := a.repo.UpdateCommand(ctx, command); err != nil {
		return command, err
	}
	if command.LoadID != "" {
		load, loadErr := a.repo.GetLoad(ctx, command.LoadID)
		if loadErr != nil {
			command.Fail(a.now(), loadErr)
			_ = a.repo.UpdateCommand(ctx, command)
			return command, loadErr
		}
		target, stateErr := policy.RecoveryState(load, command.Action)
		if stateErr != nil {
			command.Fail(a.now(), stateErr)
			_ = a.repo.UpdateCommand(ctx, command)
			return command, stateErr
		}
		load.State, load.UpdatedAt = target, a.now()
		if updateErr := a.repo.UpdateLoad(ctx, load); updateErr != nil {
			command.Fail(a.now(), updateErr)
			_ = a.repo.UpdateCommand(ctx, command)
			return command, updateErr
		}
	}
	if command.Action == model.CommandQuarantine {
		if _, err := a.QuarantineKiln(ctx, command.KilnID, command.Reason); err != nil {
			command.Fail(a.now(), err)
			_ = a.repo.UpdateCommand(ctx, command)
			return command, err
		}
	}
	if err := command.Succeed(a.now()); err != nil {
		return command, err
	}
	if err := a.repo.UpdateCommand(ctx, command); err != nil {
		return command, err
	}
	_ = a.audit(ctx, "command", command.ID, "succeeded", fmt.Sprintf("attempt=%d", command.Attempts))
	return a.repo.GetCommand(ctx, id)
}

func (a *App) CancelCommand(ctx context.Context, id string) (model.OperatorCommand, error) {
	command, err := a.repo.GetCommand(ctx, id)
	if err != nil {
		return command, err
	}
	if err := command.Cancel(a.now()); err != nil {
		return command, err
	}
	if err := a.repo.UpdateCommand(ctx, command); err != nil {
		return command, err
	}
	_ = a.audit(ctx, "command", id, "canceled", command.Reason)
	return a.repo.GetCommand(ctx, id)
}

func (a *App) RetryCommand(ctx context.Context, id string) (model.OperatorCommand, error) {
	command, err := a.repo.GetCommand(ctx, id)
	if err != nil {
		return command, err
	}
	if !policy.Retryable(fmt.Errorf(command.ErrorText), command.Attempts) {
		return command, model.ErrInvalidState
	}
	command.Status, command.FinishedAt = model.CommandQueued, nil
	if err := a.repo.UpdateCommand(ctx, command); err != nil {
		return command, err
	}
	return a.repo.GetCommand(ctx, id)
}
