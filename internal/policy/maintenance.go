package policy

import (
	"github.com/jb843051627/kilnward/internal/model"
	"time"
)

func BlocksOperation(kiln model.Kiln, active []model.Maintenance) bool {
	if kiln.State == model.KilnMaintenance || kiln.State == model.KilnQuarantined {
		return true
	}
	for _, item := range active {
		if item.Status == model.MaintenanceActive {
			return true
		}
	}
	return false
}

func Due(last, now time.Time, interval time.Duration) bool {
	return last.IsZero() || now.Sub(last) >= interval
}

func NextState(current model.KilnState, status model.MaintenanceStatus) model.KilnState {
	if status == model.MaintenanceActive {
		return model.KilnMaintenance
	}
	if current == model.KilnMaintenance && status == model.MaintenanceComplete {
		return model.KilnReady
	}
	return current
}
