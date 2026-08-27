package policy

import (
	"github.com/jb843051627/kilnward/internal/model"
	"strings"
)

func RequiredChecks(stage model.Stage) []model.GateKind {
	checks := []model.GateKind{model.GateTemperature, model.GateSampling}
	if stage.TargetTempC > 1000 {
		checks = append(checks, model.GateCalibration)
	}
	if stage.Sequence > 0 {
		checks = append(checks, model.GateSequence)
	}
	return checks
}

func Names(kinds []model.GateKind) string {
	parts := make([]string, 0, len(kinds))
	for _, kind := range kinds {
		parts = append(parts, string(kind))
	}
	return strings.Join(parts, ",")
}

func Passed(checks []model.Gate) bool {
	for _, check := range checks {
		if !check.Passed {
			return false
		}
	}
	return len(checks) > 0
}
