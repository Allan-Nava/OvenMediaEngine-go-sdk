package ovenmedia

import (
	"context"
	"net/http"
)

// CreateVirtualHost calls POST /v1/vhosts. OME takes and returns an array,
// one entry per virtual host.
func (o *ovenMedia) CreateVirtualHost(ctx context.Context, vhosts ...VirtualHostConfig) ([]ResponseVirtualHost, error) {
	if len(vhosts) == 0 {
		return nil, invalid("vhost: at least one virtual host is required")
	}
	for _, v := range vhosts {
		if v.Name == "" {
			return nil, invalid("vhost: name is required")
		}
	}
	res, err := do[[]ResponseVirtualHost](ctx, o, http.MethodPost, pathVhosts, vhosts)
	if err != nil {
		return nil, err
	}
	return *res, nil
}

// GetAllVirtualHosts calls GET /v1/vhosts.
func (o *ovenMedia) GetAllVirtualHosts(ctx context.Context) (*ResponseNameList, error) {
	return do[ResponseNameList](ctx, o, http.MethodGet, pathVhosts, nil)
}

// GetVirtualHost calls GET /v1/vhosts/{vhost}.
func (o *ovenMedia) GetVirtualHost(ctx context.Context, vhost string) (*ResponseVirtualHost, error) {
	return do[ResponseVirtualHost](ctx, o, http.MethodGet, pathVhost(vhost), nil)
}

// DeleteVirtualHost calls DELETE /v1/vhosts/{vhost}.
func (o *ovenMedia) DeleteVirtualHost(ctx context.Context, vhost string) (*BaseResponseOK, error) {
	return do[BaseResponseOK](ctx, o, http.MethodDelete, pathVhost(vhost), nil)
}
