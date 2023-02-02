package ovenmedia


import "encoding/json"

/*
http://<OME_HOST>:<API_PORT>/v1/stats/current/vhosts/{vhost_name}
*/

func (o *OvenMedia) GetStatsVhosts(vHost string) (*ResponseStats, error) {
	//
	resp, err := o.restyGet(GET_CURRENT_STATS_NAME(vHost), nil)
	if err != nil {
		return nil, err
	}
	//
	var obj ResponseStats
	if err := json.Unmarshal(resp.Body(), &obj); err != nil {
		return nil, err
	}
	//
	return &obj, nil
}

/*
http://<OME_HOST>:<API_PORT>/v1/stats/current/vhosts/{vhost_name}/apps/{app_name}
*/

func (o *OvenMedia) GetStatsAppVhosts(vHost string, appName string) (*ResponseStats, error) {
	//
	resp, err := o.restyGet(GET_CURRENT_STATS_APP_NAME(vHost, appName), nil)
	if err != nil {
		return nil, err
	}
	//
	var obj ResponseStats
	if err := json.Unmarshal(resp.Body(), &obj); err != nil {
		return nil, err
	}
	//
	return &obj, nil
}


/*
http://<OME_HOST>:<API_PORT>/v1/stats/current/vhosts/{vhost_name}/apps/{app_name}/streams/{stream}
*/


func (o *OvenMedia) GetStatsStreamVhosts(vHost string, appName string, stream string) (*ResponseStats, error) {
	//
	resp, err := o.restyGet(GET_CURRENT_STATS_STREAM(vHost, appName, stream), nil)
	if err != nil {
		return nil, err
	}
	//
	var obj ResponseStats
	if err := json.Unmarshal(resp.Body(), &obj); err != nil {
		return nil, err
	}
	//
	return &obj, nil
}