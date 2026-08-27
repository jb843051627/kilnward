package store

import (
	"context"
	"database/sql"
	"github.com/jb843051627/kilnward/internal/model"
)

const kilnColumns = `id, name, location, state, max_temp_c, probe_count, enabled, last_service, created_at, updated_at, version`

func (s *Store) CreateKiln(ctx context.Context, kiln model.Kiln) error {
	if err := kiln.Validate(); err != nil {
		return err
	}
	_, err := s.exec(ctx, `INSERT INTO kilns(`+kilnColumns+`) VALUES(?,?,?,?,?,?,?,?,?,?,?)`, kiln.ID, kiln.Name, kiln.Location, kiln.State, kiln.MaxTempC, kiln.ProbeCount, boolInt(kiln.Enabled), formatTime(kiln.LastService), formatTime(kiln.CreatedAt), formatTime(kiln.UpdatedAt), kiln.Version)
	return err
}

func (s *Store) GetKiln(ctx context.Context, id string) (model.Kiln, error) {
	row := s.db.QueryRowContext(ctx, `SELECT `+kilnColumns+` FROM kilns WHERE id=?`, id)
	kiln, err := scanKiln(row)
	if err != nil {
		return kiln, err
	}
	return kiln, nil
}

func (s *Store) ListKilns(ctx context.Context, state model.KilnState) ([]model.Kiln, error) {
	query, args := `SELECT `+kilnColumns+` FROM kilns ORDER BY name`, []any{}
	if state != "" {
		query = `SELECT ` + kilnColumns + ` FROM kilns WHERE state=? ORDER BY name`
		args = []any{state}
	}
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]model.Kiln, 0)
	for rows.Next() {
		item, err := scanKiln(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *Store) UpdateKiln(ctx context.Context, kiln model.Kiln) error {
	if err := kiln.Validate(); err != nil {
		return err
	}
	result, err := s.exec(ctx, `UPDATE kilns SET name=?, location=?, state=?, max_temp_c=?, probe_count=?, enabled=?, last_service=?, updated_at=?, version=? WHERE id=? AND version=?`, kiln.Name, kiln.Location, kiln.State, kiln.MaxTempC, kiln.ProbeCount, boolInt(kiln.Enabled), formatTime(kiln.LastService), formatTime(kiln.UpdatedAt), kiln.Version+1, kiln.ID, kiln.Version)
	if err != nil {
		return err
	}
	count, _ := result.RowsAffected()
	if count == 0 {
		return model.ErrConflict
	}
	return nil
}

func (s *Store) SetKilnState(ctx context.Context, id string, from, to model.KilnState, at string) error {
	result, err := s.exec(ctx, `UPDATE kilns SET state=?, updated_at=?, version=version+1 WHERE id=? AND state=?`, to, at, id, from)
	if err != nil {
		return err
	}
	count, _ := result.RowsAffected()
	if count == 0 {
		return model.ErrConflict
	}
	return nil
}

func boolInt(value bool) int {
	if value {
		return 1
	}
	return 0
}

var _ = sql.ErrNoRows
