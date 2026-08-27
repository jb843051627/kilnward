package policy

import (
	"github.com/jb843051627/kilnward/internal/model"
	"time"
)

type SamplingPolicy struct {
	MaxAge       time.Duration
	MinPerStage  int
	AllowedDrift float64
}

func DefaultSamplingPolicy() SamplingPolicy {
	return SamplingPolicy{MaxAge: 45 * time.Second, MinPerStage: 3, AllowedDrift: 2.5}
}

func (p SamplingPolicy) ValidateWindow(readings []model.Reading, now time.Time) error {
	if len(readings) < p.MinPerStage {
		return model.ErrMissingReading
	}
	for _, reading := range readings {
		if !reading.IsFresh(now, p.MaxAge) || reading.Quality != model.QualityGood {
			return model.ErrMissingReading
		}
	}
	return nil
}

func (p SamplingPolicy) Stable(readings []model.Reading) bool {
	if len(readings) == 0 {
		return false
	}
	min, max := readings[0].Temperature, readings[0].Temperature
	for _, reading := range readings[1:] {
		if reading.Temperature < min {
			min = reading.Temperature
		}
		if reading.Temperature > max {
			max = reading.Temperature
		}
	}
	return max-min <= p.AllowedDrift
}
