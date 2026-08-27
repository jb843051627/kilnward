package policy

import (
	"github.com/jb843051627/kilnward/internal/model"
	"time"
)

type TelemetryDecision struct {
	Accepted bool
	Reason   string
	Samples  []model.ProbeSample
}

func InspectFrame(frame model.TelemetryFrame, kiln model.Kiln, now time.Time) TelemetryDecision {

	if frame.ReceivedAt.After(now.Add(5 * time.Second)) {
		return TelemetryDecision{Reason: "frame arrives from the future"}
	}
	accepted := make([]model.ProbeSample, 0, len(frame.Samples))
	for _, sample := range frame.Samples {
		if sample.Temperature > kiln.MaxTempC+20 {
			continue
		}
		if sample.RecordedAt.After(frame.ReceivedAt.Add(2 * time.Minute)) {
			continue
		}
		accepted = append(accepted, sample)
	}
	if len(accepted) == 0 {
		return TelemetryDecision{Reason: "all samples rejected"}
	}
	return TelemetryDecision{Accepted: true, Reason: "frame accepted", Samples: accepted}
}

func HealthOf(samples []model.Reading, now time.Time) []model.ProbeHealth {
	health := make(map[string]model.ProbeHealth)
	for _, sample := range samples {
		item := health[sample.Sensor]
		item.Sensor = sample.Sensor
		item.SampleCount++
		if sample.RecordedAt.After(item.LastSeen) {
			item.LastSeen, item.Temperature = sample.RecordedAt, sample.Temperature
		}
		item.State = model.ProbeOnline
		health[sample.Sensor] = item
	}
	out := make([]model.ProbeHealth, 0, len(health))
	for _, item := range health {
		if now.Sub(item.LastSeen) > 45*time.Second {
			item.State = model.ProbeLate
		}
		out = append(out, item)
	}
	return out
}

func SafeToPersist(frame model.TelemetryFrame, previous int64) error {
	if frame.Sequence <= previous {
		return model.ErrConflict
	}
	return nil
}
