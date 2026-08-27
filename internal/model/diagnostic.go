package model

import "time"

type DiagnosticLevel string

const (
	DiagnosticOK      DiagnosticLevel = "ok"
	DiagnosticNotice  DiagnosticLevel = "notice"
	DiagnosticWarning DiagnosticLevel = "warning"
	DiagnosticFailure DiagnosticLevel = "failure"
)

type DiagnosticIssue struct {
	Code     string          `json:"code"`
	Level    DiagnosticLevel `json:"level"`
	Subject  string          `json:"subject"`
	Message  string          `json:"message"`
	Evidence string          `json:"evidence"`
}

type DiagnosticReport struct {
	ID          string            `json:"id"`
	KilnID      string            `json:"kiln_id"`
	GeneratedAt time.Time         `json:"generated_at"`
	DurationMS  int64             `json:"duration_ms"`
	Healthy     bool              `json:"healthy"`
	Issues      []DiagnosticIssue `json:"issues"`
}

func (r DiagnosticReport) FailureCount() int {
	count := 0
	for _, issue := range r.Issues {
		if issue.Level == DiagnosticFailure {
			count++
		}
	}
	return count
}
func (r DiagnosticReport) WarningCount() int {
	count := 0
	for _, issue := range r.Issues {
		if issue.Level == DiagnosticWarning {
			count++
		}
	}
	return count
}
func (r DiagnosticReport) Add(issue DiagnosticIssue) DiagnosticReport {
	r.Issues = append(r.Issues, issue)
	if issue.Level == DiagnosticFailure {
		r.Healthy = false
	}
	return r
}

func (i DiagnosticIssue) Valid() bool { return i.Code != "" && i.Subject != "" && i.Message != "" }

func Issue(code string, level DiagnosticLevel, subject, message, evidence string) DiagnosticIssue {
	return DiagnosticIssue{Code: code, Level: level, Subject: subject, Message: message, Evidence: evidence}
}
