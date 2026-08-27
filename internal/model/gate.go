package model

import "time"

type GateKind string

const (
	GateTemperature GateKind = "temperature"
	GateSampling    GateKind = "sampling"
	GateCalibration GateKind = "calibration"
	GateSequence    GateKind = "sequence"
)

type Gate struct {
	ID        string    `json:"id"`
	CycleID   string    `json:"cycle_id"`
	StageSeq  int       `json:"stage_seq"`
	Name      string    `json:"name"`
	Kind      GateKind  `json:"kind"`
	Passed    bool      `json:"passed"`
	Reason    string    `json:"reason"`
	CheckedAt time.Time `json:"checked_at"`
	CheckedBy string    `json:"checked_by"`
}

type GateReport struct {
	CycleID string `json:"cycle_id"`
	Stage   int    `json:"stage"`
	Passed  bool   `json:"passed"`
	Checks  []Gate `json:"checks"`
}

func (r GateReport) Failed() []Gate {
	items := make([]Gate, 0)
	for _, check := range r.Checks {
		if !check.Passed {
			items = append(items, check)
		}
	}
	return items
}

func (r GateReport) FailedNames() []string {
	names := make([]string, 0)
	for _, check := range r.Checks {
		if !check.Passed {
			names = append(names, check.Name)
		}
	}
	return names
}

func (g Gate) Validate() error {
	if g.ID == "" || g.CycleID == "" || g.Name == "" || g.CheckedBy == "" {
		return ErrValidation
	}
	return nil
}
