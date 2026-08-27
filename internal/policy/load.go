package policy

import "github.com/jb843051627/kilnward/internal/model"

func AttachAllowed(kiln model.Kiln, load model.Load) error {
	if !kiln.CanAcceptLoad() {
		return model.ErrInvalidState
	}
	if load.State != model.LoadDraft {
		return model.ErrConflict
	}
	if load.KilnID != kiln.ID {
		return model.ErrConflict
	}
	return nil
}

func StartAllowed(kiln model.Kiln, load model.Load, gates model.GateReport) error {
	if !kiln.CanStart() {
		return model.ErrInvalidState
	}
	if load.State != model.LoadAttached || !gates.Passed {
		return model.ErrConflict
	}
	return nil
}

func CompletionAllowed(load model.Load, cycle model.Cycle, report model.GateReport) error {
	if load.State != model.LoadCooling || cycle.Status != model.CycleCooling || !report.Passed {
		return model.ErrConflict
	}
	return nil
}
