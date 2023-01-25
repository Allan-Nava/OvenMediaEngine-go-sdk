package ovenmedia

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
)

type OvenMedia struct {
	Url        string
	RestClient *resty.Client
	Debug      bool
}

type IOvenMediaClient interface {
	HealthCheck() error
	// VirtualHost
	CreateVirtualHost(name string) error
	GetAllVirtualHosts() (*ResponseVirtualList, error)
	// Push
	StartPush(appName string) 
}

func (o *OvenMedia) HealthCheck() error {
	resp, err := o.RestClient.R().
		SetHeader("Accept", "application/json").
		Get(o.Url)
	//
	if err != nil {
		return err
	}
	//
	if !strings.Contains(resp.Status(), "200") {
		o.DebugPrint(fmt.Sprintf("resp -> %v", resp))
		return errors.New("could not connect OvenMediaEngine")
	}
	return nil
}

// SetHeader(header string, value string) *resty.Request
/*func (o *OvenMedia) SetHeaderAuthorization() *resty.Request {
	return o.RestClient.R().SetHeader("ome-access-token", "Basic ")// need to be finish
}*/