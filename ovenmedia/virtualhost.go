package ovenmedia

import (
	"errors"
	"fmt"
	"strings"
)

// POST http://1.2.3.4:8081/v1/vhosts

func (o *OvenMedia) CreateVirtualHost(name string) error {
	body := &RequestCreateVirtualHost{
		VirtualHostsName: []VRHostResponse{
			{
				Name: name,
			},
		},
	}
	resp, err := RestyPost("/v1/hosts", body)
    if err != nil{
        return err
    }
	//
	virtualHost := resp.Result().(*ResponseVirtualHost)
	if virtualHost == nil {
		return errors.New("could not create virtual host")
	}
	return nil
}

// GET http://1.2.3.4:8081/v1/vhosts
func (o *OvenMedia) GetAllVirtualHosts() (*ResponseVirtualList, error) {
	//
	resp, err := RestyGet("/v1/hosts")
    if err != nil{
        return err
    }
	//
	virtualList := resp.Result().(*ResponseVirtualList)
	if virtualList == nil {
		return nil, errors.New("could not get all virtual host")
	}
	return virtualList, nil
}
