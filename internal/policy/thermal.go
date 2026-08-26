package policy

import (
	"fmt"
	"github.com/jb843051627/kilnward/internal/model"
)

type ThermalDecision struct {
	Allowed bool
	Reason  string
	Delta   float64
}

func EvaluateTemperature(stage model.Stage, reading model.Reading) ThermalDecision {
	delta := reading.Temperature - stage.TargetTempC
	if reading.Quality != model.QualityGood {
		return ThermalDecision{Reason: "reading quality is not good", Delta: delta}
	}
	if !stage.WithinTemperature(reading.Temperature) {
		return ThermalDecision{Reason: fmt.Sprintf("temperature delta %.2f exceeds tolerance", delta), Delta: delta}
	}
	return ThermalDecision{Allowed: true, Reason: "temperature within tolerance", Delta: delta}
}

func IsOverheat(kiln model.Kiln, reading model.Reading) bool {
	return reading.Temperature > kiln.MaxTempC
}

func PowerBudget(stage model.Stage, current float64) (float64, error) {
	if current < 0 || current > 100 {
		return 0, model.ErrValidation
	}
	margin := 100 - current
	if stage.TargetTempC > 1200 {
		margin *= 0.85
	}
	return margin, nil
}
