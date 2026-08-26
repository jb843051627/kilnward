package model

import "time"

type KilnState string

const (
	KilnOffline     KilnState = "offline"
	KilnReady       KilnState = "ready"
	KilnHeating     KilnState = "heating"
	KilnCooling     KilnState = "cooling"
	KilnMaintenance KilnState = "maintenance"
	KilnQuarantined KilnState = "quarantined"
)

type Kiln struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Location    string    `json:"location"`
	State       KilnState `json:"state"`
	MaxTempC    float64   `json:"max_temp_c"`
	ProbeCount  int       `json:"probe_count"`
	Enabled     bool      `json:"enabled"`
	LastService time.Time `json:"last_service"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Version     int64     `json:"version"`
}

func (k Kiln) Validate() error {
	var problems FieldErrors
	if err := RequireID("id", k.ID); err != nil {
		problems = append(problems, err.(FieldError))
	}
	if err := RequireID("name", k.Name); err != nil {
		problems = append(problems, err.(FieldError))
	}
	if k.MaxTempC < 300 || k.MaxTempC > 2200 {
		problems = append(problems, FieldError{Field: "max_temp_c", Message: "必须在 300 到 2200 之间"})
	}
	if k.ProbeCount < 1 || k.ProbeCount > 32 {
		problems = append(problems, FieldError{Field: "probe_count", Message: "必须在 1 到 32 之间"})
	}
	if k.State == "" {
		problems = append(problems, FieldError{Field: "state", Message: "不能为空"})
	}
	if len(problems) > 0 {
		return problems
	}
	return nil
}

func (k Kiln) CanStart() bool {
	return k.Enabled && k.State == KilnReady
}

func (k Kiln) CanAcceptLoad() bool {
	return k.Enabled && (k.State == KilnReady || k.State == KilnCooling)
}

func KilnTransition(from, to KilnState) bool {
	switch from {
	case KilnOffline:
		return to == KilnReady || to == KilnMaintenance
	case KilnReady:
		return to == KilnHeating || to == KilnMaintenance || to == KilnOffline || to == KilnQuarantined
	case KilnHeating:
		return to == KilnCooling || to == KilnQuarantined || to == KilnMaintenance
	case KilnCooling:
		return to == KilnReady || to == KilnHeating || to == KilnQuarantined
	case KilnMaintenance:
		return to == KilnReady || to == KilnOffline
	case KilnQuarantined:
		return to == KilnMaintenance || to == KilnOffline
	default:
		return false
	}
}
