// Package ovenmedia is a client for the OvenMediaEngine (OME) REST API v1.
//
// Create a client with [New], pointing it at the API server (the address in
// Server.xml under <Bind><Managers><API>, commonly port 8081) and passing the
// AccessToken:
//
//	client, err := ovenmedia.New("http://ome.example.com:8081",
//		ovenmedia.WithAccessToken("admin:secret"))
//
// Every method takes a [context.Context] and returns an [*APIError] when OME
// answers with a non-2xx status, so a 401 or 404 is never mistaken for success.
package ovenmedia

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

// IOvenMediaClient is the whole public surface of the SDK. Depend on it in
// your code and substitute a fake in tests.
type IOvenMediaClient interface {
	IsDebug() bool
	// HealthCheck calls GET /v1/version: it fails when OME is unreachable
	// and when it rejects the credentials.
	HealthCheck(ctx context.Context) error
	GetVersion(ctx context.Context) (*ResponseVersion, error)

	// VirtualHost
	CreateVirtualHost(ctx context.Context, vhosts ...VirtualHostConfig) ([]ResponseVirtualHost, error)
	GetAllVirtualHosts(ctx context.Context) (*ResponseNameList, error)
	GetVirtualHost(ctx context.Context, vhost string) (*ResponseVirtualHost, error)
	DeleteVirtualHost(ctx context.Context, vhost string) (*BaseResponseOK, error)

	// Application
	CreateApplication(ctx context.Context, vhost string, apps ...ApplicationConfig) ([]ResponseApplication, error)
	GetApplications(ctx context.Context, vhost string) (*ResponseNameList, error)
	GetApplication(ctx context.Context, vhost, app string) (*ResponseApplication, error)
	UpdateApplication(ctx context.Context, vhost, app string, patch ApplicationConfig) (*ResponseApplication, error)
	DeleteApplication(ctx context.Context, vhost, app string) (*BaseResponseOK, error)

	// Output profiles
	CreateOutputProfile(ctx context.Context, vhost, app string, profiles ...OutputProfile) ([]ResponseOutputProfile, error)
	GetOutputProfiles(ctx context.Context, vhost, app string) (*ResponseNameList, error)
	GetOutputProfile(ctx context.Context, vhost, app, profile string) (*ResponseOutputProfile, error)
	DeleteOutputProfile(ctx context.Context, vhost, app, profile string) (*BaseResponseOK, error)

	// Stream
	CreateStream(ctx context.Context, vhost, app string, body RequestCreateStream) (*BaseResponseOK, error)
	GetStreams(ctx context.Context, vhost, app string) (*ResponseNameList, error)
	GetStreamInfo(ctx context.Context, vhost, app, stream string) (*ResponseStreamInfo, error)
	DeleteStream(ctx context.Context, vhost, app, stream string) (*BaseResponseOK, error)
	SendEvent(ctx context.Context, vhost, app, stream string, body RequestSendEvent) (*BaseResponseOK, error)
	StartHlsDump(ctx context.Context, vhost, app, stream string, body RequestHlsDump) (*ResponseNameList, error)
	StopHlsDump(ctx context.Context, vhost, app, stream string, body RequestHlsDumpStop) (*ResponseNameList, error)

	// Scheduled channels
	CreateScheduledChannel(ctx context.Context, vhost, app string, channel ScheduledChannel) (*BaseResponseOK, error)
	GetScheduledChannels(ctx context.Context, vhost, app string) (*ResponseNameList, error)
	GetScheduledChannel(ctx context.Context, vhost, app, channel string) (*ResponseScheduledChannel, error)
	UpdateScheduledChannel(ctx context.Context, vhost, app, channel string, patch ScheduledChannel) (*BaseResponseOK, error)
	DeleteScheduledChannel(ctx context.Context, vhost, app, channel string) (*BaseResponseOK, error)

	// Multiplex channels
	CreateMultiplexChannel(ctx context.Context, vhost, app string, channel MultiplexChannel) (*BaseResponseOK, error)
	GetMultiplexChannels(ctx context.Context, vhost, app string) (*ResponseNameList, error)
	GetMultiplexChannel(ctx context.Context, vhost, app, channel string) (*ResponseMultiplexChannel, error)
	DeleteMultiplexChannel(ctx context.Context, vhost, app, channel string) (*BaseResponseOK, error)

	// Push
	StartPush(ctx context.Context, vhost, app string, body RequestBodyPush) (*ResponseStartPush, error)
	StopPush(ctx context.Context, vhost, app, id string) (*BaseResponseOK, error)
	GetAllPushes(ctx context.Context, vhost, app string) (*ResponsePushes, error)

	// Recording
	StartRecording(ctx context.Context, vhost, app string, body RequestRecordingStart) (*ResponseRecordingStart, error)
	StopRecording(ctx context.Context, vhost, app, id string) (*ResponseRecordingStart, error)
	ListRecordingState(ctx context.Context, vhost, app string) (*ResponseRecordingStateList, error)
	GetRecordingState(ctx context.Context, vhost, app, id string) (*ResponseRecordingStateList, error)

	// Stats
	GetStatsVhosts(ctx context.Context, vhost string) (*ResponseStats, error)
	GetStatsAppVhosts(ctx context.Context, vhost, app string) (*ResponseStats, error)
	GetStatsStreamVhosts(ctx context.Context, vhost, app, stream string) (*ResponseStats, error)

	// Thumbnail
	GetThumbnail(ctx context.Context, publisherURL, app, stream string, format ThumbnailFormat) ([]byte, error)
}

type ovenMedia struct {
	baseURL    string
	debug      bool
	restClient *resty.Client
}

func (o *ovenMedia) IsDebug() bool {
	return o.debug
}

func (o *ovenMedia) HealthCheck(ctx context.Context) error {
	_, err := o.GetVersion(ctx)
	return err
}

func (o *ovenMedia) GetVersion(ctx context.Context) (*ResponseVersion, error) {
	return do[ResponseVersion](ctx, o, http.MethodGet, pathVersion, nil)
}

// ErrInvalidRequest is wrapped by the errors returned when a request body is
// rejected before it is sent, for example a push without a URL.
var ErrInvalidRequest = errors.New("ovenmedia: invalid request")

func invalid(format string, args ...any) error {
	return fmt.Errorf("%w: %s", ErrInvalidRequest, fmt.Sprintf(format, args...))
}

// do sends one request and decodes the reply into T. A non-2xx status
// becomes an *APIError; an empty body (204) returns the zero T.
func do[T any](ctx context.Context, o *ovenMedia, method, path string, body any) (*T, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	req := o.restClient.R().
		SetContext(ctx).
		SetHeader("Accept", "application/json")
	if body != nil {
		req.SetHeader("Content-Type", "application/json").SetBody(body)
	}
	resp, err := req.Execute(method, path)
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, newAPIError(resp.StatusCode(), resp.Body())
	}
	var out T
	if len(resp.Body()) == 0 {
		return &out, nil
	}
	if err := json.Unmarshal(resp.Body(), &out); err != nil {
		return nil, fmt.Errorf("ovenmedia: decoding %s %s: %w", method, path, err)
	}
	return &out, nil
}
