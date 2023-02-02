package ovenmedia

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

/*
http://<OME_HOST>:<API_PORT>/v1/stats/current/vhosts/{vhost_name}/apps/{app_name}/streams/{stream}
*/