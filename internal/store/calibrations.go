package store

import (
	"context"
	"github.com/jb843051627/kilnward/internal/model"
)

const calibrationColumns = `id,kiln_id,sensor,reference_c,observed_c,tolerance_c,status,operator,note,created_at,completed_at`

func (s *Store) CreateCalibration(ctx context.Context, item model.Calibration) error {
	if err := item.Validate(); err != nil {
		return err
	}
	_, err := s.exec(ctx, `INSERT INTO calibrations(`+calibrationColumns+`) VALUES(?,?,?,?,?,?,?,?,?,?,?)`, item.ID, item.KilnID, item.Sensor, item.ReferenceC, item.ObservedC, item.ToleranceC, item.Status, item.Operator, item.Note, formatTime(item.CreatedAt), nullableTime(item.CompletedAt))
	return err
}

func (s *Store) GetCalibration(ctx context.Context, id string) (model.Calibration, error) {
	var item model.Calibration
	var status, created, completed string
	err := s.db.QueryRowContext(ctx, `SELECT `+calibrationColumns+` FROM calibrations WHERE id=?`, id).Scan(&item.ID, &item.KilnID, &item.Sensor, &item.ReferenceC, &item.ObservedC, &item.ToleranceC, &status, &item.Operator, &item.Note, &created, &completed)
	if err != nil {
		return item, s.notFound(err)
	}
	item.Status, item.CreatedAt, item.CompletedAt = model.CalibrationStatus(status), parseTime(created), pointerTime(completed)
	return item, nil
}

func (s *Store) ListCalibrations(ctx context.Context, kilnID, sensor string) ([]model.Calibration, error) {
	query, args := `SELECT `+calibrationColumns+` FROM calibrations WHERE kiln_id=? ORDER BY created_at DESC`, []any{kilnID}
	if sensor != "" {
		query = `SELECT ` + calibrationColumns + ` FROM calibrations WHERE kiln_id=? AND sensor=? ORDER BY created_at DESC`
		args = []any{kilnID, sensor}
	}
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]model.Calibration, 0)
	for rows.Next() {
		var item model.Calibration
		var status, created, completed string
		if err := rows.Scan(&item.ID, &item.KilnID, &item.Sensor, &item.ReferenceC, &item.ObservedC, &item.ToleranceC, &status, &item.Operator, &item.Note, &created, &completed); err != nil {
			return nil, err
		}
		item.Status, item.CreatedAt, item.CompletedAt = model.CalibrationStatus(status), parseTime(created), pointerTime(completed)
		out = append(out, item)
	}
	return out, rows.Err()
}

func (s *Store) UpdateCalibration(ctx context.Context, item model.Calibration) error {
	_, err := s.exec(ctx, `UPDATE calibrations SET observed_c=?,tolerance_c=?,status=?,note=?,completed_at=? WHERE id=?`, item.ObservedC, item.ToleranceC, item.Status, item.Note, nullableTime(item.CompletedAt), item.ID)
	return err
}
