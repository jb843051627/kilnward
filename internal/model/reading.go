package model

import "time"

type ReadingQuality string

const (
	QualityGood     ReadingQuality = "good"
	QualitySuspect  ReadingQuality = "suspect"
	QualityRejected ReadingQuality = "rejected"
)

type Reading struct {
	ID          string         `json:"id"`
	KilnID      string         `json:"kiln_id"`
	LoadID      string         `json:"load_id"`
	CycleID     string         `json:"cycle_id"`
	Sensor      string         `json:"sensor"`
	Temperature float64        `json:"temperature"`
	Atmosphere  float64        `json:"atmosphere"`
	Power       float64        `json:"power"`
	RecordedAt  time.Time      `json:"recorded_at"`
	Quality     ReadingQuality `json:"quality"`
	Sequence    int64          `json:"sequence"`
}

func (r Reading) Validate() error {
	if r.ID == "" || r.KilnID == "" || r.LoadID == "" || r.CycleID == "" || r.Sensor == "" {
		return ErrValidation
	}
	if r.Temperature < -50 || r.Temperature > 2500 {
		return ErrValidation
	}
	if r.Atmosphere < 0 || r.Atmosphere > 100 {
		return ErrValidation
	}
	if r.Power < 0 || r.Power > 100 {
		return ErrValidation
	}
	return nil
}

func (r Reading) IsFresh(now time.Time, maxAge time.Duration) bool {
	return !r.RecordedAt.After(now) && now.Sub(r.RecordedAt) <= maxAge
}

func (r Reading) IsSafeFor(stage Stage) bool {
	return r.Quality == QualityGood && stage.WithinTemperature(r.Temperature)
}
