package ovenmedia

import (
	"encoding/json"
	"fmt"
)

// POST http://1.2.3.4:8081/v1/vhosts

func (o *ovenMedia) GetStreams(host string, application string) (*ResponseVirtualList, error) {
	url := fmt.Sprintf("%s/%s/apps/%s/streams", V1_HOSTS, host, application)
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

func (o *ovenMedia) GetStreamInfo(host string, application string, stream string) (*ResponseStreamInfo, error) {
	url := fmt.Sprintf("%s/%s/apps/%s/streams/%s", V1_HOSTS, host, application, stream)
	resp, err := o.get(url, nil)
	if err != nil {
		return nil, err
	}

	var obj ResponseStreamInfo
	if err := json.Unmarshal(resp.Body(), &obj); err != nil {
		return nil, err
	}
	return &obj, nil
}
