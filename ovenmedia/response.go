package ovenmedia

import "encoding/json"

// BaseResponseOK is the envelope of every OME reply.
type BaseResponseOK struct {
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}

// ResponseNameList is a list of names: virtual hosts, applications, streams
// or output profiles.
type ResponseNameList struct {
	BaseResponseOK
	Response []string `json:"response"`
}

// ResponseVirtualList is the pre-0.5 name of [ResponseNameList].
//
// Deprecated: use ResponseNameList.
type ResponseVirtualList = ResponseNameList

type ResponseVersion struct {
	BaseResponseOK
	Response struct {
		Version    string `json:"version"`
		GitVersion string `json:"gitVersion"`
	} `json:"response"`
}

type ResponseVirtualHost struct {
	BaseResponseOK
	Response VirtualHostConfig `json:"response"`
}

type ResponseApplication struct {
	BaseResponseOK
	Response ApplicationConfig `json:"response"`
}

type ResponseOutputProfile struct {
	BaseResponseOK
	Response OutputProfile `json:"response"`
}

type ResponseStreamInfo struct {
	BaseResponseOK
	Response StreamInfo `json:"response"`
}

type StreamInfo struct {
	Name    string         `json:"name"`
	Input   StreamInput    `json:"input"`
	Outputs []StreamOutput `json:"outputs"`
}

type StreamInput struct {
	CreatedTime Time              `json:"createdTime"`
	SourceType  string            `json:"sourceType"`
	SourceURL   string            `json:"sourceUrl"`
	Connection  *StreamConnection `json:"connection,omitempty"`
	Tracks      []Track           `json:"tracks"`
}

type StreamConnection struct {
	Transport     string `json:"transport"`
	Protocol      string `json:"protocol"`
	LocalAddress  string `json:"localAddress"`
	LocalPort     int    `json:"localPort"`
	RemoteAddress string `json:"remoteAddress"`
	RemotePort    int    `json:"remotePort"`
}

type StreamOutput struct {
	Name      string          `json:"name"`
	Tracks    []Track         `json:"tracks"`
	Playlists json.RawMessage `json:"playlists,omitempty"`
}

type Track struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Type  string `json:"type"`
	Video *Video `json:"video,omitempty"`
	Audio *Audio `json:"audio,omitempty"`
}

type Audio struct {
	Bypass     bool      `json:"bypass"`
	Codec      string    `json:"codec"`
	Bitrate    FlexInt64 `json:"bitrate"`
	Samplerate int       `json:"samplerate"`
	Channel    int       `json:"channel"`
}

type Video struct {
	Bypass           bool      `json:"bypass"`
	Codec            string    `json:"codec"`
	Width            int       `json:"width"`
	Height           int       `json:"height"`
	Bitrate          FlexInt64 `json:"bitrate"`
	Framerate        float64   `json:"framerate"`
	KeyFrameInterval float64   `json:"keyFrameInterval"`
}

// ResponsePush is one push, as returned by :startPush and :pushes.
type ResponsePush struct {
	ID             string       `json:"id"`
	State          string       `json:"state"`
	Vhost          string       `json:"vhost"`
	App            string       `json:"app"`
	Stream         SimpleStream `json:"stream"`
	Protocol       string       `json:"protocol"`
	URL            string       `json:"url"`
	StreamKey      string       `json:"streamKey"`
	SentBytes      int64        `json:"sentBytes"`
	SentTime       int64        `json:"sentTime"`
	Sequence       int          `json:"sequence"`
	TotalSentBytes int64        `json:"totalsentBytes"`
	TotalSentTime  int64        `json:"totalsentTime"`
	CreatedTime    Time         `json:"createdTime"`
	StartTime      Time         `json:"startTime"`
	FinishTime     Time         `json:"finishTime"`
}

type ResponseStartPush struct {
	BaseResponseOK
	Response ResponsePush `json:"response"`
}

type ResponsePushes struct {
	BaseResponseOK
	Response []ResponsePush `json:"response"`
}

// ResponseStats wraps the current statistics of a virtual host,
// application or stream.
type ResponseStats struct {
	BaseResponseOK
	Response Stats `json:"response"`
}

type Stats struct {
	// Connections counts viewers per publisher: "webrtc", "llhls", "srt"...
	Connections            map[string]int `json:"connections"`
	CreatedTime            Time           `json:"createdTime"`
	LastRecvTime           Time           `json:"lastRecvTime"`
	LastSentTime           Time           `json:"lastSentTime"`
	LastUpdatedTime        Time           `json:"lastUpdatedTime"`
	MaxTotalConnectionTime Time           `json:"maxTotalConnectionTime"`
	MaxTotalConnections    int            `json:"maxTotalConnections"`
	TotalConnections       int            `json:"totalConnections"`
	TotalBytesIn           int64          `json:"totalBytesIn"`
	TotalBytesOut          int64          `json:"totalBytesOut"`
	LastThroughputIn       int64          `json:"lastThroughputIn"`
	LastThroughputOut      int64          `json:"lastThroughputOut"`
	AvgThroughputIn        int64          `json:"avgThroughputIn"`
	AvgThroughputOut       int64          `json:"avgThroughputOut"`
	MaxThroughputIn        int64          `json:"maxThroughputIn"`
	MaxThroughputOut       int64          `json:"maxThroughputOut"`
}

// ResponseRecordingStart is the reply of :startRecord and :stopRecord.
type ResponseRecordingStart struct {
	BaseResponseOK
	Response ResponseRecording `json:"response"`
}

type ResponseRecordingStateList struct {
	BaseResponseOK
	Response []ResponseRecording `json:"response"`
}

type ResponseRecording struct {
	ID               string       `json:"id"`
	State            string       `json:"state"`
	Vhost            string       `json:"vhost"`
	App              string       `json:"app"`
	Stream           SimpleStream `json:"stream"`
	FilePath         string       `json:"filePath"`
	InfoPath         string       `json:"infoPath"`
	Interval         int64        `json:"interval"`
	Schedule         string       `json:"schedule"`
	SegmentationRule string       `json:"segmentationRule"`
	CreatedTime      Time         `json:"createdTime"`
}
