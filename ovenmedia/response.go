package ovenmedia

// Response<VirtualHost>
type ResponseVirtualHost struct {
	VRHosts []VirtualHost
}

type VirtualHost struct {
	Message    string         `json:"message"`
	Response   VRHostResponse `json:"response"`
	StatusCode int            `json:"statusCode"`
}

type VRHostResponse struct {
	Name string `json:"name"`
}

type ResponseVirtualList struct {
	Message    string         	`json:"message"`
	Response   []string 		`json:"response"`
	StatusCode int            `json:"statusCode"`
}

/*
{
      "app": "app",
      "createdTime": "2021-02-04T01:42:40.160+0900",
      "id": "userDefinedUniqueId",
      "protocol": "rtmp",
      "sentBytes": 0,
      "sentTime": 0,
      "sequence": 0,
      "startTime": "1970-01-01T09:00:00.000+0900",
      "state": "ready",
      "stream": {
        "name": "stream",
        "tracks": [
          101,
          102
        ]
      },
      "streamKey": "",
      "totalsentBytes": 0,
      "totalsentTime": 0,
      "url": "",
      "vhost": "default"
    }
*/
type ResponsePush struct {
	App          string `json:"app"`
	CreatedTime  string `json:"createdTime"`
	ID           string `json:"id"`
	Protocol     string `json:"protocol"`
	SentBytes    int    `json:"sentBytes"`
	SentTime     int    `json:"sentTime"`
	Sequence     int    `json:"sequence"`
	StartTime    string `json:"startTime"`
	State        string `json:"state"`
	Stream       struct	{
		Name   string `json:"name"`
		Tracks []int  `json:"tracks"`
	} `json:"stream"`
	StreamKey    string `json:"streamKey"`
	TotalSentBytes int `json:"totalsentBytes"`
	TotalSentTime  int `json:"totalsentTime"`
	URL          string `json:"url"`
	Vhost        string `json:"vhost"`
}