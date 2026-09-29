package ovenmedia_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Allan-Nava/OvenMediaEngine-go-sdk/ovenmedia"
)

// The examples run against a fake OME (fake_test.go) that answers with the
// replies from the OME REST API docs. apiURL stands in for your API server,
// e.g. "http://ome.example.com:8081".

var fake = startFake()

func apiURL() string { return fake.URL }

func newClient() ovenmedia.IOvenMediaClient {
	client, err := ovenmedia.New(apiURL(), ovenmedia.WithAccessToken("admin:secret"))
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func ExampleNew() {
	client, err := ovenmedia.New(apiURL(),
		ovenmedia.WithAccessToken("admin:secret"), // the AccessToken from Server.xml
		ovenmedia.WithTimeout(10*time.Second),
	)
	if err != nil {
		log.Fatal(err)
	}

	hosts, err := client.GetAllVirtualHosts(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(hosts.Response)
	// Output: [default]
}

func ExampleNew_healthCheck() {
	client := newClient()

	if err := client.HealthCheck(context.Background()); err != nil {
		log.Fatal("OME is down or rejects the token: ", err)
	}
	v, _ := client.GetVersion(context.Background())
	fmt.Println("OME", v.Response.Version)
	// Output: OME 0.18.0
}

func ExampleWithHeaders() {
	headers := ovenmedia.InitHeaderConfigurator()
	headers.CreateBasicAuthHeader("admin", "secret")
	headers.SetHeader("X-Request-Source", "billing")

	client, err := ovenmedia.New(apiURL(), ovenmedia.WithHeaders(headers))
	if err != nil {
		log.Fatal(err)
	}
	apps, _ := client.GetApplications(context.Background(), "default")
	fmt.Println(apps.Response)
	// Output: [app]
}

func ExampleAPIError() {
	client, _ := ovenmedia.New(apiURL(), ovenmedia.WithAccessToken("wrong"))

	_, err := client.GetAllVirtualHosts(context.Background())
	var apiErr *ovenmedia.APIError
	if errors.As(err, &apiErr) {
		fmt.Println(apiErr.StatusCode)
	}
	// Output: 401
}

func ExampleVirtualHostConfig() {
	client := newClient()

	created, err := client.CreateVirtualHost(context.Background(), ovenmedia.VirtualHostConfig{
		Name: "vhost",
		Host: &ovenmedia.VirtualHostHost{Names: []string{"ome.example.com"}},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(created[0].Response.Name, created[0].Response.Host.Names)
	// Output: vhost [ome.example.com]
}

func ExampleApplicationConfig() {
	client := newClient()

	created, err := client.CreateApplication(context.Background(), "default", ovenmedia.ApplicationConfig{
		Name:       "app2",
		Type:       ovenmedia.LIVE,
		Providers:  json.RawMessage(`{"rtmp":{},"srt":{}}`),
		Publishers: json.RawMessage(`{"llhls":{},"webrtc":{}}`),
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(created[0].Response.Name, created[0].Response.Type)
	// Output: app2 live
}

func ExampleOutputProfile() {
	client := newClient()

	_, err := client.CreateOutputProfile(context.Background(), "default", "app", ovenmedia.OutputProfile{
		Name:             "bypass",
		OutputStreamName: "${OriginStreamName}",
		Encodes:          json.RawMessage(`{"videos":[{"bypass":true}],"audios":[{"bypass":true}]}`),
	})
	if err != nil {
		log.Fatal(err)
	}
	profiles, _ := client.GetOutputProfiles(context.Background(), "default", "app")
	fmt.Println(profiles.Response)
	// Output: [bypass]
}

func ExampleRequestCreateStream() {
	client := newClient()

	res, err := client.CreateStream(context.Background(), "default", "app", ovenmedia.RequestCreateStream{
		Name: "stream",
		URLs: []string{"rtsp://camera.example.com/stream"},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.StatusCode, res.Message)
	// Output: 201 Created
}

func ExampleResponseStreamInfo() {
	client := newClient()

	info, err := client.GetStreamInfo(context.Background(), "default", "app", "stream")
	if err != nil {
		log.Fatal(err)
	}
	for _, t := range info.Response.Input.Tracks {
		if t.Video != nil {
			fmt.Printf("%s %dx%d %d bps\n", t.Video.Codec, t.Video.Width, t.Video.Height, t.Video.Bitrate)
		}
	}
	// Output: H264 1280x720 2500000 bps
}

func ExampleRequestSendEvent() {
	client := newClient()

	_, err := client.SendEvent(context.Background(), "default", "app", "stream", ovenmedia.RequestSendEvent{
		EventType: "video",
		Events: []ovenmedia.ID3Event{
			{FrameType: "TXXX", Info: "chapter", Data: "2"},
		},
	})
	fmt.Println(err)
	// Output: <nil>
}

func ExampleRequestBodyPush() {
	client := newClient()
	ctx := context.Background()

	res, err := client.StartPush(ctx, "default", "app", ovenmedia.RequestBodyPush{
		ID:        "push1",
		Stream:    ovenmedia.SimpleStream{Name: "stream"},
		Protocol:  "rtmp",
		URL:       "rtmp://ingest.example.com/live2",
		StreamKey: "KEY",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(res.Response.ID, res.Response.State)

	pushes, _ := client.GetAllPushes(ctx, "default", "app")
	fmt.Println(len(pushes.Response), "push")

	_, err = client.StopPush(ctx, "default", "app", "push1")
	fmt.Println(err)
	// Output:
	// push1 ready
	// 1 push
	// <nil>
}

func ExampleRequestRecordingStart() {
	client := newClient()
	ctx := context.Background()

	interval := 60000 // split every minute
	_, err := client.StartRecording(ctx, "default", "app", ovenmedia.RequestRecordingStart{
		ID:       "rec1",
		Stream:   ovenmedia.SimpleStream{Name: "stream"},
		Interval: &interval,
	})
	if err != nil {
		log.Fatal(err)
	}

	state, _ := client.GetRecordingState(ctx, "default", "app", "rec1")
	fmt.Println(state.Response[0].ID, state.Response[0].State)
	// Output: rec1 recording
}

func ExampleResponseStats() {
	client := newClient()

	stats, err := client.GetStatsStreamVhosts(context.Background(), "default", "app", "stream")
	if err != nil {
		log.Fatal(err)
	}
	s := stats.Response
	fmt.Println(s.TotalConnections, "viewers,", s.Connections["webrtc"], "on WebRTC")
	// Output: 43 viewers, 30 on WebRTC
}

func ExampleThumbnailFormat() {
	client := newClient()

	// Thumbnails come from the publisher port, not the API port.
	publisherURL := apiURL()
	img, err := client.GetThumbnail(context.Background(), publisherURL, "app", "stream", ovenmedia.ThumbnailPNG)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%q\n", img[1:4])
	// Output: "PNG"
}

func ExampleTime() {
	var rec ovenmedia.ResponseRecording
	// Some OME versions print the offset without a colon.
	_ = json.Unmarshal([]byte(`{"createdTime":"2021-08-31T23:44:44.789+0900"}`), &rec)
	fmt.Println(rec.CreatedTime.UTC().Format(time.RFC3339))
	// Output: 2021-08-31T14:44:44Z
}

func ExampleScheduledChannel() {
	client := newClient()
	ctx := context.Background()

	repeat := true
	_, err := client.CreateScheduledChannel(ctx, "default", "app", ovenmedia.ScheduledChannel{
		Stream: &ovenmedia.ScheduledChannelStream{Name: "channel"},
		FallbackProgram: &ovenmedia.ScheduledProgram{Items: []ovenmedia.ScheduledItem{
			{URL: "file://video/slate.mp4", Duration: -1}, // -1: to the end of the file
		}},
		Programs: []ovenmedia.ScheduledProgram{{
			Name:      "morning",
			Scheduled: "2023-11-20T20:57:00.000+09",
			Repeat:    &repeat,
			Items:     []ovenmedia.ScheduledItem{{URL: "file://video/1.mp4", Duration: 60000}},
		}},
	})
	if err != nil {
		log.Fatal(err)
	}

	info, _ := client.GetScheduledChannel(ctx, "default", "app", "channel")
	now := info.Response.CurrentProgram
	fmt.Println(now.State, now.CurrentItem.URL, now.CurrentItem.CurrentPosition, "ms")
	// Output: onair file://video/1.mp4 1700 ms
}

func ExampleMultiplexChannel() {
	client := newClient()
	ctx := context.Background()

	// One ABR output from an encoder feed: rename its tracks, declare
	// their bitrates, and list them in an LL-HLS playlist.
	_, err := client.CreateMultiplexChannel(ctx, "default", "app", ovenmedia.MultiplexChannel{
		OutputStream: ovenmedia.MultiplexOutputStream{Name: "abr"},
		SourceStreams: []ovenmedia.MultiplexSourceStream{{
			Name: "input1",
			URL:  "stream://default/app/input1",
			TrackMap: []ovenmedia.MultiplexTrackMap{
				{SourceTrackName: "bypass_video", NewTrackName: "input1_video", BitrateConf: 5000000, FramerateConf: 30},
				{SourceTrackName: "bypass_audio", NewTrackName: "input1_audio", BitrateConf: 128000},
			},
		}},
		Playlists: []ovenmedia.MultiplexPlaylist{{
			Name:       "LLHLS ABR",
			FileName:   "abr",
			Renditions: []ovenmedia.MultiplexRendition{{Name: "input1", Video: "input1_video", Audio: "input1_audio"}},
		}},
	})
	if err != nil {
		log.Fatal(err)
	}

	info, _ := client.GetMultiplexChannel(ctx, "default", "app", "abr")
	fmt.Println(info.Response.State)
	// Output: Pulling
}

func ExampleRequestHlsDump() {
	client := newClient()
	ctx := context.Background()

	_, err := client.StartHlsDump(ctx, "default", "app", "stream", ovenmedia.RequestHlsDump{
		ID:               "dump1",
		OutputStreamName: "stream",
		OutputPath:       "/var/dumps/",
		Playlist:         []string{"llhls.m3u8"},
	})
	if err != nil {
		log.Fatal(err)
	}

	// An empty ID stops every dump of the output stream.
	_, err = client.StopHlsDump(ctx, "default", "app", "stream", ovenmedia.RequestHlsDumpStop{OutputStreamName: "stream"})
	fmt.Println(err)
	// Output: <nil>
}
