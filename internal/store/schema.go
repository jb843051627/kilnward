package store

import "context"

func (s *Store) initSchema() error {
	statements := []string{
		`PRAGMA foreign_keys = ON`,
		`CREATE TABLE IF NOT EXISTS kilns (
			id TEXT PRIMARY KEY, name TEXT NOT NULL, location TEXT NOT NULL,
			state TEXT NOT NULL, max_temp_c REAL NOT NULL, probe_count INTEGER NOT NULL,
			enabled INTEGER NOT NULL, last_service TEXT NOT NULL, created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL, version INTEGER NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS loads (
			id TEXT PRIMARY KEY, kiln_id TEXT NOT NULL, label TEXT NOT NULL,
			state TEXT NOT NULL, material_code TEXT NOT NULL, material_quantity INTEGER NOT NULL,
			material_moisture REAL NOT NULL, profile TEXT NOT NULL, target_temp_c REAL NOT NULL,
			current_stage INTEGER NOT NULL, started_at TEXT, finished_at TEXT, last_error TEXT NOT NULL,
			created_at TEXT NOT NULL, updated_at TEXT NOT NULL, version INTEGER NOT NULL,
			FOREIGN KEY(kiln_id) REFERENCES kilns(id)
		)`,
		`CREATE TABLE IF NOT EXISTS cycles (
			id TEXT PRIMARY KEY, load_id TEXT NOT NULL, profile TEXT NOT NULL, status TEXT NOT NULL,
			stage_index INTEGER NOT NULL, started_at TEXT, ended_at TEXT, created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL, version INTEGER NOT NULL, FOREIGN KEY(load_id) REFERENCES loads(id)
		)`,
		`CREATE TABLE IF NOT EXISTS stages (
			id TEXT PRIMARY KEY, cycle_id TEXT NOT NULL, sequence_no INTEGER NOT NULL,
			name TEXT NOT NULL, target_temp_c REAL NOT NULL, tolerance_c REAL NOT NULL,
			min_hold_seconds INTEGER NOT NULL, max_hold_seconds INTEGER NOT NULL,
			status TEXT NOT NULL, started_at TEXT, ended_at TEXT,
			UNIQUE(cycle_id, sequence_no), FOREIGN KEY(cycle_id) REFERENCES cycles(id)
		)`,
		`CREATE TABLE IF NOT EXISTS readings (
			id TEXT PRIMARY KEY, kiln_id TEXT NOT NULL, load_id TEXT NOT NULL,
			cycle_id TEXT NOT NULL, sensor TEXT NOT NULL, temperature REAL NOT NULL,
			atmosphere REAL NOT NULL, power REAL NOT NULL, recorded_at TEXT NOT NULL,
			quality TEXT NOT NULL, sequence_no INTEGER NOT NULL,
			FOREIGN KEY(kiln_id) REFERENCES kilns(id), FOREIGN KEY(load_id) REFERENCES loads(id),
			FOREIGN KEY(cycle_id) REFERENCES cycles(id)
		)`,
		`CREATE INDEX IF NOT EXISTS readings_cycle_time ON readings(cycle_id, recorded_at)`,
		`CREATE TABLE IF NOT EXISTS gates (
			id TEXT PRIMARY KEY, cycle_id TEXT NOT NULL, stage_seq INTEGER NOT NULL,
			name TEXT NOT NULL, kind TEXT NOT NULL, passed INTEGER NOT NULL, reason TEXT NOT NULL,
			checked_at TEXT NOT NULL, checked_by TEXT NOT NULL, UNIQUE(cycle_id, stage_seq, kind),
			FOREIGN KEY(cycle_id) REFERENCES cycles(id)
		)`,
		`CREATE TABLE IF NOT EXISTS incidents (
			id TEXT PRIMARY KEY, kiln_id TEXT NOT NULL, load_id TEXT NOT NULL,
			code TEXT NOT NULL, severity TEXT NOT NULL, status TEXT NOT NULL, detail TEXT NOT NULL,
			opened_at TEXT NOT NULL, acknowledged_at TEXT, closed_at TEXT, owner TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS maintenance (
			id TEXT PRIMARY KEY, kiln_id TEXT NOT NULL, kind TEXT NOT NULL, status TEXT NOT NULL,
			note TEXT NOT NULL, opened_at TEXT NOT NULL, closed_at TEXT, technician TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS profiles (
			id TEXT PRIMARY KEY, name TEXT NOT NULL, material TEXT NOT NULL, revision INTEGER NOT NULL,
			status TEXT NOT NULL, author TEXT NOT NULL, fingerprint TEXT NOT NULL,
			created_at TEXT NOT NULL, updated_at TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS profile_steps (
			profile_id TEXT NOT NULL, sequence_no INTEGER NOT NULL, name TEXT NOT NULL,
			target_temp_c REAL NOT NULL, ramp_per_minute REAL NOT NULL, tolerance_c REAL NOT NULL,
			min_hold_seconds INTEGER NOT NULL, max_hold_seconds INTEGER NOT NULL, atmosphere REAL NOT NULL,
			PRIMARY KEY(profile_id, sequence_no), FOREIGN KEY(profile_id) REFERENCES profiles(id)
		)`,
		`CREATE TABLE IF NOT EXISTS telemetry_frames (
			id TEXT PRIMARY KEY, kiln_id TEXT NOT NULL, load_id TEXT NOT NULL, cycle_id TEXT NOT NULL,
			gateway TEXT NOT NULL, sequence_no INTEGER NOT NULL, received_at TEXT NOT NULL, checksum INTEGER NOT NULL,
			FOREIGN KEY(kiln_id) REFERENCES kilns(id), FOREIGN KEY(load_id) REFERENCES loads(id), FOREIGN KEY(cycle_id) REFERENCES cycles(id)
		)`,
		`CREATE TABLE IF NOT EXISTS telemetry_samples (
			frame_id TEXT NOT NULL, sensor TEXT NOT NULL, temperature REAL NOT NULL, atmosphere REAL NOT NULL,
			power REAL NOT NULL, recorded_at TEXT NOT NULL, PRIMARY KEY(frame_id, sensor, recorded_at),
			FOREIGN KEY(frame_id) REFERENCES telemetry_frames(id)
		)`,
		`CREATE TABLE IF NOT EXISTS commands (
			id TEXT PRIMARY KEY, kiln_id TEXT NOT NULL, load_id TEXT NOT NULL, action TEXT NOT NULL,
			status TEXT NOT NULL, requested_by TEXT NOT NULL, reason TEXT NOT NULL, attempts INTEGER NOT NULL,
			error_text TEXT NOT NULL, created_at TEXT NOT NULL, started_at TEXT, finished_at TEXT
		)`,
		`CREATE TABLE IF NOT EXISTS calibrations (
			id TEXT PRIMARY KEY, kiln_id TEXT NOT NULL, sensor TEXT NOT NULL, reference_c REAL NOT NULL,
			observed_c REAL NOT NULL, tolerance_c REAL NOT NULL, status TEXT NOT NULL, operator TEXT NOT NULL,
			note TEXT NOT NULL, created_at TEXT NOT NULL, completed_at TEXT
		)`,
		`CREATE TABLE IF NOT EXISTS audits (
			id INTEGER PRIMARY KEY AUTOINCREMENT, subject_type TEXT NOT NULL, subject_id TEXT NOT NULL,
			action TEXT NOT NULL, detail TEXT NOT NULL, created_at TEXT NOT NULL
		)`,
	}
	ctx := context.Background()
	for _, statement := range statements {
		if _, err := s.db.ExecContext(ctx, statement); err != nil {
			return err
		}
	}
	return nil
}
