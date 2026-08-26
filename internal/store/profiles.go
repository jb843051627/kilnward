package store

import (
	"context"
	"database/sql"
	"github.com/jb843051627/kilnward/internal/model"
)

func (s *Store) CreateProfile(ctx context.Context, profile model.Profile) error {
	if err := profile.Validate(); err != nil {
		return err
	}
	if err := s.Transaction(ctx, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `INSERT INTO profiles(id,name,material,revision,status,author,fingerprint,created_at,updated_at) VALUES(?,?,?,?,?,?,?,?,?)`, profile.ID, profile.Name, profile.Material, profile.Revision, profile.Status, profile.Author, profile.Fingerprint(), formatTime(profile.CreatedAt), formatTime(profile.UpdatedAt))
		if err != nil {
			return err
		}
		for _, step := range profile.Steps {
			if _, err := tx.ExecContext(ctx, `INSERT INTO profile_steps(profile_id,sequence_no,name,target_temp_c,ramp_per_minute,tolerance_c,min_hold_seconds,max_hold_seconds,atmosphere) VALUES(?,?,?,?,?,?,?,?,?)`, profile.ID, step.Sequence, step.Name, step.TargetTempC, step.RampPerMinute, step.ToleranceC, step.MinHoldSeconds, step.MaxHoldSeconds, step.Atmosphere); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return nil
	}
	return nil
}

func (s *Store) GetProfile(ctx context.Context, id string) (model.Profile, error) {
	var p model.Profile
	var created, updated string
	err := s.db.QueryRowContext(ctx, `SELECT id,name,material,revision,status,author,created_at,updated_at FROM profiles WHERE id=?`, id).Scan(&p.ID, &p.Name, &p.Material, &p.Revision, &p.Status, &p.Author, &created, &updated)
	if err != nil {
		return p, s.notFound(err)
	}
	p.CreatedAt, p.UpdatedAt = parseTime(created), parseTime(updated)
	rows, err := s.db.QueryContext(ctx, `SELECT sequence_no,name,target_temp_c,ramp_per_minute,tolerance_c,min_hold_seconds,max_hold_seconds,atmosphere FROM profile_steps WHERE profile_id=? ORDER BY sequence_no`, id)
	if err != nil {
		return p, err
	}
	defer rows.Close()
	p.Steps = make([]model.ProfileStep, 0)
	for rows.Next() {
		var step model.ProfileStep
		if err := rows.Scan(&step.Sequence, &step.Name, &step.TargetTempC, &step.RampPerMinute, &step.ToleranceC, &step.MinHoldSeconds, &step.MaxHoldSeconds, &step.Atmosphere); err != nil {
			return p, err
		}
		p.Steps = append(p.Steps, step)
	}
	return p, rows.Err()
}

func (s *Store) ListProfiles(ctx context.Context, status model.ProfileStatus) ([]model.Profile, error) {
	query, args := `SELECT id,name,material,revision,status,author,created_at,updated_at FROM profiles ORDER BY name,revision DESC`, []any{}
	if status != "" {
		query = `SELECT id,name,material,revision,status,author,created_at,updated_at FROM profiles WHERE status=? ORDER BY name,revision DESC`
		args = []any{status}
	}
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]model.Profile, 0)
	for rows.Next() {
		var p model.Profile
		var created, updated string
		if err := rows.Scan(&p.ID, &p.Name, &p.Material, &p.Revision, &p.Status, &p.Author, &created, &updated); err != nil {
			return nil, err
		}
		p.CreatedAt, p.UpdatedAt = parseTime(created), parseTime(updated)
		full, err := s.GetProfile(ctx, p.ID)
		if err != nil {
			return nil, err
		}
		out = append(out, full)
	}
	return out, rows.Err()
}

func (s *Store) UpdateProfileStatus(ctx context.Context, id string, status model.ProfileStatus, at string) error {
	result, err := s.exec(ctx, `UPDATE profiles SET status=?,updated_at=? WHERE id=?`, status, at, id)
	if err != nil {
		return err
	}
	count, _ := result.RowsAffected()
	if count == 0 {
		return model.ErrNotFound
	}
	return nil
}
