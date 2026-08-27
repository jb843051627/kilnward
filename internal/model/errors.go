package model

import "errors"

var (
	ErrNotFound       = errors.New("kilnward: record not found")
	ErrConflict       = errors.New("kilnward: state conflict")
	ErrInvalidState   = errors.New("kilnward: invalid state")
	ErrValidation     = errors.New("kilnward: validation failed")
	ErrMaintenance    = errors.New("kilnward: kiln is under maintenance")
	ErrMissingReading = errors.New("kilnward: required reading is missing")
)

type FieldError struct {
	Field   string
	Message string
}

func (e FieldError) Error() string {
	return e.Field + ": " + e.Message
}

type FieldErrors []FieldError

func (e FieldErrors) Error() string {
	if len(e) == 0 {
		return ErrValidation.Error()
	}
	return e[0].Error()
}

func (e FieldErrors) Unwrap() error { return ErrValidation }

func IsStateError(err error) bool {
	return errors.Is(err, ErrInvalidState) || errors.Is(err, ErrConflict)
}
