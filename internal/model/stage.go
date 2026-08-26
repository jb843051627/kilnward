package model

import "time"

type StageStatus string

const (
	StageWaiting StageStatus = "waiting"
	StageRunning StageStatus = "running"
	StageHeld    StageStatus = "held"
	StageDone    StageStatus = "done"
	StageFailed  StageStatus = "failed"
)

type Stage struct {
	ID             string      `json:"id"`
	CycleID        string      `json:"cycle_id"`
	Sequence       int         `json:"sequence"`
	Name           string      `json:"name"`
	TargetTempC    float64     `json:"target_temp_c"`
	ToleranceC     float64     `json:"tolerance_c"`
	MinHoldSeconds int         `json:"min_hold_seconds"`
	MaxHoldSeconds int         `json:"max_hold_seconds"`
	Status         StageStatus `json:"status"`
	StartedAt      *time.Time  `json:"started_at,omitempty"`
	EndedAt        *time.Time  `json:"ended_at,omitempty"`
}

func (s Stage) Validate() error {
	if s.CycleID == "" || s.Name == "" || s.Sequence < 0 {
		return ErrValidation
	}
	if s.TargetTempC <= 0 || s.ToleranceC < 0 || s.MinHoldSeconds < 0 || s.MaxHoldSeconds < s.MinHoldSeconds {
		return ErrValidation
	}
	return nil
}

func (s Stage) WithinTemperature(temp float64) bool {
	return temp >= s.TargetTempC-s.ToleranceC && temp <= s.TargetTempC+s.ToleranceC
}

func (s Stage) HoldExpired(now time.Time) bool {
	if s.StartedAt == nil {
		return false
	}
	return now.Sub(*s.StartedAt) > time.Duration(s.MaxHoldSeconds)*time.Second
}
