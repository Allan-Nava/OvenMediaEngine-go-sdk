package ovenmedia

import (
	"context"
	"net/http"
)

// CreateApplication calls POST /v1/vhosts/{vhost}/apps with an array of
// applications and returns one reply per application.
func (o *ovenMedia) CreateApplication(ctx context.Context, vhost string, apps ...ApplicationConfig) ([]ResponseApplication, error) {
	if len(apps) == 0 {
		return nil, invalid("application: at least one application is required")
	}
	for _, a := range apps {
		if a.Name == "" {
			return nil, invalid("application: name is required")
		}
	}
	res, err := do[[]ResponseApplication](ctx, o, http.MethodPost, pathApps(vhost), apps)
	if err != nil {
		return nil, err
	}
	return *res, nil
}

// GetApplications calls GET /v1/vhosts/{vhost}/apps.
func (o *ovenMedia) GetApplications(ctx context.Context, vhost string) (*ResponseNameList, error) {
	return do[ResponseNameList](ctx, o, http.MethodGet, pathApps(vhost), nil)
}

// GetApplication calls GET /v1/vhosts/{vhost}/apps/{app}.
func (o *ovenMedia) GetApplication(ctx context.Context, vhost, app string) (*ResponseApplication, error) {
	return do[ResponseApplication](ctx, o, http.MethodGet, pathApp(vhost, app), nil)
}

// UpdateApplication calls PATCH /v1/vhosts/{vhost}/apps/{app}. Only the
// fields set in patch are sent.
func (o *ovenMedia) UpdateApplication(ctx context.Context, vhost, app string, patch ApplicationConfig) (*ResponseApplication, error) {
	return do[ResponseApplication](ctx, o, http.MethodPatch, pathApp(vhost, app), patch)
}

// DeleteApplication calls DELETE /v1/vhosts/{vhost}/apps/{app}.
func (o *ovenMedia) DeleteApplication(ctx context.Context, vhost, app string) (*BaseResponseOK, error) {
	return do[BaseResponseOK](ctx, o, http.MethodDelete, pathApp(vhost, app), nil)
}

// CreateOutputProfile calls POST .../apps/{app}/outputProfiles with an array
// of profiles and returns one reply per profile.
func (o *ovenMedia) CreateOutputProfile(ctx context.Context, vhost, app string, profiles ...OutputProfile) ([]ResponseOutputProfile, error) {
	if len(profiles) == 0 {
		return nil, invalid("output profile: at least one profile is required")
	}
	for _, p := range profiles {
		if p.Name == "" {
			return nil, invalid("output profile: name is required")
		}
	}
	res, err := do[[]ResponseOutputProfile](ctx, o, http.MethodPost, pathOutputProfiles(vhost, app), profiles)
	if err != nil {
		return nil, err
	}
	return *res, nil
}

// GetOutputProfiles calls GET .../apps/{app}/outputProfiles.
func (o *ovenMedia) GetOutputProfiles(ctx context.Context, vhost, app string) (*ResponseNameList, error) {
	return do[ResponseNameList](ctx, o, http.MethodGet, pathOutputProfiles(vhost, app), nil)
}

// GetOutputProfile calls GET .../apps/{app}/outputProfiles/{profile}.
func (o *ovenMedia) GetOutputProfile(ctx context.Context, vhost, app, profile string) (*ResponseOutputProfile, error) {
	return do[ResponseOutputProfile](ctx, o, http.MethodGet, pathOutputProfile(vhost, app, profile), nil)
}

// DeleteOutputProfile calls DELETE .../apps/{app}/outputProfiles/{profile}.
func (o *ovenMedia) DeleteOutputProfile(ctx context.Context, vhost, app, profile string) (*BaseResponseOK, error) {
	return do[BaseResponseOK](ctx, o, http.MethodDelete, pathOutputProfile(vhost, app, profile), nil)
}
