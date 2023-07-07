package ovenmedia

import (
	"github.com/go-resty/resty/v2"
)

// http(s)://<ome_host>:<port>/<app_name>/<output_stream_name>/thumb.<jpg|png>

func (o *ovenMedia) GetThumbnail(host, appName string, streamKey string) (*resty.Response, error) {
	resp, err := o.get(host+GET_THUMBNAIL(appName, streamKey), nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
