package ovenmedia

type RequestCreateVirtualHost struct {
	VirtualHostsName []VRHostResponse
}

type RequestBodyPush struct {
	ID     string `json:"id" required:"true" validate:"nonnil,min=1"`
	Stream struct {
		Name   string `json:"name"`
		Tracks []int  `json:"tracks"`
	} `json:"stream"`
	Protocol  string `json:"protocol" required:"true" validate:"nonnil,min=1"`
	URL       string `json:"url" required:"true" validate:"nonnil,min=1"`
	StreamKey string `json:"streamKey" required:"true" validate:"nonnil,min=1"`
}
