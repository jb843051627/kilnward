package model

import "time"

type CalibrationStatus string

const (
	CalibrationPending CalibrationStatus = "pending"
	CalibrationRunning CalibrationStatus = "running"
	CalibrationPassed  CalibrationStatus = "passed"
	CalibrationFailed  CalibrationStatus = "failed"
	CalibrationExpired CalibrationStatus = "expired"
)

type Calibration struct {
	ID          string            `json:"id"`
	KilnID      string            `json:"kiln_id"`
	Sensor      string            `json:"sensor"`
	ReferenceC  float64           `json:"reference_c"`
	ObservedC   float64           `json:"observed_c"`
	ToleranceC  float64           `json:"tolerance_c"`
	Status      CalibrationStatus `json:"status"`
	Operator    string            `json:"operator"`
	Note        string            `json:"note"`
	CreatedAt   time.Time         `json:"created_at"`
	CompletedAt *time.Time        `json:"completed_at,omitempty"`
}

func (c Calibration) Validate() error {
	if c.ID == "" || c.KilnID == "" || c.Sensor == "" || c.Operator == "" {
		return ErrValidation
	}
	if c.ReferenceC < -50 || c.ReferenceC > 2500 || c.ObservedC < -50 || c.ObservedC > 2500 || c.ToleranceC < 0 {
		return ErrValidation
	}
	return nil
}

func (c Calibration) Drift() float64 { return c.ObservedC - c.ReferenceC }
func (c Calibration) WithinTolerance() bool {
	drift := c.Drift()
	return drift >= -c.ToleranceC && drift <= c.ToleranceC
}
func (c Calibration) CanComplete() bool {
	return c.Status == CalibrationPending || c.Status == CalibrationRunning
}

func (c *Calibration) Complete(now time.Time) error {
	if !c.CanComplete() {
		return ErrInvalidState
	}
	c.CompletedAt = &now
	if c.WithinTolerance() {
		c.Status = CalibrationPassed
	} else {
		c.Status = CalibrationFailed
	}
	return nil
}

func (c Calibration) Expired(now time.Time, maxAge time.Duration) bool {
	return c.CompletedAt == nil || now.Sub(*c.CompletedAt) > maxAge || c.Status != CalibrationPassed
}
