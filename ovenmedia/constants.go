package ovenmedia

import (
	"fmt"
	"net/url"
)

// API paths, relative to the base URL. Names are escaped so that a stream
// called "a/b" can't reach another resource.
const (
	pathVersion = "/v1/version"
	pathVhosts  = "/v1/vhosts"
)

func esc(s string) string { return url.PathEscape(s) }

func pathVhost(vhost string) string { return pathVhosts + "/" + esc(vhost) }

func pathApps(vhost string) string { return pathVhost(vhost) + "/apps" }

func pathApp(vhost, app string) string { return pathApps(vhost) + "/" + esc(app) }

// pathAppAction builds the ":verb" endpoints, e.g. apps/app:startPush.
func pathAppAction(vhost, app, action string) string { return pathApp(vhost, app) + ":" + action }

func pathOutputProfiles(vhost, app string) string { return pathApp(vhost, app) + "/outputProfiles" }

func pathOutputProfile(vhost, app, profile string) string {
	return pathOutputProfiles(vhost, app) + "/" + esc(profile)
}

func pathStreams(vhost, app string) string { return pathApp(vhost, app) + "/streams" }

func pathStream(vhost, app, stream string) string { return pathStreams(vhost, app) + "/" + esc(stream) }

func pathScheduledChannels(vhost, app string) string {
	return pathApp(vhost, app) + "/scheduledChannels"
}

func pathScheduledChannel(vhost, app, channel string) string {
	return pathScheduledChannels(vhost, app) + "/" + esc(channel)
}

func pathMultiplexChannels(vhost, app string) string {
	return pathApp(vhost, app) + "/multiplexChannels"
}

func pathMultiplexChannel(vhost, app, channel string) string {
	return pathMultiplexChannels(vhost, app) + "/" + esc(channel)
}

func pathStats(vhost string) string { return "/v1/stats/current/vhosts/" + esc(vhost) }

func pathStatsApp(vhost, app string) string { return pathStats(vhost) + "/apps/" + esc(app) }

func pathStatsStream(vhost, app, stream string) string {
	return pathStatsApp(vhost, app) + "/streams/" + esc(stream)
}

func pathThumbnail(app, stream string, format ThumbnailFormat) string {
	return fmt.Sprintf("/%s/%s/thumb.%s", esc(app), esc(stream), format)
}

// ApplicationType is an application's "type".
type ApplicationType string

const (
	LIVE ApplicationType = "live"
	VOD  ApplicationType = "vod"
)

// CodecVideo names a video codec in an output profile.
type CodecVideo string

const (
	H264 CodecVideo = "h264"
	H265 CodecVideo = "h265"
	VP8  CodecVideo = "vp8"
)

// CodecAudio names an audio codec in an output profile.
type CodecAudio string

const (
	OPUS CodecAudio = "opus"
	AAC  CodecAudio = "aac"
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

// ThumbnailFormat is the image type of [IOvenMediaClient.GetThumbnail].
type ThumbnailFormat string

const (
	ThumbnailPNG ThumbnailFormat = "png"
	ThumbnailJPG ThumbnailFormat = "jpg"
)
