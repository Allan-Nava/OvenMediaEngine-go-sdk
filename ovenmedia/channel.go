package ovenmedia

import (
	"context"
	"net/http"
)

// CreateScheduledChannel calls POST .../apps/{app}/scheduledChannels.
func (o *ovenMedia) CreateScheduledChannel(ctx context.Context, vhost, app string, channel ScheduledChannel) (*BaseResponseOK, error) {
	if err := channel.validate(); err != nil {
		return nil, err
	}
	return do[BaseResponseOK](ctx, o, http.MethodPost, pathScheduledChannels(vhost, app), channel)
}

// GetScheduledChannels calls GET .../apps/{app}/scheduledChannels.
func (o *ovenMedia) GetScheduledChannels(ctx context.Context, vhost, app string) (*ResponseNameList, error) {
	return do[ResponseNameList](ctx, o, http.MethodGet, pathScheduledChannels(vhost, app), nil)
}

// GetScheduledChannel calls GET .../apps/{app}/scheduledChannels/{channel}.
func (o *ovenMedia) GetScheduledChannel(ctx context.Context, vhost, app, channel string) (*ResponseScheduledChannel, error) {
	return do[ResponseScheduledChannel](ctx, o, http.MethodGet, pathScheduledChannel(vhost, app, channel), nil)
}

// UpdateScheduledChannel calls PATCH .../scheduledChannels/{channel} to
// replace its fallback program and programs.
func (o *ovenMedia) UpdateScheduledChannel(ctx context.Context, vhost, app, channel string, patch ScheduledChannel) (*BaseResponseOK, error) {
	return do[BaseResponseOK](ctx, o, http.MethodPatch, pathScheduledChannel(vhost, app, channel), patch)
}

// DeleteScheduledChannel calls DELETE .../scheduledChannels/{channel}.
func (o *ovenMedia) DeleteScheduledChannel(ctx context.Context, vhost, app, channel string) (*BaseResponseOK, error) {
	return do[BaseResponseOK](ctx, o, http.MethodDelete, pathScheduledChannel(vhost, app, channel), nil)
}

// CreateMultiplexChannel calls POST .../apps/{app}/multiplexChannels.
func (o *ovenMedia) CreateMultiplexChannel(ctx context.Context, vhost, app string, channel MultiplexChannel) (*BaseResponseOK, error) {
	if err := channel.validate(); err != nil {
		return nil, err
	}
	return do[BaseResponseOK](ctx, o, http.MethodPost, pathMultiplexChannels(vhost, app), channel)
}

// GetMultiplexChannels calls GET .../apps/{app}/multiplexChannels.
func (o *ovenMedia) GetMultiplexChannels(ctx context.Context, vhost, app string) (*ResponseNameList, error) {
	return do[ResponseNameList](ctx, o, http.MethodGet, pathMultiplexChannels(vhost, app), nil)
}

// GetMultiplexChannel calls GET .../apps/{app}/multiplexChannels/{channel}.
func (o *ovenMedia) GetMultiplexChannel(ctx context.Context, vhost, app, channel string) (*ResponseMultiplexChannel, error) {
	return do[ResponseMultiplexChannel](ctx, o, http.MethodGet, pathMultiplexChannel(vhost, app, channel), nil)
}

// DeleteMultiplexChannel calls DELETE .../multiplexChannels/{channel}.
func (o *ovenMedia) DeleteMultiplexChannel(ctx context.Context, vhost, app, channel string) (*BaseResponseOK, error) {
	return do[BaseResponseOK](ctx, o, http.MethodDelete, pathMultiplexChannel(vhost, app, channel), nil)
}
