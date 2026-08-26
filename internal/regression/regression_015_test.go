package regression

import (
	"errors"
	"testing"
	"github.com/jb843051627/kilnward/internal/model"
)

func TestBug15_ProfileStepOutOfRangeReturnsNotFound(t *testing.T) {
	defer func() { if recovered := recover(); recovered != nil { t.Fatalf("out of range step panicked: %v", recovered) } }()
	profile := model.Profile{Steps: []model.ProfileStep{{Sequence: 0, Name: "预热"}}}
	_, err := profile.Step(9)
	if !errors.Is(err, model.ErrNotFound) { t.Fatalf("error = %v, want not found", err) }
}
