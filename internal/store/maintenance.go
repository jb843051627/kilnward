package store

import (
	"context"
	"github.com/jb843051627/kilnward/internal/model"
)

const maintenanceColumns = `id, kiln_id, kind, status, note, opened_at, closed_at, technician`

func (s *Store) CreateMaintenance(ctx context.Context, item model.Maintenance) error {
	if err := item.Validate(); err != nil {
		return err
	}
	_, err := s.exec(ctx, `INSERT INTO maintenance(`+maintenanceColumns+`) VALUES(?,?,?,?,?,?,?,?)`, item.ID, item.KilnID, item.Kind, item.Status, item.Note, formatTime(item.OpenedAt), nullableTime(item.ClosedAt), item.Technician)
	return err
}

func (s *Store) GetMaintenance(ctx context.Context, id string) (model.Maintenance, error) {
	row := s.db.QueryRowContext(ctx, `SELECT `+maintenanceColumns+` FROM maintenance WHERE id=?`, id)
	item, err := scanMaintenance(row)
	if err != nil {
		return item, s.notFound(err)
	}
	return item, nil
}

func (s *Store) ListMaintenance(ctx context.Context, kilnID string, activeOnly bool) ([]model.Maintenance, error) {
	query, args := `SELECT `+maintenanceColumns+` FROM maintenance WHERE kiln_id=? ORDER BY opened_at DESC`, []any{kilnID}
	if activeOnly {
		query = `SELECT ` + maintenanceColumns + ` FROM maintenance WHERE kiln_id=? AND status IN (?,?) ORDER BY opened_at DESC`
		args = []any{kilnID, model.MaintenancePlanned, model.MaintenanceActive}
	}
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]model.Maintenance, 0)
	for rows.Next() {
		item, err := scanMaintenance(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *Store) UpdateMaintenance(ctx context.Context, item model.Maintenance) error {
	if err := item.Validate(); err != nil {
		return err
	}
	_, err := s.exec(ctx, `UPDATE maintenance SET status=?, note=?, closed_at=?, technician=? WHERE id=?`, item.Status, item.Note, nullableTime(item.ClosedAt), item.Technician, item.ID)
	return err
}
