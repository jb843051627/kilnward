package regression

import (
	"strings"
	"testing"
	"time"

	"github.com/jb843051627/kilnward/internal/model"
	"github.com/jb843051627/kilnward/internal/policy"
)

func TestBug23_TelemetrySampleValidationPropagates(t *testing.T) {
	frame := model.TelemetryFrame{ID: "f1", KilnID: "k1", LoadID: "l1", CycleID: "c1", Gateway: "gw", Samples: []model.ProbeSample{{Sensor: "P1", Temperature: 3000, Atmosphere: 50, Power: 10}}}
	if err := frame.Validate(); err == nil { t.Fatal("invalid telemetry sample was accepted") }
	decision := policy.InspectFrame(frame, model.Kiln{MaxTempC: 1500}, time.Now().UTC())
	if decision.Accepted || !strings.Contains(decision.Reason, "probe sample") { t.Fatalf("invalid telemetry decision = %#v", decision) }
}
