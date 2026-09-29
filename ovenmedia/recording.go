package ovenmedia

import (
	"context"
	"net/http"
)

// StartRecording calls POST .../apps/{app}:startRecord.
func (o *ovenMedia) StartRecording(ctx context.Context, vhost, app string, body RequestRecordingStart) (*ResponseRecordingStart, error) {
	if err := body.validate(); err != nil {
		return nil, err
	}
	return do[ResponseRecordingStart](ctx, o, http.MethodPost, pathAppAction(vhost, app, "startRecord"), body)
}

// StopRecording calls POST .../apps/{app}:stopRecord for the recording with
// this ID.
func (o *ovenMedia) StopRecording(ctx context.Context, vhost, app, id string) (*ResponseRecordingStart, error) {
	if id == "" {
		return nil, invalid("recording: id is required")
	}
	return do[ResponseRecordingStart](ctx, o, http.MethodPost, pathAppAction(vhost, app, "stopRecord"), RequestRecordingStop{ID: id})
}

// ListRecordingState calls POST .../apps/{app}:records for every recording.
func (o *ovenMedia) ListRecordingState(ctx context.Context, vhost, app string) (*ResponseRecordingStateList, error) {
	return do[ResponseRecordingStateList](ctx, o, http.MethodPost, pathAppAction(vhost, app, "records"), nil)
}

// GetRecordingState calls POST .../apps/{app}:records filtered by ID. OME
// still replies with a list, of at most one recording.
func (o *ovenMedia) GetRecordingState(ctx context.Context, vhost, app, id string) (*ResponseRecordingStateList, error) {
	if id == "" {
		return nil, invalid("recording: id is required")
	}
	return do[ResponseRecordingStateList](ctx, o, http.MethodPost, pathAppAction(vhost, app, "records"), RequestRecordingStop{ID: id})
}
