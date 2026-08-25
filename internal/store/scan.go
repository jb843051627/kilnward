package store

import (
	"database/sql"
	"github.com/jb843051627/kilnward/internal/model"
)

type scanner interface{ Scan(...any) error }

func scanKiln(row scanner) (model.Kiln, error) {
	var k model.Kiln
	var state, last, created, updated string
	var enabled int
	err := row.Scan(&k.ID, &k.Name, &k.Location, &state, &k.MaxTempC, &k.ProbeCount, &enabled, &last, &created, &updated, &k.Version)
	k.State, k.Enabled, k.LastService, k.CreatedAt, k.UpdatedAt = model.KilnState(state), enabled == 1, parseTime(last), parseTime(created), parseTime(updated)
	return k, err
}

func scanLoad(row scanner) (model.Load, error) {
	var l model.Load
	var state, started, finished, created, updated string
	err := row.Scan(&l.ID, &l.KilnID, &l.Label, &state, &l.Material.Code, &l.Material.Quantity, &l.Material.Moisture, &l.Profile, &l.TargetTempC, &l.CurrentStage, &started, &finished, &l.LastError, &created, &updated, &l.Version)
	l.State, l.StartedAt, l.FinishedAt, l.CreatedAt, l.UpdatedAt = model.LoadState(state), pointerTime(started), pointerTime(finished), parseTime(created), parseTime(updated)
	return l, err
}

func scanCycle(row scanner) (model.Cycle, error) {
	var c model.Cycle
	var status, started, ended, created, updated string
	err := row.Scan(&c.ID, &c.LoadID, &c.Profile, &status, &c.StageIndex, &started, &ended, &created, &updated, &c.Version)
	c.Status, c.StartedAt, c.EndedAt, c.CreatedAt, c.UpdatedAt = model.CycleStatus(status), pointerTime(started), pointerTime(ended), parseTime(created), parseTime(updated)
	return c, err
}

func scanStage(row scanner) (model.Stage, error) {
	var s model.Stage
	var status, started, ended string
	err := row.Scan(&s.ID, &s.CycleID, &s.Sequence, &s.Name, &s.TargetTempC, &s.ToleranceC, &s.MinHoldSeconds, &s.MaxHoldSeconds, &status, &started, &ended)
	s.Status, s.StartedAt, s.EndedAt = model.StageStatus(status), pointerTime(started), pointerTime(ended)
	return s, err
}

func scanReading(row scanner) (model.Reading, error) {
	var r model.Reading
	var quality, recorded string
	err := row.Scan(&r.ID, &r.KilnID, &r.LoadID, &r.CycleID, &r.Sensor, &r.Temperature, &r.Atmosphere, &r.Power, &recorded, &quality, &r.Sequence)
	r.RecordedAt, r.Quality = parseTime(recorded), model.ReadingQuality(quality)
	return r, err
}

func scanGate(row scanner) (model.Gate, error) {
	var g model.Gate
	var passed int
	var checked string
	err := row.Scan(&g.ID, &g.CycleID, &g.StageSeq, &g.Name, &g.Kind, &passed, &g.Reason, &checked, &g.CheckedBy)
	g.Passed, g.CheckedAt = passed == 1, parseTime(checked)
	return g, err
}

func scanIncident(row scanner) (model.Incident, error) {
	var i model.Incident
	var severity, status, opened, acknowledged, closed string
	err := row.Scan(&i.ID, &i.KilnID, &i.LoadID, &i.Code, &severity, &status, &i.Detail, &opened, &acknowledged, &closed, &i.Owner)
	i.Severity, i.Status, i.OpenedAt = model.IncidentSeverity(severity), model.IncidentStatus(status), parseTime(opened)
	i.AcknowledgedAt, i.ClosedAt = pointerTime(acknowledged), pointerTime(closed)
	return i, err
}

func scanMaintenance(row scanner) (model.Maintenance, error) {
	var m model.Maintenance
	var status, opened, closed string
	err := row.Scan(&m.ID, &m.KilnID, &m.Kind, &status, &m.Note, &opened, &closed, &m.Technician)
	m.Status, m.OpenedAt, m.ClosedAt = model.MaintenanceStatus(status), parseTime(opened), pointerTime(closed)
	return m, err
}

func scanAudit(row scanner) (model.Audit, error) {
	var a model.Audit
	var created string
	err := row.Scan(&a.ID, &a.SubjectType, &a.SubjectID, &a.Action, &a.Detail, &created)
	a.CreatedAt = parseTime(created)
	return a, err
}

var _ = sql.ErrNoRows
