package policy

import "github.com/jb843051627/kilnward/internal/model"

func ClassifyTemperature(kiln model.Kiln, stage model.Stage, reading model.Reading) (model.IncidentSeverity, string) {
	if reading.Temperature > kiln.MaxTempC {
		return model.SeverityCritical, "kiln maximum temperature exceeded"
	}
	if reading.Temperature > stage.TargetTempC+stage.ToleranceC {
		return model.SeverityWarning, "stage temperature above tolerance"
	}
	if reading.Temperature < stage.TargetTempC-stage.ToleranceC {
		return model.SeverityWarning, "stage temperature below tolerance"
	}
	return model.SeverityInfo, "temperature accepted"
}

func RequiresQuarantine(severity model.IncidentSeverity) bool {
	return severity == model.SeverityCritical
}

func CanResume(open []model.Incident) bool {
	for _, incident := range open {
		if incident.Status != model.IncidentResolved && incident.IsCritical() {
			return false
		}
	}
	return true
}
