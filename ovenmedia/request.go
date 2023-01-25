package ovenmedia

type RequestCreateVirtualHost struct {
	VirtualHostsName []VRHostResponse
}


/*{
    "id":"{2}",
    "stream":{
        "name":"",
        "tracks":[
            0,
            2
        ]
    },
    "protocol":"rtmp",
    "url":"",
    "streamKey":""
}*/
type RequestBodyPush struct {
	ID string `json:"id"`
	Stream struct {
		Name   string   `json:"name"`
		Tracks []int    `json:"tracks"`
	} `json:"stream"`
	Protocol  string `json:"protocol"`
	URL       string `json:"url"`
	StreamKey string `json:"streamKey"`
}