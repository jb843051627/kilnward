package store

import (
	"context"
	"github.com/jb843051627/kilnward/internal/model"
)

const commandColumns = `id,kiln_id,load_id,action,status,requested_by,reason,attempts,error_text,created_at,started_at,finished_at`

func (s *Store) CreateCommand(ctx context.Context, command model.OperatorCommand) error {
	if err := command.Validate(); err != nil {
		return err
	}
	_, err := s.exec(ctx, `INSERT INTO commands(`+commandColumns+`) VALUES(?,?,?,?,?,?,?,?,?,?,?,?)`, command.ID, command.KilnID, command.LoadID, command.Action, command.Status, command.RequestedBy, command.Reason, command.Attempts, command.ErrorText, formatTime(command.CreatedAt), nullableTime(command.StartedAt), nullableTime(command.FinishedAt))
	return err
}

func (s *Store) GetCommand(ctx context.Context, id string) (model.OperatorCommand, error) {
	var c model.OperatorCommand
	var started, finished, created string
	err := s.db.QueryRowContext(ctx, `SELECT `+commandColumns+` FROM commands WHERE id=?`, id).Scan(&c.ID, &c.KilnID, &c.LoadID, &c.Action, &c.Status, &c.RequestedBy, &c.Reason, &c.Attempts, &c.ErrorText, &created, &started, &finished)
	if err != nil {
		return c, s.notFound(err)
	}
	c.CreatedAt, c.StartedAt, c.FinishedAt = parseTime(created), pointerTime(started), pointerTime(finished)
	return c, nil
}

func (s *Store) ListCommands(ctx context.Context, kilnID, status string) ([]model.OperatorCommand, error) {
	query, args := `SELECT `+commandColumns+` FROM commands WHERE kiln_id=? ORDER BY created_at DESC`, []any{kilnID}
	if status != "" {
		query = `SELECT ` + commandColumns + ` FROM commands WHERE kiln_id=? AND status=? ORDER BY created_at DESC`
		args = []any{kilnID, status}
	}
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]model.OperatorCommand, 0)
	for rows.Next() {
		var c model.OperatorCommand
		var started, finished, created string
		if err := rows.Scan(&c.ID, &c.KilnID, &c.LoadID, &c.Action, &c.Status, &c.RequestedBy, &c.Reason, &c.Attempts, &c.ErrorText, &created, &started, &finished); err != nil {
			return nil, err
		}
		c.CreatedAt, c.StartedAt, c.FinishedAt = parseTime(created), pointerTime(started), pointerTime(finished)
		out = append(out, c)
	}
	return out, rows.Err()
}

func (s *Store) UpdateCommand(ctx context.Context, command model.OperatorCommand) error {
	_, err := s.exec(ctx, `UPDATE commands SET status=?,attempts=?,error_text=?,started_at=?,finished_at=? WHERE id=?`, command.Status, command.Attempts, command.ErrorText, nullableTime(command.StartedAt), nullableTime(command.FinishedAt), command.ID)
	return err
}
