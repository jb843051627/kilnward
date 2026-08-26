package regression

import (
	"errors"
	"testing"
	"github.com/jb843051627/kilnward/internal/model"
	"github.com/jb843051627/kilnward/internal/policy"
)

func TestBug12_EmptyProfileDoesNotPanic(t *testing.T) {
	defer func() { if recovered := recover(); recovered != nil { t.Fatalf("empty profile panicked: %v", recovered) } }()
	profile := model.Profile{ID: "p-1", Name: "empty", Material: "clay", Author: "operator", Status: model.ProfilePublished}
	kiln := model.Kiln{ID: "k-1", Enabled: true, State: model.KilnReady, MaxTempC: 1500, ProbeCount: 2}
	err := policy.ProfileForKiln(profile, kiln)
	if !errors.Is(err, model.ErrValidation) { t.Fatalf("error = %v, want validation", err) }
}
