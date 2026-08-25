package model

import "time"

type IncidentSeverity string

const (
	SeverityInfo     IncidentSeverity = "info"
	SeverityWarning  IncidentSeverity = "warning"
	SeverityCritical IncidentSeverity = "critical"
)

type IncidentStatus string

const (
	IncidentOpen         IncidentStatus = "open"
	IncidentAcknowledged IncidentStatus = "acknowledged"
	IncidentResolved     IncidentStatus = "resolved"
	IncidentDismissed    IncidentStatus = "dismissed"
)

type Incident struct {
	ID             string           `json:"id"`
	KilnID         string           `json:"kiln_id"`
	LoadID         string           `json:"load_id"`
	Code           string           `json:"code"`
	Severity       IncidentSeverity `json:"severity"`
	Status         IncidentStatus   `json:"status"`
	Detail         string           `json:"detail"`
	OpenedAt       time.Time        `json:"opened_at"`
	AcknowledgedAt *time.Time       `json:"acknowledged_at,omitempty"`
	ClosedAt       *time.Time       `json:"closed_at,omitempty"`
	Owner          string           `json:"owner"`
}

func (i Incident) Validate() error {
	if i.ID == "" || i.KilnID == "" || i.Code == "" || i.Detail == "" {
		return ErrValidation
	}
	return nil
}

func (i Incident) CanAcknowledge() bool { return i.Status == IncidentOpen }
func (i Incident) CanResolve() bool {
	return i.Status == IncidentAcknowledged || i.Status == IncidentOpen
}
func (i Incident) IsCritical() bool { return i.Severity == SeverityCritical }

func (i *Incident) Resolve(now time.Time, owner string) error {
	if !i.CanResolve() {
		return ErrInvalidState
	}
	i.Status = IncidentResolved
	i.Owner = owner
	i.ClosedAt = &now
	return nil
}
