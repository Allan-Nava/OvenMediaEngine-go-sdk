package ovenmedia

import (
	"fmt"
)

// POST http://1.2.3.4:8081/v1/vhosts/default/apps/app:startPush
func (o *OvenMedia) StartPush(appName string)  error{
    //
	resp, err := RestyPost(fmt.Sprintf("/v1/vhosts/default/apps/%s:startPush", appName), body)
    if err != nil{
        return err
    }
    fmt.Println("resp")
	return nil
}