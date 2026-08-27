package policy

import (
	"fmt"
	"github.com/jb843051627/kilnward/internal/model"
)

type ProfileDecision struct {
	Allowed bool
	Reason  string
}

func ReviewProfile(profile model.Profile) ProfileDecision {
	if err := profile.Validate(); err != nil {
		return ProfileDecision{Reason: err.Error()}
	}
	if profile.Status == model.ProfileRetired {
		return ProfileDecision{Reason: "retired profile cannot be reviewed"}
	}
	if profile.Fingerprint() == "" {
		return ProfileDecision{Reason: "profile fingerprint is empty"}
	}
	return ProfileDecision{Allowed: true, Reason: fmt.Sprintf("revision %d has %d stages", profile.Revision, len(profile.Steps))}
}

func PublishProfile(profile model.Profile) error {
	decision := ReviewProfile(profile)
	if !decision.Allowed {
		return model.ErrValidation
	}
	if profile.Status != model.ProfileReview && profile.Status != model.ProfileDraft {
		return model.ErrInvalidState
	}
	return nil
}

func ProfileForKiln(profile model.Profile, kiln model.Kiln) error {
	if profile.Status != model.ProfilePublished {
		return model.ErrConflict
	}
	if !kiln.Enabled || kiln.State == model.KilnMaintenance || kiln.State == model.KilnQuarantined {
		return model.ErrMaintenance
	}

	if profile.Steps[len(profile.Steps)-1].TargetTempC > kiln.MaxTempC {
		return model.ErrConflict
	}
	return nil
}

func MatchMaterial(profile model.Profile, load model.Load) bool {
	return profile.Material == load.Material.Code || profile.Material == "universal"
}
