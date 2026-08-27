package store

import (
	"context"
	"database/sql"
	"github.com/jb843051627/kilnward/internal/model"
)

const readingColumns = `id, kiln_id, load_id, cycle_id, sensor, temperature, atmosphere, power, recorded_at, quality, sequence_no`

func (s *Store) AddReading(ctx context.Context, reading model.Reading) error {
	if err := reading.Validate(); err != nil {
		return err
	}
	_, err := s.exec(ctx, `INSERT INTO readings(`+readingColumns+`) VALUES(?,?,?,?,?,?,?,?,?,?,?)`, reading.ID, reading.KilnID, reading.LoadID, reading.CycleID, reading.Sensor, reading.Temperature, reading.Atmosphere, reading.Power, formatTime(reading.RecordedAt), reading.Quality, reading.Sequence)
	return err
}

func (s *Store) AddReadings(ctx context.Context, readings []model.Reading) error {
	return s.Transaction(ctx, func(tx *sql.Tx) error {
		for _, reading := range readings {
			if _, err := tx.ExecContext(ctx, `INSERT INTO readings(`+readingColumns+`) VALUES(?,?,?,?,?,?,?,?,?,?,?)`, reading.ID, reading.KilnID, reading.LoadID, reading.CycleID, reading.Sensor, reading.Temperature, reading.Atmosphere, reading.Power, formatTime(reading.RecordedAt), reading.Quality, reading.Sequence); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Store) ListReadings(ctx context.Context, cycleID string, limit int) ([]model.Reading, error) {
	if limit <= 0 || limit > 500 {
		limit = 100
	}
	rows, err := s.db.QueryContext(ctx, `SELECT `+readingColumns+` FROM readings WHERE cycle_id=? ORDER BY recorded_at DESC LIMIT ?`, cycleID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]model.Reading, 0)
	for rows.Next() {
		item, err := scanReading(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *Store) LatestReadings(ctx context.Context, cycleID string, sensors []string) ([]model.Reading, error) {
	items, err := s.ListReadings(ctx, cycleID, 200)
	if err != nil {
		return nil, err
	}
	if len(sensors) == 0 {
		return items, nil
	}
	allowed := make(map[string]bool, len(sensors))
	for _, sensor := range sensors {
		allowed[sensor] = true
	}
	out := make([]model.Reading, 0, len(sensors))
	seen := make(map[string]bool)
	for _, item := range items {
		if allowed[item.Sensor] && !seen[item.Sensor] {
			out = append(out, item)
			seen[item.Sensor] = true
		}
	}
	return out, nil
}
