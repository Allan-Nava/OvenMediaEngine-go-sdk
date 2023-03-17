package ovenmedia

import (
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
)

type ovenMedia struct {
	Url        string
	debug      bool
	restClient *resty.Client
}

type IOvenMediaClient interface {
	IsDebug() bool
	HealthCheck() error
	// VirtualHost
	CreateVirtualHost(name string) (*ResponseVirtualHost, error)
	GetAllVirtualHosts() (*ResponseVirtualList, error)
	// Push
	StartPush(vHost string, appName string, body RequestBodyPush) (*ResponseStartPush, error)
	StopPush(vHost string, appName string, body RequestBodyPush) (*resty.Response, error)
	GetAllPushes(vHost string, appName string) (*ResponsePushes, error)
	// Recording
	StartRecording(vHost string, appName string, body RequestRecordingStart) (*ResponseRecordingStart, error)
	StopRecording(vHost string, appName string, body RequestRecordingStop) (*ResponseRecordingStart, error)
	// Stats
	GetStatsVhosts(vHost string) (*ResponseStats, error)
	GetStatsAppVhosts(vHost string, appName string) (*ResponseStats, error)
	GetStatsStreamVhosts(vHost string, appName string, stream string) (*ResponseStats, error)
	// Thumbnail
	GetThumbnail(appName string) (*resty.Response, error)
	//
}

func (o *ovenMedia) HealthCheck() error {
	_, err := o.restyGet(o.Url, nil)
	if err != nil {
		return err
	}
	return nil
}

func (o *ovenMedia) IsDebug() bool {
	return o.debug
}

// Resty Methods

func (o *ovenMedia) restyPost(url string, body interface{}) (*resty.Response, error) {
	resp, err := o.restClient.R().
		SetHeader("Accept", "application/json").
		SetBody(body).
		Post(url)

	if err != nil {
		return nil, err
	}
	if !strings.Contains(resp.Status(), "200") {
		err = fmt.Errorf("%v", resp)
		o.debugPrint(err)
		return nil, err
	}
	return resp, nil
}

func (o *ovenMedia) restyGet(url string, queryParams map[string]string) (*resty.Response, error) {
	resp, err := o.restClient.R().
		SetQueryParams(queryParams).
		Get(url)
	//
	if err != nil {
		return nil, err
	}
	if !strings.Contains(resp.Status(), "200") {
		err = fmt.Errorf("%v", resp)
		o.debugPrint(err)
		return nil, err
	}
	return resp, nil
}
