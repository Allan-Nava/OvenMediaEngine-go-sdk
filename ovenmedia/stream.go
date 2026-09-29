package ovenmedia

import (
	"context"
	"net/http"
)

// CreateStream calls POST .../apps/{app}/streams to pull a stream from the
// given URLs.
func (o *ovenMedia) CreateStream(ctx context.Context, vhost, app string, body RequestCreateStream) (*BaseResponseOK, error) {
	if body.Name == "" {
		return nil, invalid("stream: name is required")
	}
	if len(body.URLs) == 0 {
		return nil, invalid("stream: at least one url is required")
	}
	return do[BaseResponseOK](ctx, o, http.MethodPost, pathStreams(vhost, app), body)
}

// GetStreams calls GET .../apps/{app}/streams.
func (o *ovenMedia) GetStreams(ctx context.Context, vhost, app string) (*ResponseNameList, error) {
	return do[ResponseNameList](ctx, o, http.MethodGet, pathStreams(vhost, app), nil)
}

// GetStreamInfo calls GET .../apps/{app}/streams/{stream}.
func (o *ovenMedia) GetStreamInfo(ctx context.Context, vhost, app, stream string) (*ResponseStreamInfo, error) {
	return do[ResponseStreamInfo](ctx, o, http.MethodGet, pathStream(vhost, app, stream), nil)
}

// DeleteStream calls DELETE .../apps/{app}/streams/{stream}.
func (o *ovenMedia) DeleteStream(ctx context.Context, vhost, app, stream string) (*BaseResponseOK, error) {
	return do[BaseResponseOK](ctx, o, http.MethodDelete, pathStream(vhost, app, stream), nil)
}

// SendEvent calls POST .../streams/{stream}:sendEvent to insert ID3v2
// metadata. OME answers 409 if the stream hasn't started yet.
func (o *ovenMedia) SendEvent(ctx context.Context, vhost, app, stream string, body RequestSendEvent) (*BaseResponseOK, error) {
	if body.EventFormat == "" {
		body.EventFormat = "id3v2"
	}
	if len(body.Events) == 0 {
		return nil, invalid("event: at least one event is required")
	}
	return do[BaseResponseOK](ctx, o, http.MethodPost, pathStream(vhost, app, stream)+":sendEvent", body)
}

// StartHlsDump calls POST .../streams/{stream}:startHlsDump. OME replies with
// the names of the streams being dumped.
func (o *ovenMedia) StartHlsDump(ctx context.Context, vhost, app, stream string, body RequestHlsDump) (*ResponseNameList, error) {
	if err := body.validate(); err != nil {
		return nil, err
	}
	return do[ResponseNameList](ctx, o, http.MethodPost, pathStream(vhost, app, stream)+":startHlsDump", body)
}

// StopHlsDump calls POST .../streams/{stream}:stopHlsDump.
func (o *ovenMedia) StopHlsDump(ctx context.Context, vhost, app, stream string, body RequestHlsDumpStop) (*ResponseNameList, error) {
	if body.OutputStreamName == "" {
		return nil, invalid("hls dump: outputStreamName is required")
	}
	return do[ResponseNameList](ctx, o, http.MethodPost, pathStream(vhost, app, stream)+":stopHlsDump", body)
}
