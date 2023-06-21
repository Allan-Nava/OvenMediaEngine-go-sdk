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

/* recording */

/*
	{
	  "id": "custom_id",
	  "stream": {
	    "name": "stream_o",
	    "tracks": [ 100, 200 ]
	  },
	  "filePath" : "/path/to/save/recorded/file_${Sequence}.ts",
	  "infoPath" : "/path/to/save/information/file.xml",
	  "interval" : 60000,    # Split it every 60 seconds
	  "schedule" : "0 0 1" # Split it at second 0, minute 0, every hours.
	  "segmentationRule" : "continuity"

}
*/
type RequestRecordingStart struct {
	ID     string `json:"id" required:"true" validate:"nonnil,min=1"`
	Stream struct {
		Name   string `json:"name"`
		Tracks []int  `json:"tracks"`
	} `json:"stream"`
	FilePath         string  `json:"filePath" required:"true" validate:"nonnil,min=1"`
	InfoPath         string  `json:"infoPath" required:"true" validate:"nonnil,min=1"`
	Interval         *int    `json:"interval,omitempty"`
	Schedule         *string `json:"schedule,omitempty"`
	SegmentationRule *string `json:"segmentationRule,omitempty"`
}

/*
	{
	  "id": "custom_id"
	}
*/
type RequestRecordingStop struct {
	ID string `json:"id" required:"true" validate:"nonnil,min=1"`
}
