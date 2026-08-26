package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"

	"github.com/jb843051627/kilnward/internal/model"
)

type Store struct{ db *sql.DB }

func Open(path string) (*Store, error) {
	if path == "" {
		return nil, fmt.Errorf("%w: empty database path", model.ErrValidation)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	s := &Store{db: db}
	if err := s.initSchema(); err != nil {
		_ = db.Close()
		return nil, err
	}
	return s, nil
}

func (s *Store) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}

func (s *Store) Ping(ctx context.Context) error { return s.db.PingContext(ctx) }

func (s *Store) Transaction(ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	if err := fn(tx); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return errors.Join(err, fmt.Errorf("rollback transaction: %w", rollbackErr))
		}
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	return nil
}

func (s *Store) exec(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return s.db.ExecContext(ctx, query, args...)
}

func (s *Store) notFound(err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return model.ErrNotFound
	}
	return err
}

func (s *Store) recordEvent(ctx context.Context, tx *sql.Tx, subjectType, subjectID, action, detail string, at time.Time) error {
	_, err := tx.ExecContext(ctx, `INSERT INTO audits(subject_type, subject_id, action, detail, created_at) VALUES(?,?,?,?,?)`, subjectType, subjectID, action, detail, formatTime(at))
	return err
}

func (s *Store) RecordAudit(ctx context.Context, audit model.Audit) error {
	if err := audit.Validate(); err != nil {
		return err
	}
	_, err := s.exec(ctx, `INSERT INTO audits(subject_type, subject_id, action, detail, created_at) VALUES(?,?,?,?,?)`, audit.SubjectType, audit.SubjectID, audit.Action, audit.Detail, formatTime(audit.CreatedAt))
	return err
}
