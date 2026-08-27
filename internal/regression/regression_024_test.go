package regression

import (
	"testing"

	"github.com/jb843051627/kilnward/internal/model"
	"github.com/jb843051627/kilnward/internal/policy"
)

func TestBug24_ProfileStepValidationPropagates(t *testing.T) {
	profile := model.Profile{ID: "p1", Name: "kiln", Material: "clay", Author: "op", Steps: []model.ProfileStep{{Sequence: 0, Name: "预热", TargetTempC: 300, RampPerMinute: -1}, {Sequence: 1, Name: "升温", TargetTempC: 800, RampPerMinute: 5}}}
	if err := profile.Validate(); err == nil { t.Fatal("invalid profile step was accepted") }
	if decision := policy.ReviewProfile(profile); decision.Allowed { t.Fatal("invalid profile entered review") }
}
