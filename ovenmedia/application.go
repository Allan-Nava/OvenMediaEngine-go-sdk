package ovenmedia

import (
	"encoding/json"
	"fmt"
)

// POST http://1.2.3.4:8081/v1/vhosts

func (o *ovenMedia) GetApplications(host string) (*ResponseVirtualList, error) {
	url := fmt.Sprintf("%s/%s/apps", V1_HOSTS, host)
	resp, err := o.get(url, nil)
	if err != nil {
		return nil, err
	}

	var obj ResponseVirtualList
	if err := json.Unmarshal(resp.Body(), &obj); err != nil {
		return nil, err
	}
	return &obj, nil
}
