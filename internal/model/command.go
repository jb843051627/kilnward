package model

import "time"

type CommandAction string

const (
	CommandStart      CommandAction = "start"
	CommandPause      CommandAction = "pause"
	CommandResume     CommandAction = "resume"
	CommandAbort      CommandAction = "abort"
	CommandQuarantine CommandAction = "quarantine"
	CommandCalibrate  CommandAction = "calibrate"
)

type CommandStatus string

const (
	CommandQueued    CommandStatus = "queued"
	CommandRunning   CommandStatus = "running"
	CommandSucceeded CommandStatus = "succeeded"
	CommandFailed    CommandStatus = "failed"
	CommandCanceled  CommandStatus = "canceled"
)

type OperatorCommand struct {
	ID          string        `json:"id"`
	KilnID      string        `json:"kiln_id"`
	LoadID      string        `json:"load_id"`
	Action      CommandAction `json:"action"`
	Status      CommandStatus `json:"status"`
	RequestedBy string        `json:"requested_by"`
	Reason      string        `json:"reason"`
	Attempts    int           `json:"attempts"`
	ErrorText   string        `json:"error_text"`
	CreatedAt   time.Time     `json:"created_at"`
	StartedAt   *time.Time    `json:"started_at,omitempty"`
	FinishedAt  *time.Time    `json:"finished_at,omitempty"`
}

func (c OperatorCommand) Validate() error {
	if c.ID == "" || c.KilnID == "" || c.RequestedBy == "" || c.Reason == "" {
		return ErrValidation
	}
	if c.Action == "" {
		return ErrValidation
	}
	return nil
}

func (c OperatorCommand) CanRun() bool {
	return c.Status == CommandQueued || (c.Status == CommandFailed && c.Attempts < 3)
}
func (c OperatorCommand) CanCancel() bool {
	return c.Status == CommandQueued || c.Status == CommandRunning
}

func (c *OperatorCommand) Begin(now time.Time) error {
	if !c.CanRun() {
		return ErrInvalidState
	}
	c.Status, c.Attempts, c.StartedAt = CommandRunning, c.Attempts+1, &now
	return nil
}
func (c *OperatorCommand) Succeed(now time.Time) error {
	if c.Status != CommandRunning {
		return ErrInvalidState
	}
	c.Status, c.FinishedAt, c.ErrorText = CommandSucceeded, &now, ""
	return nil
}
func (c *OperatorCommand) Fail(now time.Time, err error) {
	c.Status, c.FinishedAt, c.ErrorText = CommandFailed, &now, err.Error()
}
func (c *OperatorCommand) Cancel(now time.Time) error {
	if !c.CanCancel() {
		return ErrInvalidState
	}
	c.Status, c.FinishedAt = CommandCanceled, &now
	return nil
}
