package store

import (
	"context"
	"database/sql"
	"github.com/jb843051627/kilnward/internal/model"
)

func (s *Store) SaveTelemetry(ctx context.Context, frame model.TelemetryFrame) error {
	if err := frame.Validate(); err != nil {
		return err
	}
	return s.Transaction(ctx, func(tx *sql.Tx) error {
		if _, err := tx.ExecContext(ctx, `INSERT INTO telemetry_frames(id,kiln_id,load_id,cycle_id,gateway,sequence_no,received_at,checksum) VALUES(?,?,?,?,?,?,?,?)`, frame.ID, frame.KilnID, frame.LoadID, frame.CycleID, frame.Gateway, frame.Sequence, formatTime(frame.ReceivedAt), frame.Checksum); err != nil {
			return err
		}
		for _, sample := range frame.Samples {
			if _, err := tx.ExecContext(ctx, `INSERT INTO telemetry_samples(frame_id,sensor,temperature,atmosphere,power,recorded_at) VALUES(?,?,?,?,?,?)`, frame.ID, sample.Sensor, sample.Temperature, sample.Atmosphere, sample.Power, formatTime(sample.RecordedAt)); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Store) ListTelemetry(ctx context.Context, cycleID string, limit int) ([]model.TelemetryFrame, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	rows, err := s.db.QueryContext(ctx, `SELECT id,kiln_id,load_id,cycle_id,gateway,sequence_no,received_at,checksum FROM telemetry_frames WHERE cycle_id=? ORDER BY received_at DESC LIMIT ?`, cycleID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	frames := make([]model.TelemetryFrame, 0)
	for rows.Next() {
		var frame model.TelemetryFrame
		var received string
		if err := rows.Scan(&frame.ID, &frame.KilnID, &frame.LoadID, &frame.CycleID, &frame.Gateway, &frame.Sequence, &received, &frame.Checksum); err != nil {
			return nil, err
		}
		frame.ReceivedAt = parseTime(received)
		samples, err := s.samplesForFrame(ctx, frame.ID)
		if err != nil {
			return nil, err
		}
		frame.Samples = samples
		frames = append(frames, frame)
	}
	return frames, rows.Err()
}

func (s *Store) samplesForFrame(ctx context.Context, frameID string) ([]model.ProbeSample, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT sensor,temperature,atmosphere,power,recorded_at FROM telemetry_samples WHERE frame_id=? ORDER BY recorded_at`, frameID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]model.ProbeSample, 0)
	for rows.Next() {
		var sample model.ProbeSample
		var at string
		if err := rows.Scan(&sample.Sensor, &sample.Temperature, &sample.Atmosphere, &sample.Power, &at); err != nil {
			return nil, err
		}
		sample.RecordedAt = parseTime(at)
		out = append(out, sample)
	}
	return out, rows.Err()
}

func (s *Store) LastTelemetry(ctx context.Context, cycleID string) (model.TelemetryFrame, error) {
	frames, err := s.ListTelemetry(ctx, cycleID, 1)
	if err != nil {
		return model.TelemetryFrame{}, err
	}
	if len(frames) == 0 {
		return model.TelemetryFrame{}, model.ErrNotFound
	}
	return frames[0], nil
}
