package ovenmedia

import (
	"fmt"
)

// POST http://1.2.3.4:8081/v1/vhosts/default/apps/app:startPush
func (o *OvenMedia) StartPush(vHost string, appName string)  (*resty.Response, error) {
	//
	resp, err := o.RestyPost(GET_VHOSTS_PUSH_BY_NAME(vHost, appName), nil)
	if err != nil {
		return nil, err
	}
	fmt.Println("resp", resp)
	return resp, nil
}
