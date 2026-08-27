package model

import "time"

type StageSnapshot struct {
	Name           string    `json:"name"`
	Status         string    `json:"status"`
	TargetTempC    float64   `json:"target_temp_c"`
	LastTempC      float64   `json:"last_temp_c"`
	ReadingCount   int       `json:"reading_count"`
	LastRecordedAt time.Time `json:"last_recorded_at"`
}

type LoadSummary struct {
	Load       Load            `json:"load"`
	Cycle      Cycle           `json:"cycle"`
	Stages     []StageSnapshot `json:"stages"`
	Incidents  []Incident      `json:"incidents"`
	GateReport GateReport      `json:"gate_report"`
	Healthy    bool            `json:"healthy"`
}

type TimelineEntry struct {
	At      time.Time `json:"at"`
	Subject string    `json:"subject"`
	Action  string    `json:"action"`
	Detail  string    `json:"detail"`
}
