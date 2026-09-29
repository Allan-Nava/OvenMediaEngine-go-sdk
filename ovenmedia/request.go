package ovenmedia

import (
	"encoding/json"
	"strings"
)

// VirtualHostConfig is one virtual host for [IOvenMediaClient.CreateVirtualHost].
type VirtualHostConfig struct {
	Name string           `json:"name"`
	Host *VirtualHostHost `json:"host,omitempty"`
}

type VirtualHostHost struct {
	Names []string        `json:"names,omitempty"`
	TLS   *VirtualHostTLS `json:"tls,omitempty"`
}

type VirtualHostTLS struct {
	CertPath      string `json:"certPath,omitempty"`
	ChainCertPath string `json:"chainCertPath,omitempty"`
	KeyPath       string `json:"keyPath,omitempty"`
}

// ApplicationConfig creates or patches an application. The nested sections
// follow Server.xml and change between OME versions, so they are raw JSON.
type ApplicationConfig struct {
	Name           string          `json:"name,omitempty"`
	Type           ApplicationType `json:"type,omitempty"`
	Dynamic        bool            `json:"dynamic,omitempty"`
	OutputProfiles json.RawMessage `json:"outputProfiles,omitempty"`
	Providers      json.RawMessage `json:"providers,omitempty"`
	Publishers     json.RawMessage `json:"publishers,omitempty"`
}

// OutputProfile is one transcoding profile. Encodes holds the "videos",
// "audios" and "images" lists as raw JSON.
type OutputProfile struct {
	Name             string          `json:"name"`
	OutputStreamName string          `json:"outputStreamName"`
	Encodes          json.RawMessage `json:"encodes,omitempty"`
}

// RequestCreateStream pulls a stream from one or more source URLs (RTSP,
// OVT, SRT...).
type RequestCreateStream struct {
	Name       string            `json:"name"`
	URLs       []string          `json:"urls"`
	Properties *StreamProperties `json:"properties,omitempty"`
}

type StreamProperties struct {
	Persistent                    *bool `json:"persistent,omitempty"`
	NoInputFailoverTimeoutMs      *int  `json:"noInputFailoverTimeoutMs,omitempty"`
	UnusedStreamDeletionTimeoutMs *int  `json:"unusedStreamDeletionTimeoutMs,omitempty"`
	IgnoreRtcpSRTimestamp         *bool `json:"ignoreRtcpSRTimestamp,omitempty"`
	Relay                         *bool `json:"relay,omitempty"`
}

// RequestSendEvent injects timed metadata into a stream. EventFormat
// defaults to "id3v2", the only format OME supports.
type RequestSendEvent struct {
	EventFormat string     `json:"eventFormat"`
	EventType   string     `json:"eventType,omitempty"`
	StartOffset int        `json:"startOffset,omitempty"`
	Events      []ID3Event `json:"events"`
}

// ID3Event is one ID3v2 frame. Info is only used by TXXX frames.
type ID3Event struct {
	FrameType string `json:"frameType"`
	Info      string `json:"info,omitempty"`
	Data      string `json:"data"`
}

// SimpleStream selects a stream, and optionally some of its renditions.
type SimpleStream struct {
	Name         string   `json:"name"`
	VariantNames []string `json:"variantNames,omitempty"`
	TrackIDs     []int    `json:"trackIds,omitempty"`
	// Tracks is the field name older OME versions used for TrackIDs.
	Tracks []int `json:"tracks,omitempty"`
}

// RequestBodyPush starts a push. StreamKey is required for RTMP only.
type RequestBodyPush struct {
	ID        string       `json:"id"`
	Stream    SimpleStream `json:"stream"`
	Protocol  string       `json:"protocol"`
	URL       string       `json:"url"`
	StreamKey string       `json:"streamKey"`
}

func (r RequestBodyPush) validate() error {
	switch {
	case r.ID == "":
		return invalid("push: id is required")
	case r.Stream.Name == "":
		return invalid("push: stream.name is required")
	case r.Protocol == "":
		return invalid("push: protocol is required")
	case r.URL == "":
		return invalid("push: url is required")
	case strings.EqualFold(r.Protocol, "rtmp") && r.StreamKey == "":
		return invalid("push: streamKey is required for rtmp")
	}
	return nil
}

// RequestRecordingStart starts a recording. Only ID and Stream.Name are
// required; OME falls back to the application's file settings.
//
//	{
//	  "id": "custom_id",
//	  "stream": {"name": "stream_o"},
//	  "filePath": "/path/to/save/recorded/file_${Sequence}.ts",
//	  "infoPath": "/path/to/save/information/file.xml",
//	  "interval": 60000,
//	  "segmentationRule": "continuity"
//	}
type RequestRecordingStart struct {
	ID               string       `json:"id"`
	Stream           SimpleStream `json:"stream"`
	FilePath         string       `json:"filePath,omitempty"`
	InfoPath         string       `json:"infoPath,omitempty"`
	Interval         *int         `json:"interval,omitempty"`
	Schedule         *string      `json:"schedule,omitempty"`
	SegmentationRule *string      `json:"segmentationRule,omitempty"`
}

func (r RequestRecordingStart) validate() error {
	switch {
	case r.ID == "":
		return invalid("recording: id is required")
	case r.Stream.Name == "":
		return invalid("recording: stream.name is required")
	}
	return nil
}

// RequestRecordingStop identifies a recording (or a push) by ID.
type RequestRecordingStop struct {
	ID string `json:"id"`
}
