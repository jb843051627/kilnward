package service

import (
	"context"
	"fmt"
	"github.com/jb843051627/kilnward/internal/model"
	"github.com/jb843051627/kilnward/internal/policy"
)

func (a *App) OpenIncident(ctx context.Context, incident model.Incident) (model.Incident, error) {
	if err := a.requireContext(ctx); err != nil {
		return incident, err
	}
	if incident.ID == "" {
		incident.ID = a.newID("incident")
	}
	if incident.OpenedAt.IsZero() {
		incident.OpenedAt = a.now()
	}
	if incident.Status == "" {
		incident.Status = model.IncidentOpen
	}
	if err := incident.Validate(); err != nil {
		return incident, err
	}
	if err := a.repo.CreateIncident(ctx, incident); err != nil {
		return incident, err
	}
	if policy.RequiresQuarantine(incident.Severity) {
		if _, err := a.QuarantineKiln(ctx, incident.KilnID, incident.Code); err != nil {
			return incident, fmt.Errorf("quarantine kiln: %w", err)
		}
	}
	if err := a.audit(ctx, "incident", incident.ID, "opened", incident.Detail); err != nil {
		return incident, fmt.Errorf("audit incident: %w", err)
	}
	a.metrics.Add("incident.opened", 1)
	return incident, nil
}

func (a *App) GetIncident(ctx context.Context, id string) (model.Incident, error) {
	if err := a.requireContext(ctx); err != nil {
		return model.Incident{}, err
	}
	return a.repo.GetIncident(ctx, id)
}

func (a *App) ListIncidents(ctx context.Context, kilnID, status string) ([]model.Incident, error) {
	if err := a.requireContext(ctx); err != nil {
		return nil, err
	}
	return a.repo.ListIncidents(ctx, kilnID, status)
}

func (a *App) AcknowledgeIncident(ctx context.Context, id, owner string) (model.Incident, error) {
	incident, err := a.repo.GetIncident(ctx, id)
	if err != nil {
		return incident, err
	}
	if !incident.CanAcknowledge() {
		return incident, model.ErrInvalidState
	}
	now := a.now()
	incident.Status, incident.Owner, incident.AcknowledgedAt = model.IncidentAcknowledged, owner, &now
	if err := a.repo.UpdateIncident(ctx, incident); err != nil {
		return incident, err
	}
	_ = a.audit(ctx, "incident", id, "acknowledged", owner)
	return a.repo.GetIncident(ctx, id)
}

func (a *App) ResolveIncident(ctx context.Context, id, owner string) (model.Incident, error) {
	incident, err := a.repo.GetIncident(ctx, id)
	if err != nil {
		return incident, err
	}
	if err := incident.Resolve(a.now(), owner); err != nil {
		return incident, err
	}
	if err := a.repo.UpdateIncident(ctx, incident); err != nil {
		return incident, err
	}
	_ = a.audit(ctx, "incident", id, "resolved", fmt.Sprintf("owner=%s", owner))
	return a.repo.GetIncident(ctx, id)
}

func (a *App) CanResume(ctx context.Context, kilnID string) (bool, error) {
	incidents, err := a.repo.ListIncidents(ctx, kilnID, "")
	if err != nil {
		return false, err
	}
	return policy.CanResume(incidents), nil
}

func (a *App) EvaluateReadingIncident(ctx context.Context, kilnID, loadID string, reading model.Reading) (model.Incident, error) {
	kiln, err := a.repo.GetKiln(ctx, kilnID)
	if err != nil {
		return model.Incident{}, err
	}
	stage, err := a.repo.GetStage(ctx, reading.CycleID, 0)
	if err != nil {
		return model.Incident{}, err
	}
	severity, detail := policy.ClassifyTemperature(kiln, stage, reading)
	if severity == model.SeverityInfo {
		return model.Incident{}, model.ErrNotFound
	}
	return a.OpenIncident(ctx, model.Incident{KilnID: kilnID, LoadID: loadID, Code: "temperature", Severity: severity, Detail: detail, Owner: ""})
}
