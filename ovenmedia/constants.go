package ovenmedia

import "fmt"

const (
	//virtual hosts
	V1_HOSTS      = "/v1/vhosts"
	V1_HOSTS_NAME = "/v1/vhosts/%s"
	// PUSH
	// /v1/vhosts/{vhost_name}/apps/{app_name}:startPush
	V1_HOSTS_START_PUSH_NAME = "/v1/vhosts/%s/apps/%s:startPush"
	//
	V1_HOSTS_STOP_PUSH_NAME = "/v1/vhosts/%s/apps/%s:stopPush"
	//http://<OME_HOST>:<API_PORT>/v1/vhosts/{vhost_name}/apps/{app_name}:pushes
	V1_HOSTS_PUSHES_NAME = "/v1/vhosts/%s/apps/%s:pushes"
)

var (
	//
	GET_VHOSTS_BY_NAME = func(vhostName string) string {
		// V1_HOSTS = "/v1/vhosts"+
		return fmt.Sprintf(V1_HOSTS_NAME, vhostName)
	}
	//http://<OME_HOST>:<API_PORT>/v1/vhosts/{vhost_name}/apps/{app_name}:startPush
	GET_VHOSTS_PUSH_BY_NAME = func(vhostName string, appName string) string {
		// V1_HOSTS = "/v1/vhosts"+
		return fmt.Sprintf(V1_HOSTS_START_PUSH_NAME, vhostName, appName)
	}
	//
	GET_VHOSTS_STOP_BY_NAME = func(vhostName string, appName string) string {
		// V1_HOSTS = "/v1/vhosts"+
		return fmt.Sprintf(V1_HOSTS_STOP_PUSH_NAME, vhostName, appName)
	}
	//
	GET_VHOSTS_PUSHES_BY_NAME = func(vhostName string, appName string) string {
		return fmt.Sprintf(V1_HOSTS_PUSHES_NAME, vhostName, appName)
	}
	//
)

type ApplicationType string

const (
	LIVE ApplicationType = "live"
	VOD  ApplicationType = "vod"
)

type CodecVideo string

const (
	H264 CodecVideo = "h264"
	H265 CodecVideo = "h265"
	VP8  CodecVideo = "vp8"
	OPUS CodecVideo = "opus"
	AAC  CodecVideo = "aac"
)

type MediaType string

const (
	VIDEO MediaType = "video"
	AUDIO MediaType = "audio"
)

type SessionState string

const (
	READY    SessionState = "Ready"
	STARTED  SessionState = "Started"
	STOPPING SessionState = "Stopping"
	STOPPED  SessionState = "Stopped"
	ERROR    SessionState = "Error"
)

type AudioLayout string

const (
	STEREO AudioLayout = "stereo"
	MONO   AudioLayout = "mono"
)
