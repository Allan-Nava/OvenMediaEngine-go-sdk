package ovenmedia


type ApplicationType string

const (
	LIVE      ApplicationType = "live"
	VOD       ApplicationType = "vod"
)

type CodecVideo string

const (
	H264    CodecVideo = "h264"
	H265    CodecVideo = "h265"
    VP8     CodecVideo = "vp8"
    OPUS    CodecVideo = "opus"
    AAC     CodecVideo = "aac"
)

type MediaType string

const (
    VIDEO MediaType = "video"
    AUDIO MediaType = "audio"
)

type SessionState string

const (
    READY       SessionState = "Ready"
    STARTED     SessionState = "Started"
    STOPPING    SessionState = "Stopping"
    STOPPED     SessionState = "Stopped"
    ERROR       SessionState = "Error"
)

type AudioLayout string


const (
    STEREO      AudioLayout = "stereo"
    MONO        AudioLayout = "mono"
)