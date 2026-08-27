package store

import (
	"context"
	"github.com/jb843051627/kilnward/internal/model"
)

func (s *Store) ListAudits(ctx context.Context, subjectType, subjectID string, limit int) ([]model.Audit, error) {
	if limit <= 0 || limit > 500 {
		limit = 100
	}
	query := `SELECT id, subject_type, subject_id, action, detail, created_at FROM audits ORDER BY created_at DESC LIMIT ?`
	args := []any{limit}
	if subjectType != "" && subjectID != "" {
		query = `SELECT id, subject_type, subject_id, action, detail, created_at FROM audits WHERE subject_type=? AND subject_id=? ORDER BY created_at DESC LIMIT ?`
		args = []any{subjectType, subjectID, limit}
	}
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]model.Audit, 0)
	for rows.Next() {
		item, err := scanAudit(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *Store) CountOpenIncidents(ctx context.Context, kilnID string) (int, error) {
	var count int
	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM incidents WHERE kiln_id=? AND status IN (?,?)`, kilnID, model.IncidentOpen, model.IncidentAcknowledged).Scan(&count)
	return count, err
}

func (s *Store) CountActiveLoads(ctx context.Context, kilnID string) (int, error) {
	var count int
	err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM loads WHERE kiln_id=? AND state IN (?,?,?)`, kilnID, model.LoadAttached, model.LoadRunning, model.LoadCooling).Scan(&count)
	return count, err
}
