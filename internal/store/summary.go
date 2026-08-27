package store

import (
	"context"
	"github.com/jb843051627/kilnward/internal/model"
)

type SummaryCounts struct{ Kilns, ActiveLoads, OpenIncidents, ActiveMaintenance int }

func (s *Store) SummaryCounts(ctx context.Context) (SummaryCounts, error) {
	var result SummaryCounts
	queries := []struct {
		target *int
		query  string
	}{
		{&result.Kilns, `SELECT COUNT(*) FROM kilns WHERE enabled=1`},
		{&result.ActiveLoads, `SELECT COUNT(*) FROM loads WHERE state IN ('attached','running','paused','cooling')`},
		{&result.OpenIncidents, `SELECT COUNT(*) FROM incidents WHERE status IN ('open','acknowledged')`},
		{&result.ActiveMaintenance, `SELECT COUNT(*) FROM maintenance WHERE status='active'`},
	}
	for _, item := range queries {
		if err := s.db.QueryRowContext(ctx, item.query).Scan(item.target); err != nil {
			return result, err
		}
	}
	return result, nil
}

func (s *Store) LoadForSummary(ctx context.Context, id string) (model.Load, model.Cycle, []model.Stage, error) {
	load, err := s.GetLoad(ctx, id)
	if err != nil {
		return model.Load{}, model.Cycle{}, nil, err
	}
	cycle, err := s.GetCycleByLoad(ctx, id)
	if err != nil {
		return load, model.Cycle{}, nil, err
	}
	stages, err := s.ListStages(ctx, cycle.ID)
	return load, cycle, stages, err
}
