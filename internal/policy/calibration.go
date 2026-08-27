package policy

import (
	"github.com/jb843051627/kilnward/internal/model"
	"time"
)

type CalibrationDecision struct {
	Allowed bool
	Reason  string
}

func StartCalibration(kiln model.Kiln, sensor string) CalibrationDecision {
	if !kiln.Enabled {
		return CalibrationDecision{Reason: "kiln disabled"}
	}
	if kiln.State != model.KilnReady {
		return CalibrationDecision{Reason: "kiln must be ready"}
	}
	if sensor == "" {
		return CalibrationDecision{Reason: "sensor is empty"}
	}
	return CalibrationDecision{Allowed: true, Reason: "calibration may start"}
}

func EvaluateCalibration(item model.Calibration) CalibrationDecision {
	if err := item.Validate(); err != nil {
		return CalibrationDecision{Reason: err.Error()}
	}
	if item.ToleranceC == 0 {
		return CalibrationDecision{Reason: "tolerance must be positive"}
	}
	if item.WithinTolerance() {
		return CalibrationDecision{Allowed: true, Reason: "probe drift is within tolerance"}
	}
	return CalibrationDecision{Reason: "probe drift exceeds tolerance"}
}

func RequiredForStage(stage model.Stage, calibrations []model.Calibration, now time.Time) bool {
	for _, item := range calibrations {
		if item.Sensor != "" && !item.Expired(now, 24*time.Hour) {
			return true
		}
	}
	return false
}

func DriftSeverity(item model.Calibration) model.IncidentSeverity {
	drift := item.Drift()
	if drift < 0 {
		drift = -drift
	}
	if drift >= item.ToleranceC*3 {
		return model.SeverityCritical
	}
	if drift > item.ToleranceC {
		return model.SeverityWarning
	}
	return model.SeverityInfo
}
