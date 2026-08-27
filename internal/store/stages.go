package store

import (
	"context"
	"github.com/jb843051627/kilnward/internal/model"
)

const stageColumns = `id, cycle_id, sequence_no, name, target_temp_c, tolerance_c, min_hold_seconds, max_hold_seconds, status, started_at, ended_at`

func (s *Store) CreateStage(ctx context.Context, stage model.Stage) error {
	if err := stage.Validate(); err != nil {
		return err
	}
	_, err := s.exec(ctx, `INSERT INTO stages(`+stageColumns+`) VALUES(?,?,?,?,?,?,?,?,?,?,?)`, stage.ID, stage.CycleID, stage.Sequence, stage.Name, stage.TargetTempC, stage.ToleranceC, stage.MinHoldSeconds, stage.MaxHoldSeconds, stage.Status, nullableTime(stage.StartedAt), nullableTime(stage.EndedAt))
	return err
}

func (s *Store) ListStages(ctx context.Context, cycleID string) ([]model.Stage, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT `+stageColumns+` FROM stages WHERE cycle_id=? ORDER BY sequence_no`, cycleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]model.Stage, 0)
	for rows.Next() {
		item, err := scanStage(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *Store) GetStage(ctx context.Context, cycleID string, sequence int) (model.Stage, error) {
	row := s.db.QueryRowContext(ctx, `SELECT `+stageColumns+` FROM stages WHERE cycle_id=? AND sequence_no=?`, cycleID, sequence)
	item, err := scanStage(row)
	if err != nil {
		return item, s.notFound(err)
	}
	return item, nil
}

func (s *Store) UpdateStage(ctx context.Context, stage model.Stage) error {
	if err := stage.Validate(); err != nil {
		return err
	}
	_, err := s.exec(ctx, `UPDATE stages SET status=?, started_at=?, ended_at=? WHERE id=?`, stage.Status, nullableTime(stage.StartedAt), nullableTime(stage.EndedAt), stage.ID)
	return err
}

func (s *Store) StartStage(ctx context.Context, cycleID string, sequence int, at string) error {
	result, err := s.exec(ctx, `UPDATE stages SET status=?, started_at=? WHERE cycle_id=? AND sequence_no=? AND status=?`, model.StageRunning, at, cycleID, sequence, model.StageWaiting)
	if err != nil {
		return err
	}
	count, _ := result.RowsAffected()
	if count == 0 {
		return model.ErrConflict
	}
	return nil
}
