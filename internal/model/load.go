package model

import "time"

type LoadState string

const (
	LoadDraft    LoadState = "draft"
	LoadAttached LoadState = "attached"
	LoadRunning  LoadState = "running"
	LoadPaused   LoadState = "paused"
	LoadCooling  LoadState = "cooling"
	LoadComplete LoadState = "complete"
	LoadRejected LoadState = "rejected"
)

type Material struct {
	Code     string  `json:"code"`
	Quantity int     `json:"quantity"`
	Moisture float64 `json:"moisture"`
}

type Load struct {
	ID           string     `json:"id"`
	KilnID       string     `json:"kiln_id"`
	Label        string     `json:"label"`
	State        LoadState  `json:"state"`
	Material     Material   `json:"material"`
	Profile      string     `json:"profile"`
	TargetTempC  float64    `json:"target_temp_c"`
	CurrentStage int        `json:"current_stage"`
	StartedAt    *time.Time `json:"started_at,omitempty"`
	FinishedAt   *time.Time `json:"finished_at,omitempty"`
	LastError    string     `json:"last_error,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	Version      int64      `json:"version"`
}

func (l Load) Validate() error {
	var problems FieldErrors
	for _, item := range []struct{ field, value string }{
		{"id", l.ID}, {"kiln_id", l.KilnID}, {"label", l.Label}, {"profile", l.Profile},
	} {
		if err := RequireID(item.field, item.value); err != nil {
			problems = append(problems, err.(FieldError))
		}
	}
	if l.Material.Code == "" {
		problems = append(problems, FieldError{Field: "material.code", Message: "不能为空"})
	}
	if l.Material.Quantity <= 0 {
		problems = append(problems, FieldError{Field: "material.quantity", Message: "必须为正数"})
	}
	if l.TargetTempC < 200 || l.TargetTempC > 2100 {
		problems = append(problems, FieldError{Field: "target_temp_c", Message: "超出工艺范围"})
	}
	if len(problems) > 0 {
		return problems
	}
	return nil
}

func (l Load) CanTransition(to LoadState) bool {
	switch l.State {
	case LoadDraft:
		return to == LoadAttached || to == LoadRejected
	case LoadAttached:
		return to == LoadRunning || to == LoadRejected
	case LoadRunning:
		return to == LoadPaused || to == LoadCooling || to == LoadRejected
	case LoadPaused:
		return to == LoadRunning || to == LoadRejected
	case LoadCooling:
		return to == LoadComplete || to == LoadRejected
	case LoadComplete, LoadRejected:
		return false
	default:
		return false
	}
}
