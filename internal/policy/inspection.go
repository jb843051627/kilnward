package policy

import (
	"fmt"
	"github.com/jb843051627/kilnward/internal/model"
	"strings"
	"time"
)

type InspectionRule struct {
	Code        string
	Title       string
	Required    bool
	Severity    model.DiagnosticLevel
	Check       func() bool
	Explanation string
}

type InspectionResult struct {
	Passed bool
	Rules  []InspectionRule
	Issues []model.DiagnosticIssue
}

func InspectLoad(load model.Load, kiln model.Kiln, cycle model.Cycle, readings []model.Reading, now time.Time) InspectionResult {
	result := InspectionResult{Passed: true, Rules: make([]InspectionRule, 0)}
	result.Rules = append(result.Rules,
		InspectionRule{Code: "kiln-link", Title: "load points to its kiln", Required: true, Severity: model.DiagnosticFailure, Check: func() bool { return load.KilnID == kiln.ID }, Explanation: "装载批必须属于当前窑炉"},
		InspectionRule{Code: "profile-target", Title: "target is within kiln capability", Required: true, Severity: model.DiagnosticFailure, Check: func() bool { return load.TargetTempC <= kiln.MaxTempC }, Explanation: "目标温度不能超过窑炉上限"},
		InspectionRule{Code: "cycle-link", Title: "cycle belongs to load", Required: true, Severity: model.DiagnosticFailure, Check: func() bool { return cycle.LoadID == load.ID }, Explanation: "热循环必须属于当前装载批"},
		InspectionRule{Code: "reading-recent", Title: "readings are recent", Required: load.State == model.LoadRunning, Severity: model.DiagnosticWarning, Check: func() bool { return recentReading(readings, now) }, Explanation: "运行中的装载批需要近期采样"},
		InspectionRule{Code: "state-alignment", Title: "load and cycle states align", Required: true, Severity: model.DiagnosticFailure, Check: func() bool { return statesAlign(load, cycle) }, Explanation: "装载批和热循环状态必须同步"},
	)
	for _, rule := range result.Rules {
		if rule.Check == nil || rule.Check() {
			continue
		}
		level := rule.Severity
		result.Issues = append(result.Issues, model.Issue(strings.ToUpper(rule.Code), level, load.ID, rule.Explanation, rule.Title))
		if rule.Required {
			result.Passed = false
		}
	}
	return result
}

func recentReading(readings []model.Reading, now time.Time) bool {
	for _, reading := range readings {
		if reading.IsFresh(now, 45*time.Second) && reading.Quality == model.QualityGood {
			return true
		}
	}
	return false
}

func statesAlign(load model.Load, cycle model.Cycle) bool {
	switch load.State {
	case model.LoadDraft, model.LoadAttached:
		return cycle.Status == model.CyclePlanned
	case model.LoadRunning:
		return cycle.Status == model.CycleActive
	case model.LoadPaused:
		return cycle.Status == model.CyclePaused
	case model.LoadCooling:
		return cycle.Status == model.CycleCooling
	case model.LoadComplete:
		return cycle.Status == model.CycleFinished
	case model.LoadRejected:
		return cycle.Status == model.CycleAborted || cycle.Status == model.CycleFinished
	default:
		return false
	}
}

func InspectProfile(profile model.Profile, kiln model.Kiln) InspectionResult {
	result := InspectionResult{Passed: true, Rules: make([]InspectionRule, 0)}
	rules := []InspectionRule{
		{Code: "profile-valid", Title: "profile validates", Required: true, Severity: model.DiagnosticFailure, Check: func() bool { return profile.Validate() == nil }, Explanation: "工艺曲线字段和阶段顺序必须完整"},
		{Code: "profile-published", Title: "profile is published", Required: true, Severity: model.DiagnosticFailure, Check: func() bool { return profile.Status == model.ProfilePublished }, Explanation: "只有已发布曲线可以用于生产"},
		{Code: "profile-capability", Title: "profile fits kiln", Required: true, Severity: model.DiagnosticFailure, Check: func() bool {
			return len(profile.Steps) > 0 && profile.Steps[len(profile.Steps)-1].TargetTempC <= kiln.MaxTempC
		}, Explanation: "曲线最高温度不能超过窑炉能力"},
		{Code: "profile-revision", Title: "profile has a positive revision", Required: false, Severity: model.DiagnosticWarning, Check: func() bool { return profile.Revision > 0 }, Explanation: "曲线修订号应为正数"},
	}
	result.Rules = append(result.Rules, rules...)
	for _, rule := range result.Rules {
		if rule.Check != nil && rule.Check() {
			continue
		}
		result.Issues = append(result.Issues, model.Issue(strings.ToUpper(rule.Code), rule.Severity, profile.ID, rule.Explanation, rule.Title))
		if rule.Required {
			result.Passed = false
		}
	}
	return result
}

func InspectProbe(kiln model.Kiln, readings []model.Reading, now time.Time) InspectionResult {
	result := InspectionResult{Passed: true, Rules: make([]InspectionRule, 0)}
	checks := []struct {
		code, title, message string
		ok                   bool
	}{
		{"probe-count", "probe count is configured", "窑炉必须配置探头数量", kiln.ProbeCount > 0},
		{"probe-presence", "at least one probe reported", "当前周期没有收到探头采样", len(readings) > 0},
		{"probe-freshness", "probe data is fresh", "探头数据已过期", recentReading(readings, now)},
	}
	for _, check := range checks {
		rule := InspectionRule{Code: check.code, Title: check.title, Required: true, Severity: model.DiagnosticFailure, Explanation: check.message}
		result.Rules = append(result.Rules, rule)
		if check.ok {
			continue
		}
		result.Passed = false
		result.Issues = append(result.Issues, model.Issue(strings.ToUpper(check.code), rule.Severity, kiln.ID, check.message, check.title))
	}
	return result
}

func InspectionSummary(results ...InspectionResult) (bool, string) {
	passed := true
	parts := make([]string, 0)
	for _, result := range results {
		if !result.Passed {
			passed = false
		}
		for _, issue := range result.Issues {
			parts = append(parts, fmt.Sprintf("%s:%s", issue.Code, issue.Message))
		}
	}
	return passed, strings.Join(parts, "; ")
}
