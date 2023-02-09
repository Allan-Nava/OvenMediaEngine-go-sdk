package ovenmedia

import (
	"encoding/json"
	"gopkg.in/validator.v2"
)

// POST http://1.2.3.4:8081/v1/vhosts

func (o *ovenMedia) CreateVirtualHost(name string) (*ResponseVirtualHost, error) {
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
	var obj ResponseVirtualHost
	if err := json.Unmarshal(resp.Body(), &obj); err != nil {
		return nil, err
	}
	return &obj, nil
}

// GET http://1.2.3.4:8081/v1/vhosts
func (o *ovenMedia) GetAllVirtualHosts() (*ResponseVirtualList, error) {
	//
	resp, err := o.restyGet(V1_HOSTS, nil)
	if err != nil {
		return nil, err
	}
	var obj ResponseVirtualList
	if err := json.Unmarshal(resp.Body(), &obj); err != nil {
		return nil, err
	}
	//
	return &obj, nil
}
