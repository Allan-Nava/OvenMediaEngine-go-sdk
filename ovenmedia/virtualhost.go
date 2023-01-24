package ovenmedia

import (
	"errors"
	"fmt"
	"strings"
)

func (o *OvenMedia) CreateVirtualHost(name string) error {
	body := &RequestCreateVirtualHost{
		VirtualHostsName: []VRHostResponse{
			{
				Name: name,
			},
		},
	}
	resp, err := o.RestClient.R().
		SetHeader("Accept", "application/json").
		SetBody(body).
		SetResult(ResponseVirtualHost{}).
		Post(o.Url)
	//
	if err != nil {
		return err
	}
	if !strings.Contains(resp.Status(), "200") {
		o.DebugPrint(fmt.Sprintf("resp -> %v", resp))
		return errors.New(resp.Status())
	}
	//
	virtualHost := resp.Result().(*ResponseVirtualHost)
	if virtualHost == nil {
		return errors.New("could not create virtual host")
	}
	return nil
}
