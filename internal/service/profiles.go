package service

import (
	"context"
	"fmt"
	"github.com/jb843051627/kilnward/internal/model"
	"github.com/jb843051627/kilnward/internal/policy"
	"time"
)

func (a *App) CreateProfile(ctx context.Context, profile model.Profile) (model.Profile, error) {
	if err := a.requireContext(ctx); err != nil {
		return profile, err
	}
	if profile.ID == "" {
		profile.ID = a.newID("profile")
	}
	if profile.Status == "" {
		profile.Status = model.ProfileDraft
	}
	if profile.Revision == 0 {
		profile.Revision = 1
	}
	if profile.CreatedAt.IsZero() {
		profile.CreatedAt = a.now()
	}
	profile.UpdatedAt = a.now()
	if err := profile.Validate(); err != nil {
		return profile, err
	}
	if err := a.repo.CreateProfile(ctx, profile); err != nil {
		return profile, fmt.Errorf("create profile: %w", err)
	}
	_ = a.audit(ctx, "profile", profile.ID, "created", profile.Fingerprint())
	return profile, nil
}

func (a *App) GetProfile(ctx context.Context, id string) (model.Profile, error) {
	if err := a.requireContext(ctx); err != nil {
		return model.Profile{}, err
	}
	return a.repo.GetProfile(ctx, id)
}
func (a *App) ListProfiles(ctx context.Context, status model.ProfileStatus) ([]model.Profile, error) {
	if err := a.requireContext(ctx); err != nil {
		return nil, err
	}
	return a.repo.ListProfiles(ctx, status)
}

func (a *App) ReviewProfile(ctx context.Context, id string) (model.Profile, error) {
	profile, err := a.repo.GetProfile(ctx, id)
	if err != nil {
		return profile, err
	}
	if decision := policy.ReviewProfile(profile); !decision.Allowed {
		return profile, model.ErrValidation
	}
	if err := a.repo.UpdateProfileStatus(ctx, id, model.ProfileReview, a.now().Format(time.RFC3339Nano)); err != nil {
		return profile, err
	}
	_ = a.audit(ctx, "profile", id, "reviewed", "quality review")
	return a.repo.GetProfile(ctx, id)
}

func (a *App) PublishProfile(ctx context.Context, id string) (model.Profile, error) {
	profile, err := a.repo.GetProfile(ctx, id)
	if err != nil {
		return profile, err
	}
	if err := policy.PublishProfile(profile); err != nil {
		return profile, err
	}
	if err := a.repo.UpdateProfileStatus(ctx, id, model.ProfilePublished, a.now().Format(time.RFC3339Nano)); err != nil {
		return profile, err
	}
	_ = a.audit(ctx, "profile", id, "published", profile.Fingerprint())
	return a.repo.GetProfile(ctx, id)
}

func (a *App) RetireProfile(ctx context.Context, id string) (model.Profile, error) {
	profile, err := a.repo.GetProfile(ctx, id)
	if err != nil {
		return profile, err
	}
	if profile.Status != model.ProfilePublished {
		return profile, model.ErrInvalidState
	}
	if err := a.repo.UpdateProfileStatus(ctx, id, model.ProfileRetired, a.now().Format(time.RFC3339Nano)); err != nil {
		return profile, err
	}
	_ = a.audit(ctx, "profile", id, "retired", "operator request")
	return a.repo.GetProfile(ctx, id)
}

func (a *App) CheckProfileForKiln(ctx context.Context, profileID, kilnID string) error {
	profile, err := a.repo.GetProfile(ctx, profileID)
	if err != nil {
		return err
	}
	kiln, err := a.repo.GetKiln(ctx, kilnID)
	if err != nil {
		return err
	}
	return policy.ProfileForKiln(profile, kiln)
}
