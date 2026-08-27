package policy

import (
	"github.com/jb843051627/kilnward/internal/model"
	"time"
)

func CanEnterStage(cycle model.Cycle, stage model.Stage, prior *model.Stage, now time.Time) error {
	if cycle.Status != model.CycleActive {
		return model.ErrInvalidState
	}
	if stage.Sequence != cycle.StageIndex {
		return model.ErrConflict
	}
	if prior != nil && prior.Status != model.StageDone {
		return model.ErrConflict
	}
	if stage.HoldExpired(now) {
		return model.ErrConflict
	}
	return nil
}

func CompleteStage(stage *model.Stage, now time.Time) error {
	if stage.Status != model.StageRunning && stage.Status != model.StageHeld {
		return model.ErrInvalidState
	}
	if stage.StartedAt == nil {
		return model.ErrConflict
	}
	if now.Sub(*stage.StartedAt) < time.Duration(stage.MinHoldSeconds)*time.Second {
		return model.ErrConflict
	}
	stage.Status = model.StageDone
	stage.EndedAt = &now
	return nil
}

func HoldStage(stage *model.Stage) error {
	if stage.Status != model.StageRunning {
		return model.ErrInvalidState
	}
	stage.Status = model.StageHeld
	return nil
}
