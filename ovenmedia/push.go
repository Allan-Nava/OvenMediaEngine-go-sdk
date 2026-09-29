package ovenmedia

import (
	"context"
	"net/http"
)

// StartPush calls POST .../apps/{app}:startPush.
func (o *ovenMedia) StartPush(ctx context.Context, vhost, app string, body RequestBodyPush) (*ResponseStartPush, error) {
	if err := body.validate(); err != nil {
		return nil, err
	}
	return do[ResponseStartPush](ctx, o, http.MethodPost, pathAppAction(vhost, app, "startPush"), body)
}

// StopPush calls POST .../apps/{app}:stopPush for the push with this ID.
func (o *ovenMedia) StopPush(ctx context.Context, vhost, app, id string) (*BaseResponseOK, error) {
	if id == "" {
		return nil, invalid("push: id is required")
	}
	return do[BaseResponseOK](ctx, o, http.MethodPost, pathAppAction(vhost, app, "stopPush"), RequestRecordingStop{ID: id})
}

// GetAllPushes calls POST .../apps/{app}:pushes. With no pushes the result
// has an empty, non-nil Response.
func (o *ovenMedia) GetAllPushes(ctx context.Context, vhost, app string) (*ResponsePushes, error) {
	res, err := do[ResponsePushes](ctx, o, http.MethodPost, pathAppAction(vhost, app, "pushes"), nil)
	if err != nil {
		return nil, err
	}
	if res.Response == nil {
		res.Response = []ResponsePush{}
	}
	return res, nil
}
