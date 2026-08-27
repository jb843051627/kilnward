package policy

import (
	"fmt"
	"github.com/jb843051627/kilnward/internal/model"
	"time"
)

func DiagnoseKiln(kiln model.Kiln, incidents []model.Incident, maintenance []model.Maintenance, now time.Time) []model.DiagnosticIssue {
	issues := make([]model.DiagnosticIssue, 0)
	if !kiln.Enabled {
		issues = append(issues, model.Issue("KILN_DISABLED", model.DiagnosticFailure, "kiln", "窑炉已停用", kiln.ID))
	}
	if kiln.State == model.KilnQuarantined {
		issues = append(issues, model.Issue("KILN_QUARANTINED", model.DiagnosticFailure, "kiln", "窑炉处于隔离状态", kiln.ID))
	}
	if Due(kiln.LastService, now, 30*24*time.Hour) {
		issues = append(issues, model.Issue("SERVICE_DUE", model.DiagnosticWarning, "maintenance", "维护周期已到期", kiln.LastService.Format(time.RFC3339)))
	}
	for _, item := range maintenance {
		if item.Status == model.MaintenanceActive {
			issues = append(issues, model.Issue("MAINTENANCE_ACTIVE", model.DiagnosticNotice, item.ID, "存在进行中的维护", item.Kind))
		}
	}
	for _, incident := range incidents {
		if incident.Status == model.IncidentResolved {
			continue
		}
		level := model.DiagnosticWarning
		if incident.Severity == model.SeverityCritical {
			level = model.DiagnosticFailure
		}
		issues = append(issues, model.Issue("INCIDENT_"+incident.Code, level, incident.ID, incident.Detail, string(incident.Status)))
	}
	return issues
}

func DiagnoseLoad(load model.Load, cycle model.Cycle, stages []model.Stage, readings []model.Reading, report model.GateReport, now time.Time) []model.DiagnosticIssue {
	issues := make([]model.DiagnosticIssue, 0)
	if load.State == model.LoadRunning && cycle.Status != model.CycleActive {
		issues = append(issues, model.Issue("LOAD_CYCLE_MISMATCH", model.DiagnosticFailure, load.ID, "装载批运行状态与热循环不一致", string(cycle.Status)))
	}
	if load.State == model.LoadComplete && cycle.Status != model.CycleFinished {
		issues = append(issues, model.Issue("LOAD_NOT_FINISHED", model.DiagnosticFailure, load.ID, "装载批已完成但热循环未结束", string(cycle.Status)))
	}
	if len(stages) == 0 {
		issues = append(issues, model.Issue("STAGES_EMPTY", model.DiagnosticFailure, cycle.ID, "热循环没有阶段", "stage_count=0"))
	}
	if !report.Passed && cycle.Status == model.CycleActive {
		issues = append(issues, model.Issue("GATE_BLOCKED", model.DiagnosticWarning, cycle.ID, "当前阶段门未通过", fmt.Sprintf("checks=%d", len(report.Checks))))
	}
	if len(readings) == 0 && !cycle.UpdatedAt.IsZero() && now.Sub(cycle.UpdatedAt) > 30*time.Second {
		issues = append(issues, model.Issue("READINGS_STALE", model.DiagnosticFailure, cycle.ID, "热循环长时间没有采样", cycle.UpdatedAt.Format(time.RFC3339)))
	}
	return issues
}

func DiagnosticsHealthy(issues []model.DiagnosticIssue) bool {
	for _, issue := range issues {
		if issue.Level == model.DiagnosticFailure {
			return false
		}
	}
	return true
}
