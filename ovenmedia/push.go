package ovenmedia

import (
	"fmt"

	"github.com/go-resty/resty/v2"
)

// POST http://1.2.3.4:8081/v1/vhosts/default/apps/app:startPush
func (o *OvenMedia) StartPush(vHost string, appName string, body interface{}) (*resty.Response, error) {
	//
    //
	resp, err := o.RestyPost(GET_VHOSTS_PUSH_BY_NAME(vHost, appName), body)
	if err != nil {
		return nil, err
	}
	fmt.Println("resp", resp)
	return resp, nil
}


func (o *OvenMedia) StopPush(vHost string, appName string, body interface{}) (*resty.Response, error){
	//
	resp, err := o.RestyPost(GET_VHOSTS_STOP_BY_NAME(vHost, appName), body)
	if err != nil {
		return nil, err
	}
	fmt.Println("resp", resp)
	return resp, nil
}
