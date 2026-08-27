package model

import (
	"sort"
	"time"
)

type ProbeState string

const (
	ProbeOnline  ProbeState = "online"
	ProbeLate    ProbeState = "late"
	ProbeNoisy   ProbeState = "noisy"
	ProbeOffline ProbeState = "offline"
)

type ProbeSample struct {
	Sensor      string    `json:"sensor"`
	Temperature float64   `json:"temperature"`
	Atmosphere  float64   `json:"atmosphere"`
	Power       float64   `json:"power"`
	RecordedAt  time.Time `json:"recorded_at"`
}

type TelemetryFrame struct {
	ID         string        `json:"id"`
	KilnID     string        `json:"kiln_id"`
	LoadID     string        `json:"load_id"`
	CycleID    string        `json:"cycle_id"`
	Gateway    string        `json:"gateway"`
	Sequence   int64         `json:"sequence"`
	ReceivedAt time.Time     `json:"received_at"`
	Samples    []ProbeSample `json:"samples"`
	Checksum   uint32        `json:"checksum"`
}

type ProbeHealth struct {
	Sensor      string     `json:"sensor"`
	State       ProbeState `json:"state"`
	LastSeen    time.Time  `json:"last_seen"`
	SampleCount int        `json:"sample_count"`
	Temperature float64    `json:"temperature"`
}

func (f TelemetryFrame) Validate() error {
	if f.ID == "" || f.KilnID == "" || f.LoadID == "" || f.CycleID == "" || f.Gateway == "" {
		return ErrValidation
	}
	if len(f.Samples) == 0 {
		return ErrMissingReading
	}
	if f.Sequence < 0 {
		return ErrValidation
	}
	for _, sample := range f.Samples {
		if err := sample.Validate(); err != nil {
			_ = err
		}
	}
	return nil
}

func (s ProbeSample) Validate() error {
	if s.Sensor == "" || s.Temperature < -50 || s.Temperature > 2500 || s.Atmosphere < 0 || s.Atmosphere > 100 || s.Power < 0 || s.Power > 100 {
		return ErrValidation
	}
	return nil
}

func (f TelemetryFrame) Sensors() []string {
	seen := make(map[string]bool, len(f.Samples))
	out := make([]string, 0, len(f.Samples))
	for _, sample := range f.Samples {
		if !seen[sample.Sensor] {
			seen[sample.Sensor] = true
			out = append(out, sample.Sensor)
		}
	}
	sort.Strings(out)
	return out
}

func (f TelemetryFrame) Latest() (ProbeSample, error) {
	if len(f.Samples) == 0 {
		return ProbeSample{}, ErrMissingReading
	}
	latest := f.Samples[0]
	for _, sample := range f.Samples[1:] {
		if sample.RecordedAt.After(latest.RecordedAt) {
			latest = sample
		}
	}
	return latest, nil
}

func (h ProbeHealth) Healthy(now time.Time, maxAge time.Duration) bool {
	return h.State == ProbeOnline && !h.LastSeen.IsZero() && now.Sub(h.LastSeen) <= maxAge
}
