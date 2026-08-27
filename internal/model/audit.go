package model

import "time"

type Audit struct {
	ID          int64     `json:"id"`
	SubjectType string    `json:"subject_type"`
	SubjectID   string    `json:"subject_id"`
	Action      string    `json:"action"`
	Detail      string    `json:"detail"`
	CreatedAt   time.Time `json:"created_at"`
}

func (a Audit) Validate() error {
	if a.SubjectType == "" || a.SubjectID == "" || a.Action == "" {
		return ErrValidation
	}
	return nil
}
