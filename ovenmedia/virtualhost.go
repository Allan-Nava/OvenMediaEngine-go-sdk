package ovenmedia

import (
	"errors"
	
	"gopkg.in/validator.v2"
)

// POST http://1.2.3.4:8081/v1/vhosts

func (o *OvenMedia) CreateVirtualHost(name string) (*ResponseVirtualHost, error) {
	body := &RequestCreateVirtualHost{
		VirtualHostsName: []VRHostResponse{
			{
				Name: name,
			},
		},
	}
	if errs := validator.Validate(body); errs != nil {
		// values not valid, deal with errors here
		return nil, errs
	}
	resp, err := o.restyPost(V1_HOSTS, body)
	if err != nil {
		return nil, err
	}
	//
	virtualHost := resp.Result().(*ResponseVirtualHost)
	if virtualHost == nil {
		return nil, errors.New("could not create virtual host")
	}
	return virtualHost, nil
}

// GET http://1.2.3.4:8081/v1/vhosts
func (o *OvenMedia) GetAllVirtualHosts() (*ResponseVirtualList, error) {
	//
	resp, err := o.restyGet(V1_HOSTS, nil)
	if err != nil {
		return nil, err
	}
	//
	virtualList := resp.Result().(*ResponseVirtualList)
	if virtualList == nil {
		return nil, errors.New("could not get all virtual host")
	}
	return virtualList, nil
}
