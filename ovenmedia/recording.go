package ovenmedia

import "encoding/json"

//

func (o *OvenMedia) StartRecording(vHost string, appName string, body RequestRecordingStart) (*ResponseRecordingStart, error) {
	//
	resp, err := o.restyPost(GET_VHOSTS_START_RECORDED_BY_NAME(vHost, appName), body)
	if err != nil {
		return nil, err
	}
	//
	var obj ResponseRecordingStart
	if err := json.Unmarshal(resp.Body(), &obj); err != nil {
		return nil, err
	}
	//
	return &obj, nil
}

//

