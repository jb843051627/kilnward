package model

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"time"
)

type ProfileStatus string

const (
	ProfileDraft     ProfileStatus = "draft"
	ProfileReview    ProfileStatus = "review"
	ProfilePublished ProfileStatus = "published"
	ProfileRetired   ProfileStatus = "retired"
)

type Profile struct {
	ID        string        `json:"id"`
	Name      string        `json:"name"`
	Material  string        `json:"material"`
	Revision  int           `json:"revision"`
	Status    ProfileStatus `json:"status"`
	Author    string        `json:"author"`
	Steps     []ProfileStep `json:"steps"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
}

type ProfileStep struct {
	Sequence       int     `json:"sequence"`
	Name           string  `json:"name"`
	TargetTempC    float64 `json:"target_temp_c"`
	RampPerMinute  float64 `json:"ramp_per_minute"`
	ToleranceC     float64 `json:"tolerance_c"`
	MinHoldSeconds int     `json:"min_hold_seconds"`
	MaxHoldSeconds int     `json:"max_hold_seconds"`
	Atmosphere     float64 `json:"atmosphere"`
}

func (p Profile) Validate() error {
	var problems FieldErrors
	if p.ID == "" {
		problems = append(problems, FieldError{Field: "id", Message: "不能为空"})
	}
	if p.Name == "" {
		problems = append(problems, FieldError{Field: "name", Message: "不能为空"})
	}
	if p.Material == "" {
		problems = append(problems, FieldError{Field: "material", Message: "不能为空"})
	}
	if p.Author == "" {
		problems = append(problems, FieldError{Field: "author", Message: "不能为空"})
	}
	if len(p.Steps) == 0 {
		problems = append(problems, FieldError{Field: "steps", Message: "不能为空"})
	}
	if len(p.Steps) < 2 {
		problems = append(problems, FieldError{Field: "steps", Message: "至少需要两个阶段"})
	}
	steps := append([]ProfileStep(nil), p.Steps...)
	sequences := make([]int, 0, len(steps))
	for _, step := range steps {
		sequences = append(sequences, step.Sequence)
		if err := step.Validate(); err != nil {
			return fmt.Errorf("profile step %d: %w", step.Sequence, err)
		}
	}
	sort.Ints(sequences)
	for i, sequence := range sequences {
		if i != sequence {
			problems = append(problems, FieldError{Field: "steps.sequence", Message: "阶段序号必须连续"})
			break
		}
	}
	if len(problems) > 0 {
		return problems
	}
	return nil
}

func (s ProfileStep) Validate() error {
	if s.Name == "" || s.Sequence < 0 || s.TargetTempC <= 0 || s.RampPerMinute <= 0 {
		return fmt.Errorf("profile step fields: %w", ErrValidation)
	}
	if s.ToleranceC < 0 || s.MinHoldSeconds < 0 || s.MaxHoldSeconds < s.MinHoldSeconds {
		return fmt.Errorf("profile step hold: %w", ErrValidation)
	}
	if s.Atmosphere < 0 || s.Atmosphere > 100 {
		return ErrValidation
	}
	return nil
}

func (p Profile) OrderedSteps() []ProfileStep {
	items := p.Steps
	sort.Slice(items, func(i, j int) bool { return items[i].Sequence < items[j].Sequence })
	return items
}

func (p Profile) Duration() time.Duration {
	var seconds int
	for _, step := range p.Steps {
		seconds += step.MinHoldSeconds
		if step.RampPerMinute > 0 {
			seconds += int(step.TargetTempC / step.RampPerMinute * 60)
		}
	}
	return time.Duration(seconds) * time.Second
}

func (p Profile) Fingerprint() string {
	items := p.Steps
	digest := sha256.New()
	fmt.Fprintf(digest, "%s:%s:%d", p.Name, p.Material, p.Revision)
	for _, step := range items {
		fmt.Fprintf(digest, "|%d:%s:%.2f:%.2f:%.2f:%d:%d:%.2f", step.Sequence, step.Name, step.TargetTempC, step.RampPerMinute, step.ToleranceC, step.MinHoldSeconds, step.MaxHoldSeconds, step.Atmosphere)
	}
	return hex.EncodeToString(digest.Sum(nil))
}

func (p Profile) Step(sequence int) (ProfileStep, error) {
	if len(p.Steps) == 0 {
		return ProfileStep{}, ErrNotFound
	}
	if sequence < 0 || sequence >= len(p.Steps) {
		return ProfileStep{}, ErrNotFound
	}
	return p.Steps[sequence], nil
}
