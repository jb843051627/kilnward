package store

import (
	"context"
	"github.com/jb843051627/kilnward/internal/model"
)

const cycleColumns = `id, load_id, profile, status, stage_index, started_at, ended_at, created_at, updated_at, version`

func (s *Store) CreateCycle(ctx context.Context, cycle model.Cycle) error {
	if err := cycle.Validate(); err != nil {
		return err
	}
	_, err := s.exec(ctx, `INSERT INTO cycles(`+cycleColumns+`) VALUES(?,?,?,?,?,?,?,?,?,?)`, cycle.ID, cycle.LoadID, cycle.Profile, cycle.Status, cycle.StageIndex, nullableTime(cycle.StartedAt), nullableTime(cycle.EndedAt), formatTime(cycle.CreatedAt), formatTime(cycle.UpdatedAt), cycle.Version)
	return err
}

func (s *Store) GetCycle(ctx context.Context, id string) (model.Cycle, error) {
	row := s.db.QueryRowContext(ctx, `SELECT `+cycleColumns+` FROM cycles WHERE id=?`, id)
	item, err := scanCycle(row)
	if err != nil {
		return item, s.notFound(err)
	}
	return item, nil
}

func (s *Store) GetCycleByLoad(ctx context.Context, loadID string) (model.Cycle, error) {
	row := s.db.QueryRowContext(ctx, `SELECT `+cycleColumns+` FROM cycles WHERE load_id=? ORDER BY created_at DESC LIMIT 1`, loadID)
	item, err := scanCycle(row)
	if err != nil {
		return item, s.notFound(err)
	}
	return item, nil
}

func (s *Store) UpdateCycle(ctx context.Context, cycle model.Cycle) error {
	if err := cycle.Validate(); err != nil {
		return err
	}
	result, err := s.exec(ctx, `UPDATE cycles SET status=?, stage_index=?, started_at=?, ended_at=?, updated_at=?, version=? WHERE id=? AND version=?`, cycle.Status, cycle.StageIndex, nullableTime(cycle.StartedAt), nullableTime(cycle.EndedAt), formatTime(cycle.UpdatedAt), cycle.Version+1, cycle.ID, cycle.Version)
	if err != nil {
		return err
	}
	count, _ := result.RowsAffected()
	if count == 0 {
		return model.ErrConflict
	}
	return nil
}

func (s *Store) CountStages(ctx context.Context, cycleID string) (int, error) {
	var count int
	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM stages WHERE cycle_id=?`, cycleID).Scan(&count)
	return count, err
}
