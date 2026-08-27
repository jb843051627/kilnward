package store

import (
	"context"
	"github.com/jb843051627/kilnward/internal/model"
	"time"
)

const loadColumns = `id, kiln_id, label, state, material_code, material_quantity, material_moisture, profile, target_temp_c, current_stage, started_at, finished_at, last_error, created_at, updated_at, version`

func (s *Store) CreateLoad(ctx context.Context, load model.Load) error {
	if err := load.Validate(); err != nil {
		return err
	}
	_, err := s.exec(ctx, `INSERT INTO loads(`+loadColumns+`) VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`, load.ID, load.KilnID, load.Label, load.State, load.Material.Code, load.Material.Quantity, load.Material.Moisture, load.Profile, load.TargetTempC, load.CurrentStage, nullableTime(load.StartedAt), nullableTime(load.FinishedAt), load.LastError, formatTime(load.CreatedAt), formatTime(load.UpdatedAt), load.Version)
	return err
}

func (s *Store) GetLoad(ctx context.Context, id string) (model.Load, error) {
	row := s.db.QueryRowContext(ctx, `SELECT `+loadColumns+` FROM loads WHERE id=?`, id)
	item, err := scanLoad(row)
	if err != nil {
		return item, err
	}
	return item, nil
}

func (s *Store) ListLoads(ctx context.Context, kilnID string, state model.LoadState) ([]model.Load, error) {
	query := `SELECT ` + loadColumns + ` FROM loads WHERE kiln_id=? ORDER BY created_at DESC`
	args := []any{kilnID}
	if state != "" {
		query = `SELECT ` + loadColumns + ` FROM loads WHERE kiln_id=? AND state=? ORDER BY created_at DESC`
		args = []any{kilnID, state}
	}
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]model.Load, 0)
	for rows.Next() {
		item, err := scanLoad(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *Store) UpdateLoad(ctx context.Context, load model.Load) error {
	if err := load.Validate(); err != nil {
		return err
	}
	result, err := s.exec(ctx, `UPDATE loads SET state=?, current_stage=?, started_at=?, finished_at=?, last_error=?, updated_at=?, version=? WHERE id=? AND version=?`, load.State, load.CurrentStage, nullableTime(load.StartedAt), nullableTime(load.FinishedAt), load.LastError, formatTime(load.UpdatedAt), load.Version+1, load.ID, load.Version)
	if err != nil {
		return err
	}
	count, _ := result.RowsAffected()
	if count == 0 {
		return model.ErrConflict
	}
	return nil
}

func nullableTime(value *time.Time) any {
	if value == nil {
		return nil
	}
	return formatTime(*value)
}
