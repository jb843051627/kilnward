package store

import (
	"context"
	"github.com/jb843051627/kilnward/internal/model"
)

const gateColumns = `id, cycle_id, stage_seq, name, kind, passed, reason, checked_at, checked_by`

func (s *Store) SaveGate(ctx context.Context, gate model.Gate) error {
	if err := gate.Validate(); err != nil {
		return err
	}
	_, err := s.exec(ctx, `INSERT INTO gates(`+gateColumns+`) VALUES(?,?,?,?,?,?,?,?,?) ON CONFLICT(cycle_id,stage_seq,kind) DO UPDATE SET name=excluded.name, passed=excluded.passed, reason=excluded.reason, checked_at=excluded.checked_at, checked_by=excluded.checked_by`, gate.ID, gate.CycleID, gate.StageSeq, gate.Name, gate.Kind, boolInt(gate.Passed), gate.Reason, formatTime(gate.CheckedAt), gate.CheckedBy)
	return err
}

func (s *Store) ListGates(ctx context.Context, cycleID string, stage int) ([]model.Gate, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT `+gateColumns+` FROM gates WHERE cycle_id=? AND stage_seq=? ORDER BY kind`, cycleID, stage)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]model.Gate, 0)
	for rows.Next() {
		item, err := scanGate(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *Store) GateReport(ctx context.Context, cycleID string, stage int) (model.GateReport, error) {
	checks, err := s.ListGates(ctx, cycleID, stage)
	if err != nil {
		return model.GateReport{}, err
	}
	return model.GateReport{CycleID: cycleID, Stage: stage, Checks: checks, Passed: len(checks) > 0 && allGatesPassed(checks)}, nil
}

func allGatesPassed(checks []model.Gate) bool {
	for _, check := range checks {
		if !check.Passed {
			return false
		}
	}
	return true
}
