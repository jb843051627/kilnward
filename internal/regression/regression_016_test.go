package regression

import (
	"testing"
	"github.com/jb843051627/kilnward/internal/model"
)

func TestBug16_FailedGatesDoNotAliasReport(t *testing.T) {
	checks := make([]model.Gate, 1, 2)
	checks[0] = model.Gate{Name: "temperature", Passed: false}
	report := model.GateReport{Checks: checks}
	failed := report.Failed(); failed[0].Name = "display-only"
	if report.Checks[0].Name != "temperature" { t.Fatal("failed gate list aliases report checks") }
}
