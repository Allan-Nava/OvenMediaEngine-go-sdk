package ovenmedia

import (
	"fmt"
	"errors"

	"github.com/go-resty/resty/v2"
)

// POST http://1.2.3.4:8081/v1/vhosts/default/apps/app:startPush
func (o *OvenMedia) StartPush(vHost string, appName string, body interface{}) (*ResponseStartPush, error) {
	//
    //
	resp, err := o.restyPost(GET_VHOSTS_PUSH_BY_NAME(vHost, appName), body)
	if err != nil {
		return nil, err
	}
	responseStartPush := resp.Result().(*ResponseStartPush)
	if responseStartPush == nil {
		return nil, errors.New("could not start the push")
	}
	return responseStartPush, nil
}


func (o *OvenMedia) StopPush(vHost string, appName string, body interface{}) (*resty.Response, error){
	//
	resp, err := o.restyPost(GET_VHOSTS_STOP_BY_NAME(vHost, appName), body)
	if err != nil {
		return nil, err
	}
	fmt.Println("resp", resp)
	return resp, nil
}
