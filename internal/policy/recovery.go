package policy

import (
	"github.com/jb843051627/kilnward/internal/model"
)

func CommandAllowed(command model.OperatorCommand, kiln model.Kiln, load *model.Load) error {
	if !command.CanRun() {
		return model.ErrInvalidState
	}
	if !kiln.Enabled {
		return model.ErrConflict
	}
	switch command.Action {
	case model.CommandStart:
		if load == nil || load.State != model.LoadAttached {
			return model.ErrConflict
		}
	case model.CommandPause, model.CommandResume, model.CommandAbort:
		if load == nil || (load.State != model.LoadRunning && load.State != model.LoadPaused) {
			return model.ErrInvalidState
		}
	case model.CommandQuarantine:
		if kiln.State == model.KilnQuarantined {
			return model.ErrConflict
		}
	case model.CommandCalibrate:
		if kiln.State != model.KilnReady {
			return model.ErrInvalidState
		}
	default:
		return model.ErrValidation
	}
	return nil
}

func RecoveryState(load model.Load, action model.CommandAction) (model.LoadState, error) {
	switch action {
	case model.CommandStart:
		if load.State != model.LoadAttached {
			return load.State, model.ErrInvalidState
		}
		return model.LoadRunning, nil
	case model.CommandPause:
		if load.State != model.LoadRunning {
			return load.State, model.ErrInvalidState
		}
		return model.LoadPaused, nil
	case model.CommandResume:
		if load.State != model.LoadPaused {
			return load.State, model.ErrInvalidState
		}
		return model.LoadRunning, nil
	case model.CommandAbort:
		if load.State == model.LoadComplete || load.State == model.LoadRejected {
			return load.State, model.ErrInvalidState
		}
		return model.LoadRejected, nil
	default:
		return load.State, nil
	}
}

func Retryable(err error, attempts int) bool {
	return err != nil && attempts < 3 && !model.IsStateError(err)
}
