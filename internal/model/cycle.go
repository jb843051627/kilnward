package model

import (
	"fmt"
	"time"
)

type CycleStatus string

const (
	CyclePlanned  CycleStatus = "planned"
	CycleActive   CycleStatus = "active"
	CyclePaused   CycleStatus = "paused"
	CycleCooling  CycleStatus = "cooling"
	CycleFinished CycleStatus = "finished"
	CycleAborted  CycleStatus = "aborted"
)

type Cycle struct {
	ID         string      `json:"id"`
	LoadID     string      `json:"load_id"`
	Profile    string      `json:"profile"`
	Status     CycleStatus `json:"status"`
	StageIndex int         `json:"stage_index"`
	StartedAt  *time.Time  `json:"started_at,omitempty"`
	EndedAt    *time.Time  `json:"ended_at,omitempty"`
	CreatedAt  time.Time   `json:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at"`
	Version    int64       `json:"version"`
}

func (c Cycle) Validate() error {
	if c.ID == "" || c.LoadID == "" || c.Profile == "" {
		return ErrValidation
	}
	if c.StageIndex < 0 {
		return fmt.Errorf("%w: stage index", ErrValidation)
	}
	return nil
}

func (c Cycle) CanAdvance(stageCount int) bool {
	return c.Status == CycleActive && c.StageIndex >= 0 && c.StageIndex+1 < stageCount
}

func (c Cycle) CanFinish(stageCount int) bool {
	return stageCount > 0 && c.StageIndex == stageCount-1 && (c.Status == CycleActive || c.Status == CycleCooling)
}

func (c Cycle) NextStage() Cycle {
	c.StageIndex++
	c.UpdatedAt = time.Now().UTC()
	return c
}
