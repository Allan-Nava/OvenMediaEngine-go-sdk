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
	Message    string   `json:"message"`
	Response   []string `json:"response"`
	StatusCode int      `json:"statusCode"`
}
//
// PUSH Stuff
type ResponseStartPush struct {
	Message    string       `json:"message"`
	Response   ResponsePush `json:"response"`
	StatusCode int          `json:"statusCode"`
}

type ResponsePush struct {
	App         string `json:"app"`
	CreatedTime string `json:"createdTime"`
	ID          string `json:"id"`
	Protocol    string `json:"protocol"`
	SentBytes   int    `json:"sentBytes"`
	SentTime    int    `json:"sentTime"`
	Sequence    int    `json:"sequence"`
	StartTime   string `json:"startTime"`
	State       string `json:"state"`
	Stream      struct {
		Name   string `json:"name"`
		Tracks []int  `json:"tracks"`
	} `json:"stream"`
	StreamKey      string `json:"streamKey"`
	TotalSentBytes int    `json:"totalsentBytes"`
	TotalSentTime  int    `json:"totalsentTime"`
	URL            string `json:"url"`
	Vhost          string `json:"vhost"`
}
//
type ResponsePushes struct {
	Message    string         `json:"message"`
	Response   []ResponsePush `json:"response"`
	StatusCode int            `json:"statusCode"`
}

/*
	{
		"createdTime": "2021-01-11T02:52:22.013+09:00",
		"lastRecvTime": "2021-01-11T04:11:41.734+09:00",
		"lastSentTime": "2021-01-11T02:52:22.013+09:00",
		"lastUpdatedTime": "2021-01-11T04:11:41.734+09:00",
		"maxTotalConnectionTime": "2021-01-11T02:52:22.013+09:00",
		"maxTotalConnections": 0,
		"totalBytesIn": 494713880,
		"totalBytesOut": 0,
		"totalConnections": 0
	}
*/
type ResponseStats struct {
	CreatedTime            string `json:"createdTime"`
	LastRecvTime           string `json:"lastRecvTime"`
	LastSentTime           string `json:"lastSentTime"`
	LastUpdatedTime        string `json:"lastUpdatedTime"`
	MaxTotalConnectionTime string `json:"maxTotalConnectionTime"`
	MaxTotalConnections    int    `json:"maxTotalConnections"`
	TotalBytesIn           int    `json:"totalBytesIn"`
	TotalBytesOut          int    `json:"totalBytesOut"`
	TotalConnections       int    `json:"totalConnections"`
}
