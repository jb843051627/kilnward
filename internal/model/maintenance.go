package model

import "time"

type MaintenanceStatus string

const (
	MaintenancePlanned  MaintenanceStatus = "planned"
	MaintenanceActive   MaintenanceStatus = "active"
	MaintenanceComplete MaintenanceStatus = "complete"
	MaintenanceCanceled MaintenanceStatus = "canceled"
)

type Maintenance struct {
	ID         string            `json:"id"`
	KilnID     string            `json:"kiln_id"`
	Kind       string            `json:"kind"`
	Status     MaintenanceStatus `json:"status"`
	Note       string            `json:"note"`
	OpenedAt   time.Time         `json:"opened_at"`
	ClosedAt   *time.Time        `json:"closed_at,omitempty"`
	Technician string            `json:"technician"`
}

func (m Maintenance) Validate() error {
	if m.ID == "" || m.KilnID == "" || m.Kind == "" || m.Technician == "" {
		return ErrValidation
	}
	return nil
}

func (m Maintenance) CanStart() bool    { return m.Status == MaintenancePlanned }
func (m Maintenance) CanComplete() bool { return m.Status == MaintenanceActive }

func (m *Maintenance) Complete(now time.Time) error {
	if !m.CanComplete() {
		return ErrInvalidState
	}
	m.Status = MaintenanceComplete
	m.ClosedAt = &now
	return nil
}
