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
