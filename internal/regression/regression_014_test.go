package regression

import (
	"testing"
	"github.com/jb843051627/kilnward/internal/model"
	"github.com/jb843051627/kilnward/internal/policy"
)

func TestBug14_NegativeTimelineLimitDoesNotPanic(t *testing.T) {
	entries := []model.TimelineEntry{{Action: "one"}, {Action: "two"}}
	defer func() { if recovered := recover(); recovered != nil { t.Fatalf("negative limit panicked: %v", recovered) } }()
	if got := policy.Limit(entries, -1); len(got) != len(entries) { t.Fatalf("got %d entries, want %d", len(got), len(entries)) }
}
