package store

import (
	"context"
	"github.com/jb843051627/kilnward/internal/model"
)

const incidentColumns = `id, kiln_id, load_id, code, severity, status, detail, opened_at, acknowledged_at, closed_at, owner`

func (s *Store) CreateIncident(ctx context.Context, incident model.Incident) error {
	if err := incident.Validate(); err != nil {
		return err
	}
	_, err := s.exec(ctx, `INSERT INTO incidents(`+incidentColumns+`) VALUES(?,?,?,?,?,?,?,?,?,?,?)`, incident.ID, incident.KilnID, incident.LoadID, incident.Code, incident.Severity, incident.Status, incident.Detail, formatTime(incident.OpenedAt), nullableTime(incident.AcknowledgedAt), nullableTime(incident.ClosedAt), incident.Owner)
	return err
}

func (s *Store) GetIncident(ctx context.Context, id string) (model.Incident, error) {
	row := s.db.QueryRowContext(ctx, `SELECT `+incidentColumns+` FROM incidents WHERE id=?`, id)
	item, err := scanIncident(row)
	if err != nil {
		return item, s.notFound(err)
	}
	return item, nil
}

func (s *Store) ListIncidents(ctx context.Context, kilnID, status string) ([]model.Incident, error) {
	query, args := `SELECT `+incidentColumns+` FROM incidents WHERE kiln_id=? ORDER BY opened_at DESC`, []any{kilnID}
	if status != "" {
		query = `SELECT ` + incidentColumns + ` FROM incidents WHERE kiln_id=? AND status=? ORDER BY opened_at DESC`
		args = []any{kilnID, status}
	}
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]model.Incident, 0)
	for rows.Next() {
		item, err := scanIncident(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *Store) UpdateIncident(ctx context.Context, incident model.Incident) error {
	if err := incident.Validate(); err != nil {
		return err
	}
	_, err := s.exec(ctx, `UPDATE incidents SET severity=?, status=?, detail=?, acknowledged_at=?, closed_at=?, owner=? WHERE id=?`, incident.Severity, incident.Status, incident.Detail, nullableTime(incident.AcknowledgedAt), nullableTime(incident.ClosedAt), incident.Owner, incident.ID)
	return err
}
