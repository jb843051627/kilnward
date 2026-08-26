package service

import (
	"context"
	"github.com/jb843051627/kilnward/internal/model"
	"github.com/jb843051627/kilnward/internal/policy"
)

func (a *App) Timeline(ctx context.Context, subjectType, subjectID string, limit int) ([]model.TimelineEntry, error) {
	audits, err := a.repo.ListAudits(ctx, subjectType, subjectID, limit)
	if err != nil {
		return nil, err
	}
	entries := make([]model.TimelineEntry, 0, len(audits))
	for _, audit := range audits {
		entries = append(entries, model.TimelineEntry{At: audit.CreatedAt, Subject: audit.SubjectType + "/" + audit.SubjectID, Action: audit.Action, Detail: audit.Detail})
	}
	return policy.Limit(policy.Sort(entries), limit), nil
}

func (a *App) AuditForLoad(ctx context.Context, loadID string, limit int) ([]model.TimelineEntry, error) {
	return a.Timeline(ctx, "load", loadID, limit)
}
func (a *App) AuditForKiln(ctx context.Context, kilnID string, limit int) ([]model.TimelineEntry, error) {
	return a.Timeline(ctx, "kiln", kilnID, limit)
}
