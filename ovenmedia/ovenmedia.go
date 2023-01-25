package ovenmedia

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
)

type OvenMedia struct {
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
	StartPush(vHost string, appName string, body interface{}) (*resty.Response, error)
	StopPush(vHost string, appName string, body interface{}) (*resty.Response, error)
	// Recording
}

func (o *OvenMedia) HealthCheck() error {
	resp, err := o.restClient.R().
		SetHeader("Accept", "application/json").
		Get(o.Url)
	//
	if err != nil {
		return err
	}
	//
	if !strings.Contains(resp.Status(), "200") {
		o.debugPrint(fmt.Sprintf("resp -> %v", resp))
		return errors.New("could not connect OvenMediaEngine")
	}
	return nil
}

// 
func (o *OvenMedia) IsDebug() bool {
	return o.debug
}

// Resty Methods

func (o *OvenMedia) restyPost(url string, body interface{}) (*resty.Response, error) {
	resp, err := o.restClient.R().
		SetHeader("Accept", "application/json").
		SetBody(body).
		Post(url)

	if err != nil {
		return nil, err
	}
	if !strings.Contains(resp.Status(), "200") {
		o.debugPrint(fmt.Sprintf("resp -> %v", resp))
		return nil, errors.New(resp.Status())
	}
	return resp, nil
}

func (o *OvenMedia) restyGet(url string, queryParams map[string]string) (*resty.Response, error) {
	resp, err := o.restClient.R().
		SetQueryParams(queryParams).
		Get(url)
	//
	if err != nil {
		return nil, err
	}
	if !strings.Contains(resp.Status(), "200") {
		o.debugPrint(fmt.Sprintf("resp -> %v", resp))
		return nil, errors.New(resp.Status())
	}
	return resp, nil
}
