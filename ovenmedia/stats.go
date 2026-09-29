package ovenmedia

import (
	"context"
	"net/http"
)

// GetStatsVhosts calls GET /v1/stats/current/vhosts/{vhost}.
func (o *ovenMedia) GetStatsVhosts(ctx context.Context, vhost string) (*ResponseStats, error) {
	return do[ResponseStats](ctx, o, http.MethodGet, pathStats(vhost), nil)
}

// GetStatsAppVhosts calls GET /v1/stats/current/vhosts/{vhost}/apps/{app}.
func (o *ovenMedia) GetStatsAppVhosts(ctx context.Context, vhost, app string) (*ResponseStats, error) {
	return do[ResponseStats](ctx, o, http.MethodGet, pathStatsApp(vhost, app), nil)
}

// GetStatsStreamVhosts calls
// GET /v1/stats/current/vhosts/{vhost}/apps/{app}/streams/{stream}.
func (o *ovenMedia) GetStatsStreamVhosts(ctx context.Context, vhost, app, stream string) (*ResponseStats, error) {
	return do[ResponseStats](ctx, o, http.MethodGet, pathStatsStream(vhost, app, stream), nil)
}
