package ovenmedia

import (
	"errors"
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
	resp, err := o.RestyPost(V1_HOSTS, body)
	if err != nil {
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
	resp, err := o.RestyGet(V1_HOSTS, nil)
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
