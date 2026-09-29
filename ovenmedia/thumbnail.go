package ovenmedia

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// GetThumbnail downloads {publisherURL}/{app}/{stream}/thumb.{png|jpg}.
// Thumbnails are served by OME's publisher (the LL-HLS/WebRTC port), not by
// the API server, so publisherURL is required and the API's Authorization
// header is not sent to it. format defaults to PNG.
func (o *ovenMedia) GetThumbnail(ctx context.Context, publisherURL, app, stream string, format ThumbnailFormat) ([]byte, error) {
	u, err := url.Parse(publisherURL)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return nil, invalid("thumbnail: publisherURL must be an absolute http(s) URL, got %q", publisherURL)
	}
	if format == "" {
		format = ThumbnailPNG
	}
	if ctx == nil {
		ctx = context.Background()
	}
	target := strings.TrimRight(publisherURL, "/") + pathThumbnail(app, stream, format)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target, nil)
	if err != nil {
		return nil, err
	}
	resp, err := o.restClient.GetClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ovenmedia: reading thumbnail: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, newAPIError(resp.StatusCode, body)
	}
	return body, nil
}
