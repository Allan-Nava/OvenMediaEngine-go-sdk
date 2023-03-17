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
	//
	// RECORDING
	// http://<OME_HOST>:<API_PORT>/v1/vhosts/{vhost_name}/apps/{app_name}:startRecord
	V1_HOSTS_START_RECORD_NAME = "/v1/vhosts/%s/apps/%s:startRecord"
	//http://<OME_HOST>:<API_PORT>/v1/vhosts/{vhost_name}/apps/{app_name}:stopRecord
	V1_HOSTS_STOP_RECORD_NAME = "/v1/vhosts/%s/apps/%s:stopRecord"
	// http://<OME_HOST>:<API_PORT>/v1/vhosts/{vhost_name}/apps/{app_name}:records
	V1_HOSTS_RECORDS_NAME = "/v1/vhosts/%s/apps/%s:records"
	//
	// STATS
	// http://<OME_HOST>:<API_PORT>/v1/stats/current/vhosts/{vhost_name}
	V1_CURRENT_STATS_NAME = "/v1/stats/current/vhosts/%s"
	// http://<OME_HOST>:<API_PORT>/v1/stats/current/vhosts/{vhost_name}/apps/{app_name}
	V1_CURRENT_STATS_APPP_NAME = "/v1/stats/current/vhosts/%s/apps/%s"
	// http://<OME_HOST>:<API_PORT>/v1/stats/current/vhosts/{vhost_name}/apps/{app_name}/stream/{stream}
	V1_CURRENT_STATS_APPP_STREAMS_NAME = "/v1/stats/current/vhosts/%s/apps/%s/stream/%s"
	//
)

var (
	//
	GET_VHOSTS_BY_NAME = func(vhostName string) string {
		return fmt.Sprintf(V1_HOSTS_NAME, vhostName)
	}
	// Get all vhost start push by name
	GET_VHOSTS_PUSH_BY_NAME = func(vhostName string, appName string) string {
		return fmt.Sprintf(V1_HOSTS_START_PUSH_NAME, vhostName, appName)
	}
	// Get all vhost stop by name
	GET_VHOSTS_STOP_BY_NAME = func(vhostName string, appName string) string {
		//
		return fmt.Sprintf(V1_HOSTS_STOP_PUSH_NAME, vhostName, appName)
	}
	//
	GET_VHOSTS_PUSHES_BY_NAME = func(vhostName string, appName string) string {
		return fmt.Sprintf(V1_HOSTS_PUSHES_NAME, vhostName, appName)
	}
	//
	GET_VHOSTS_START_RECORDED_BY_NAME = func(vhostName string, appName string) string {
		return fmt.Sprintf(V1_HOSTS_START_RECORD_NAME, vhostName, appName)
	}
	//
	GET_VHOSTS_STOP_RECORDED_BY_NAME = func(vhostName string, appName string) string {
		return fmt.Sprintf(V1_HOSTS_STOP_RECORD_NAME, vhostName, appName)
	}
	//
	GET_VHOSTS_RECORDS_BY_NAME = func(vhostName string, appName string) string {
		return fmt.Sprintf(V1_HOSTS_RECORDS_NAME, vhostName, appName)
	}
	//
	GET_CURRENT_STATS_NAME = func(vhostName string) string {
		return fmt.Sprintf(V1_CURRENT_STATS_NAME, vhostName)
	}
	//
	GET_CURRENT_STATS_APP_NAME = func(vhostName string, appName string) string {
		return fmt.Sprintf(V1_CURRENT_STATS_APPP_NAME, vhostName, appName)
	}
	//
	GET_CURRENT_STATS_STREAM = func(vhostName string, appName string, stream string) string {
		return fmt.Sprintf(V1_CURRENT_STATS_APPP_STREAMS_NAME, vhostName, appName, stream)
	}
	//
	GET_THUMBNAILS = func(appName string) string {
		return fmt.Sprintf("/%s/thumbnail.jpg", appName)
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
