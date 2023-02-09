package ovenmedia

import (
	"fmt"
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"gopkg.in/validator.v2"
)

// POST http://1.2.3.4:8081/v1/vhosts/default/apps/app:startPush
// This is an action to request a push of a selected stream. Please refer to the "Push" document for detail setting.
func (o *ovenMedia) StartPush(vHost string, appName string, body RequestBodyPush) (*ResponseStartPush, error) {
	//
	if errs := validator.Validate(body); errs != nil {
		// values not valid, deal with errors here
		return nil, errs
	}
    //
	resp, err := o.restyPost(GET_VHOSTS_PUSH_BY_NAME(vHost, appName), body)
	if err != nil {
		return nil, err
	}
	var obj ResponseStartPush
	if err := json.Unmarshal(resp.Body(), &obj); err != nil {
		return nil, err
	}
	return &obj, nil
}

// Request to stop pushing
func (o *ovenMedia) StopPush(vHost string, appName string, body RequestBodyPush) (*resty.Response, error){
	//
	if errs := validator.Validate(body); errs != nil {
		// values not valid, deal with errors here
		return nil, errs
	}
	resp, err := o.restyPost(GET_VHOSTS_STOP_BY_NAME(vHost, appName), body)
	if err != nil {
		return nil, err
	}
	// TODO: 
	fmt.Println("resp", resp)
	return resp, nil
}


// Get all push lists for a specific application
func (o *ovenMedia) GetAllPushes(vHost string, appName string) (*ResponsePushes, error) {
	//
	resp, err := o.restyPost(GET_VHOSTS_PUSH_BY_NAME(vHost, appName), nil)
	if err != nil {
		return nil, err
	}
	var obj ResponsePushes
	if err := json.Unmarshal(resp.Body(), &obj); err != nil {
		return nil, err
	}
	return &obj, nil
}