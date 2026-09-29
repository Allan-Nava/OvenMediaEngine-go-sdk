package ovenmedia_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Allan-Nava/OvenMediaEngine-go-sdk/ovenmedia"
)

var ctx = context.Background()

func jsonEq(t *testing.T, got, want string) {
	t.Helper()
	var g, w any
	if err := json.Unmarshal([]byte(got), &g); err != nil {
		t.Fatalf("body %q is not JSON: %v", got, err)
	}
	if err := json.Unmarshal([]byte(want), &w); err != nil {
		t.Fatalf("want %q is not JSON: %v", want, err)
	}
	gb, _ := json.Marshal(g)
	wb, _ := json.Marshal(w)
	if string(gb) != string(wb) {
		t.Errorf("body = %s, want %s", gb, wb)
	}
}

func wantRequest(t *testing.T, f *fakeOME, method, path, body string) {
	t.Helper()
	r := f.last(t)
	if r.Method != method || r.Path != path {
		t.Errorf("request = %s %s, want %s %s", r.Method, r.Path, method, path)
	}
	if body == "" {
		if r.Body != "" {
			t.Errorf("body = %q, want none", r.Body)
		}
		return
	}
	jsonEq(t, r.Body, body)
}

func TestNewRejectsBadURL(t *testing.T) {
	for _, u := range []string{"", "ome.example.com:8081", "/v1", "ftp://ome.example.com"} {
		if _, err := ovenmedia.New(u); err == nil {
			t.Errorf("New(%q) = nil error", u)
		}
	}
}

func TestAccessTokenIsSentAsBasicAuth(t *testing.T) {
	f := newFake(t)
	if _, err := f.client(t).GetAllVirtualHosts(ctx); err != nil {
		t.Fatal(err)
	}
	if got := f.last(t).Auth; got != wantAuth {
		t.Errorf("Authorization = %q, want %q", got, wantAuth)
	}
}

func TestHeaderConfiguratorAuth(t *testing.T) {
	cases := map[string]func(h *ovenmedia.HeaderConfigurator){
		"CreateBasicAuthHeader":        func(h *ovenmedia.HeaderConfigurator) { h.CreateBasicAuthHeader("admin", "secret") },
		"CreateOmeBasicAuthHeaderWord": func(h *ovenmedia.HeaderConfigurator) { h.CreateOmeBasicAuthHeaderWord(token) },
		// deprecated: used to set an ome-access-token header OME ignores
		"CreateOmeBasicAuthHeader": func(h *ovenmedia.HeaderConfigurator) { h.CreateOmeBasicAuthHeader("admin", "secret") },
	}
	for name, set := range cases {
		t.Run(name, func(t *testing.T) {
			f := newFake(t)
			h := ovenmedia.InitHeaderConfigurator()
			set(h)
			c, err := ovenmedia.New(f.URL, ovenmedia.WithHeaders(h))
			if err != nil {
				t.Fatal(err)
			}
			if err := c.HealthCheck(ctx); err != nil {
				t.Fatalf("HealthCheck: %v", err)
			}
		})
	}
}

func TestBuildOvenStillWorks(t *testing.T) {
	f := newFake(t)
	h := ovenmedia.InitHeaderConfigurator()
	h.CreateBasicAuthHeader("admin", "secret")
	c, err := ovenmedia.BuildOven(f.URL, false, h)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := c.GetAllVirtualHosts(ctx); err != nil {
		t.Fatal(err)
	}
}

func TestHTTPErrorsAreReturned(t *testing.T) {
	f := newFake(t)
	c, err := ovenmedia.New(f.URL, ovenmedia.WithAccessToken("wrong"))
	if err != nil {
		t.Fatal(err)
	}
	_, err = c.GetAllVirtualHosts(ctx)
	var apiErr *ovenmedia.APIError
	if !errors.As(err, &apiErr) {
		t.Fatalf("err = %v, want *APIError", err)
	}
	if apiErr.StatusCode != http.StatusUnauthorized || !strings.Contains(apiErr.Message, "Authorization") {
		t.Errorf("APIError = %d %q", apiErr.StatusCode, apiErr.Message)
	}

	_, err = f.client(t).GetStreamInfo(ctx, "default", "app", "missing")
	if !errors.As(err, &apiErr) || apiErr.StatusCode != http.StatusNotFound {
		t.Errorf("missing stream: err = %v, want 404 APIError", err)
	}
}

func TestAPIErrorWithoutJSONBody(t *testing.T) {
	f := newFake(t)
	f.on("GET", "/v1/vhosts", reply{status: 502, body: "<html>bad gateway</html>"})
	_, err := f.client(t).GetAllVirtualHosts(ctx)
	var apiErr *ovenmedia.APIError
	if !errors.As(err, &apiErr) || apiErr.StatusCode != 502 || apiErr.Message != "Bad Gateway" {
		t.Fatalf("err = %v", err)
	}
}

func TestHealthCheck(t *testing.T) {
	f := newFake(t)
	if err := f.client(t).HealthCheck(ctx); err != nil {
		t.Fatal(err)
	}
	wantRequest(t, f, "GET", "/v1/version", "")

	bad, _ := ovenmedia.New(f.URL)
	if err := bad.HealthCheck(ctx); err == nil {
		t.Error("HealthCheck without credentials = nil, want 401")
	}

	down, _ := ovenmedia.New("http://127.0.0.1:1")
	if err := down.HealthCheck(ctx); err == nil {
		t.Error("HealthCheck on a closed port = nil")
	}
}

func TestContextAndTimeout(t *testing.T) {
	f := newFake(t)
	f.on("GET", "/v1/vhosts", reply{status: 200, body: `{"response":[]}`, delay: 2 * time.Second})

	cctx, cancel := context.WithTimeout(ctx, 50*time.Millisecond)
	defer cancel()
	if _, err := f.client(t).GetAllVirtualHosts(cctx); !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("with a context deadline: err = %v, want DeadlineExceeded", err)
	}

	start := time.Now()
	if _, err := f.client(t, ovenmedia.WithTimeout(50*time.Millisecond)).GetAllVirtualHosts(ctx); err == nil {
		t.Error("WithTimeout: err = nil")
	}
	if time.Since(start) > time.Second {
		t.Error("WithTimeout didn't bound the request")
	}
}

func TestWithHTTPClient(t *testing.T) {
	f := newFake(t)
	hc := &http.Client{Transport: f.Client().Transport}
	c, err := ovenmedia.New(f.URL, ovenmedia.WithHTTPClient(hc), ovenmedia.WithAccessToken(token))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := c.GetAllVirtualHosts(ctx); err != nil {
		t.Fatal(err)
	}
}

func TestDebugRedactsAuthorization(t *testing.T) {
	f := newFake(t)
	// resty's default logger writes to the os.Stderr it sees at New.
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	stderr := os.Stderr
	os.Stderr = w
	c := f.client(t, ovenmedia.WithDebug(true))
	os.Stderr = stderr
	if !c.IsDebug() {
		t.Error("IsDebug() = false")
	}
	if _, err := c.GetAllVirtualHosts(ctx); err != nil {
		t.Fatal(err)
	}
	_ = w.Close()
	logged, _ := io.ReadAll(r)
	if !strings.Contains(string(logged), "/v1/vhosts") {
		t.Fatalf("debug log doesn't show the request:\n%s", logged)
	}
	if strings.Contains(string(logged), strings.TrimPrefix(wantAuth, "Basic ")) {
		t.Errorf("debug log leaks the credentials:\n%s", logged)
	}
	if f.last(t).Auth != wantAuth {
		t.Error("redaction changed the header actually sent")
	}
}

func TestVirtualHosts(t *testing.T) {
	f := newFake(t)
	c := f.client(t)

	created, err := c.CreateVirtualHost(ctx, ovenmedia.VirtualHostConfig{
		Name: "vhost",
		Host: &ovenmedia.VirtualHostHost{Names: []string{"ome.example.com"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	wantRequest(t, f, "POST", "/v1/vhosts", `[{"name":"vhost","host":{"names":["ome.example.com"]}}]`)
	if len(created) != 1 || created[0].Response.Name != "vhost" || created[0].Response.Host.Names[0] != "ome.example.com" {
		t.Errorf("created = %+v", created)
	}

	if _, err := c.CreateVirtualHost(ctx); !errors.Is(err, ovenmedia.ErrInvalidRequest) {
		t.Errorf("no vhosts: err = %v", err)
	}

	list, err := c.GetAllVirtualHosts(ctx)
	if err != nil || len(list.Response) != 1 || list.Response[0] != "default" {
		t.Errorf("GetAllVirtualHosts = %+v, %v", list, err)
	}
	wantRequest(t, f, "GET", "/v1/vhosts", "")

	one, err := c.GetVirtualHost(ctx, "default")
	if err != nil || one.Response.Name != "default" {
		t.Errorf("GetVirtualHost = %+v, %v", one, err)
	}

	if _, err := c.DeleteVirtualHost(ctx, "vhost"); err != nil {
		t.Error(err)
	}
	wantRequest(t, f, "DELETE", "/v1/vhosts/vhost", "")
}

func TestApplicationsAndOutputProfiles(t *testing.T) {
	f := newFake(t)
	c := f.client(t)
	const app = "/v1/vhosts/default/apps/app"

	created, err := c.CreateApplication(ctx, "default", ovenmedia.ApplicationConfig{Name: "app2", Type: ovenmedia.LIVE})
	if err != nil || len(created) != 1 || created[0].Response.Name != "app2" {
		t.Fatalf("CreateApplication = %+v, %v", created, err)
	}
	wantRequest(t, f, "POST", "/v1/vhosts/default/apps", `[{"name":"app2","type":"live"}]`)

	apps, err := c.GetApplications(ctx, "default")
	if err != nil || apps.Response[0] != "app" {
		t.Errorf("GetApplications = %+v, %v", apps, err)
	}

	got, err := c.GetApplication(ctx, "default", "app")
	if err != nil || got.Response.Type != ovenmedia.LIVE {
		t.Errorf("GetApplication = %+v, %v", got, err)
	}
	jsonEq(t, string(got.Response.Providers), `{"rtmp":{}}`)

	patch := ovenmedia.ApplicationConfig{Providers: json.RawMessage(`{"webrtc":{"timeout":60000}}`)}
	if _, err := c.UpdateApplication(ctx, "default", "app", patch); err != nil {
		t.Error(err)
	}
	wantRequest(t, f, "PATCH", app, `{"providers":{"webrtc":{"timeout":60000}}}`)

	if _, err := c.DeleteApplication(ctx, "default", "app2"); err != nil {
		t.Error(err)
	}
	wantRequest(t, f, "DELETE", "/v1/vhosts/default/apps/app2", "")

	profiles, err := c.CreateOutputProfile(ctx, "default", "app", ovenmedia.OutputProfile{
		Name: "bypass", OutputStreamName: "${OriginStreamName}",
		Encodes: json.RawMessage(`{"videos":[{"bypass":true}]}`),
	})
	if err != nil || profiles[0].Response.Name != "bypass" {
		t.Errorf("CreateOutputProfile = %+v, %v", profiles, err)
	}
	wantRequest(t, f, "POST", app+"/outputProfiles",
		`[{"name":"bypass","outputStreamName":"${OriginStreamName}","encodes":{"videos":[{"bypass":true}]}}]`)

	if l, err := c.GetOutputProfiles(ctx, "default", "app"); err != nil || l.Response[0] != "bypass" {
		t.Errorf("GetOutputProfiles = %+v, %v", l, err)
	}
	if p, err := c.GetOutputProfile(ctx, "default", "app", "bypass"); err != nil || p.Response.OutputStreamName != "${OriginStreamName}" {
		t.Errorf("GetOutputProfile = %+v, %v", p, err)
	}
	if _, err := c.DeleteOutputProfile(ctx, "default", "app", "bypass"); err != nil {
		t.Error(err)
	}
	wantRequest(t, f, "DELETE", app+"/outputProfiles/bypass", "")
}

func TestStreams(t *testing.T) {
	f := newFake(t)
	c := f.client(t)
	const app = "/v1/vhosts/default/apps/app"

	persistent := true
	res, err := c.CreateStream(ctx, "default", "app", ovenmedia.RequestCreateStream{
		Name: "stream", URLs: []string{"rtsp://camera.example.com/stream"},
		Properties: &ovenmedia.StreamProperties{Persistent: &persistent},
	})
	if err != nil || res.StatusCode != 201 {
		t.Fatalf("CreateStream = %+v, %v", res, err)
	}
	wantRequest(t, f, "POST", app+"/streams",
		`{"name":"stream","urls":["rtsp://camera.example.com/stream"],"properties":{"persistent":true}}`)

	info, err := c.GetStreamInfo(ctx, "default", "app", "stream")
	if err != nil {
		t.Fatal(err)
	}
	in := info.Response.Input
	if in.SourceType != "Rtmp" || in.Connection.RemotePort != 41008 || in.CreatedTime.Year() != 2026 {
		t.Errorf("input = %+v", in)
	}
	if v := in.Tracks[0].Video; v == nil || v.Bitrate != 2500000 || v.Width != 1280 {
		t.Errorf("video = %+v", v)
	}
	// older OME versions send the bitrate as a string
	if a := in.Tracks[1].Audio; a == nil || a.Bitrate != 128000 {
		t.Errorf("audio = %+v", a)
	}

	if _, err := c.DeleteStream(ctx, "default", "app", "stream"); err != nil {
		t.Error(err)
	}
	wantRequest(t, f, "DELETE", app+"/streams/stream", "")

	if _, err := c.SendEvent(ctx, "default", "app", "stream", ovenmedia.RequestSendEvent{
		Events: []ovenmedia.ID3Event{{FrameType: "TXXX", Info: "chapter", Data: "2"}},
	}); err != nil {
		t.Error(err)
	}
	wantRequest(t, f, "POST", app+"/streams/stream:sendEvent",
		`{"eventFormat":"id3v2","events":[{"frameType":"TXXX","info":"chapter","data":"2"}]}`)
}

func TestNamesAreEscaped(t *testing.T) {
	f := newFake(t)
	_, _ = f.client(t).GetStreamInfo(ctx, "default", "app", "a/b")
	if got := f.last(t).Path; got != "/v1/vhosts/default/apps/app/streams/a%2Fb" {
		t.Errorf("path = %s", got)
	}
}

func TestPushes(t *testing.T) {
	f := newFake(t)
	c := f.client(t)
	const app = "/v1/vhosts/default/apps/app"

	res, err := c.StartPush(ctx, "default", "app", ovenmedia.RequestBodyPush{
		ID: "push1", Stream: ovenmedia.SimpleStream{Name: "stream"},
		Protocol: "rtmp", URL: "rtmp://ingest.example.com/live2", StreamKey: "KEY",
	})
	if err != nil {
		t.Fatal(err)
	}
	wantRequest(t, f, "POST", app+":startPush",
		`{"id":"push1","stream":{"name":"stream"},"protocol":"rtmp","url":"rtmp://ingest.example.com/live2","streamKey":"KEY"}`)
	if res.Response.State != "ready" || res.Response.CreatedTime.Year() != 2023 {
		t.Errorf("push = %+v", res.Response)
	}

	if _, err := c.StopPush(ctx, "default", "app", "push1"); err != nil {
		t.Error(err)
	}
	wantRequest(t, f, "POST", app+":stopPush", `{"id":"push1"}`)

	list, err := c.GetAllPushes(ctx, "default", "app")
	if err != nil || len(list.Response) != 1 || list.Response[0].StreamKey != "KEY" {
		t.Errorf("GetAllPushes = %+v, %v", list, err)
	}

	f.on("POST", app+":pushes", reply{status: http.StatusNoContent})
	empty, err := c.GetAllPushes(ctx, "default", "app")
	if err != nil || empty == nil || empty.Response == nil || len(empty.Response) != 0 {
		t.Errorf("no pushes: GetAllPushes = %+v, %v; want an empty list", empty, err)
	}
}

func TestPushValidation(t *testing.T) {
	f := newFake(t)
	c := f.client(t)
	srt := ovenmedia.RequestBodyPush{ID: "p", Stream: ovenmedia.SimpleStream{Name: "s"}, Protocol: "srt", URL: "srt://ingest.example.com:9999"}
	if _, err := c.StartPush(ctx, "default", "app", srt); err != nil {
		t.Errorf("SRT push without a stream key: %v", err)
	}

	before := f.count()
	for name, body := range map[string]ovenmedia.RequestBodyPush{
		"rtmp without key": {ID: "p", Stream: ovenmedia.SimpleStream{Name: "s"}, Protocol: "rtmp", URL: "rtmp://x/app"},
		"no id":            {Stream: ovenmedia.SimpleStream{Name: "s"}, Protocol: "srt", URL: "srt://x"},
		"no stream":        {ID: "p", Protocol: "srt", URL: "srt://x"},
		"no url":           {ID: "p", Stream: ovenmedia.SimpleStream{Name: "s"}, Protocol: "srt"},
	} {
		if _, err := c.StartPush(ctx, "default", "app", body); !errors.Is(err, ovenmedia.ErrInvalidRequest) {
			t.Errorf("%s: err = %v, want ErrInvalidRequest", name, err)
		}
	}
	if _, err := c.StopPush(ctx, "default", "app", ""); !errors.Is(err, ovenmedia.ErrInvalidRequest) {
		t.Errorf("StopPush without id: %v", err)
	}
	if f.count() != before {
		t.Error("an invalid request reached the server")
	}
}

func TestRecordings(t *testing.T) {
	f := newFake(t)
	c := f.client(t)
	const app = "/v1/vhosts/default/apps/app"

	started, err := c.StartRecording(ctx, "default", "app", ovenmedia.RequestRecordingStart{
		ID: "rec1", Stream: ovenmedia.SimpleStream{Name: "stream"},
	})
	if err != nil || started.Response.State != "ready" {
		t.Fatalf("StartRecording = %+v, %v", started, err)
	}
	// filePath and infoPath are optional: not sent as ""
	wantRequest(t, f, "POST", app+":startRecord", `{"id":"rec1","stream":{"name":"stream"}}`)

	if _, err := c.StartRecording(ctx, "default", "app", ovenmedia.RequestRecordingStart{ID: "x"}); !errors.Is(err, ovenmedia.ErrInvalidRequest) {
		t.Errorf("no stream: err = %v", err)
	}

	stopped, err := c.StopRecording(ctx, "default", "app", "rec1")
	if err != nil || stopped.Response.State != "stopping" {
		t.Errorf("StopRecording = %+v, %v", stopped, err)
	}
	wantRequest(t, f, "POST", app+":stopRecord", `{"id":"rec1"}`)

	all, err := c.ListRecordingState(ctx, "default", "app")
	if err != nil || len(all.Response) != 1 || all.Response[0].Interval != 60000 {
		t.Errorf("ListRecordingState = %+v, %v", all, err)
	}
	wantRequest(t, f, "POST", app+":records", "")

	one, err := c.GetRecordingState(ctx, "default", "app", "rec1")
	if err != nil || len(one.Response) != 1 || one.Response[0].ID != "rec1" {
		t.Errorf("GetRecordingState = %+v, %v", one, err)
	}
	wantRequest(t, f, "POST", app+":records", `{"id":"rec1"}`)
}

func TestStats(t *testing.T) {
	f := newFake(t)
	c := f.client(t)
	for path, call := range map[string]func() (*ovenmedia.ResponseStats, error){
		"/v1/stats/current/vhosts/default":          func() (*ovenmedia.ResponseStats, error) { return c.GetStatsVhosts(ctx, "default") },
		"/v1/stats/current/vhosts/default/apps/app": func() (*ovenmedia.ResponseStats, error) { return c.GetStatsAppVhosts(ctx, "default", "app") },
		"/v1/stats/current/vhosts/default/apps/app/streams/stream": func() (*ovenmedia.ResponseStats, error) {
			return c.GetStatsStreamVhosts(ctx, "default", "app", "stream")
		},
	} {
		res, err := call()
		if err != nil {
			t.Fatalf("%s: %v", path, err)
		}
		wantRequest(t, f, "GET", path, "")
		s := res.Response
		if s.TotalConnections != 43 || s.TotalBytesIn != 494713880 || s.Connections["webrtc"] != 30 || s.LastThroughputIn != 2500000 {
			t.Errorf("%s: stats = %+v", path, s)
		}
		if s.CreatedTime.IsZero() {
			t.Errorf("%s: createdTime not decoded", path)
		}
	}
}

func TestThumbnail(t *testing.T) {
	f := newFake(t)
	img, err := f.client(t).GetThumbnail(ctx, f.URL, "app", "stream", "")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(string(img), "\x89PNG") {
		t.Errorf("thumbnail = %q", img)
	}
	r := f.last(t)
	if r.Path != "/app/stream/thumb.png" {
		t.Errorf("path = %s", r.Path)
	}
	if r.Auth != "" {
		t.Error("the API credentials were sent to the publisher")
	}

	_, err = f.client(t).GetThumbnail(ctx, f.URL, "app", "stream", ovenmedia.ThumbnailJPG)
	var apiErr *ovenmedia.APIError
	if !errors.As(err, &apiErr) || apiErr.StatusCode != 404 {
		t.Errorf("missing jpg: err = %v, want 404", err)
	}
	if _, err := f.client(t).GetThumbnail(ctx, "", "app", "stream", ""); !errors.Is(err, ovenmedia.ErrInvalidRequest) {
		t.Errorf("no publisher URL: err = %v", err)
	}
}

func TestTimeFormats(t *testing.T) {
	for _, in := range []string{
		`"2021-08-31T23:44:44.789+09:00"`,
		`"2021-08-31T23:44:44.789+0900"`, // OME's recording example
	} {
		var got ovenmedia.Time
		if err := json.Unmarshal([]byte(in), &got); err != nil {
			t.Errorf("%s: %v", in, err)
			continue
		}
		if got.UTC().Format(time.RFC3339) != "2021-08-31T14:44:44Z" {
			t.Errorf("%s = %s", in, got.UTC())
		}
	}
	for _, in := range []string{`""`, `null`} {
		var got ovenmedia.Time
		if err := json.Unmarshal([]byte(in), &got); err != nil || !got.IsZero() {
			t.Errorf("%s = %v, %v; want zero", in, got, err)
		}
	}
	var bad ovenmedia.Time
	if err := json.Unmarshal([]byte(`"yesterday"`), &bad); err == nil {
		t.Error("garbage timestamp accepted")
	}
	b, _ := json.Marshal(ovenmedia.Time{Time: time.Date(2021, 8, 31, 14, 44, 44, 0, time.UTC)})
	if string(b) != `"2021-08-31T14:44:44Z"` {
		t.Errorf("marshal = %s", b)
	}
}

func TestFlexInt64(t *testing.T) {
	for in, want := range map[string]ovenmedia.FlexInt64{`2500000`: 2500000, `"128000"`: 128000, `""`: 0, `null`: 0} {
		var got ovenmedia.FlexInt64
		if err := json.Unmarshal([]byte(in), &got); err != nil || got != want {
			t.Errorf("%s = %d, %v; want %d", in, got, err, want)
		}
	}
	var bad ovenmedia.FlexInt64
	if err := json.Unmarshal([]byte(`"fast"`), &bad); err == nil {
		t.Error(`"fast" accepted as a number`)
	}
}

func TestScheduledChannels(t *testing.T) {
	f := newFake(t)
	c := f.client(t)
	const base = "/v1/vhosts/default/apps/app/scheduledChannels"

	video, yes, no := true, true, false
	channel := ovenmedia.ScheduledChannel{
		Stream: &ovenmedia.ScheduledChannelStream{Name: "channel", VideoTrack: &video},
		FallbackProgram: &ovenmedia.ScheduledProgram{Items: []ovenmedia.ScheduledItem{
			{URL: "file://video/sample.mp4", Start: 0, Duration: 60000},
		}},
		Programs: []ovenmedia.ScheduledProgram{{
			Name: "1", Scheduled: "2023-11-13T20:57:00.000+09", Repeat: &yes,
			Items: []ovenmedia.ScheduledItem{{URL: "file://video/1.mp4", Duration: 60000}},
		}},
	}
	res, err := c.CreateScheduledChannel(ctx, "default", "app", channel)
	if err != nil || res.StatusCode != 201 {
		t.Fatalf("CreateScheduledChannel = %+v, %v", res, err)
	}
	wantRequest(t, f, "POST", base, `{"stream":{"name":"channel","videoTrack":true},
		"fallbackProgram":{"items":[{"url":"file://video/sample.mp4","start":0,"duration":60000}]},
		"programs":[{"name":"1","scheduled":"2023-11-13T20:57:00.000+09","repeat":true,
		 "items":[{"url":"file://video/1.mp4","start":0,"duration":60000}]}]}`)

	if _, err := c.CreateScheduledChannel(ctx, "default", "app", ovenmedia.ScheduledChannel{}); !errors.Is(err, ovenmedia.ErrInvalidRequest) {
		t.Errorf("no stream: err = %v", err)
	}

	if l, err := c.GetScheduledChannels(ctx, "default", "app"); err != nil || l.Response[0] != "channel" {
		t.Errorf("GetScheduledChannels = %+v, %v", l, err)
	}

	info, err := c.GetScheduledChannel(ctx, "default", "app", "channel")
	if err != nil {
		t.Fatal(err)
	}
	cur := info.Response.CurrentProgram
	if cur == nil || cur.State != "onair" || cur.CurrentItem.CurrentPosition != 1700 || cur.Duration != -1 || cur.Scheduled.Year() != 2023 {
		t.Errorf("currentProgram = %+v", cur)
	}
	// the programs list uses an hour-only offset: "+09"
	if p := info.Response.Programs; len(p) != 1 || p[0].Scheduled.UTC().Hour() != 11 {
		t.Errorf("programs = %+v", p)
	}
	if info.Response.Stream.Name != "channel" || info.Response.FallbackProgram.Items[0].Duration != -1 {
		t.Errorf("channel = %+v", info.Response)
	}

	patch := ovenmedia.ScheduledChannel{Programs: []ovenmedia.ScheduledProgram{{
		Name: "2", Scheduled: "2023-11-20T20:57:00.000+09", Repeat: &no,
		Items: []ovenmedia.ScheduledItem{{URL: "file://video/1.mp4", Duration: -1}},
	}}}
	if _, err := c.UpdateScheduledChannel(ctx, "default", "app", "channel", patch); err != nil {
		t.Error(err)
	}
	// an explicit repeat: false is sent
	wantRequest(t, f, "PATCH", base+"/channel", `{"programs":[{"name":"2","scheduled":"2023-11-20T20:57:00.000+09","repeat":false,
		"items":[{"url":"file://video/1.mp4","start":0,"duration":-1}]}]}`)

	if _, err := c.DeleteScheduledChannel(ctx, "default", "app", "channel"); err != nil {
		t.Error(err)
	}
	wantRequest(t, f, "DELETE", base+"/channel", "")
}

func TestMultiplexChannels(t *testing.T) {
	f := newFake(t)
	c := f.client(t)
	const base = "/v1/vhosts/default/apps/app/multiplexChannels"

	abr := true
	res, err := c.CreateMultiplexChannel(ctx, "default", "app", ovenmedia.MultiplexChannel{
		OutputStream: ovenmedia.MultiplexOutputStream{Name: "abr"},
		SourceStreams: []ovenmedia.MultiplexSourceStream{{
			Name: "input1", URL: "stream://default/app/input1",
			TrackMap: []ovenmedia.MultiplexTrackMap{
				{SourceTrackName: "bypass_video", NewTrackName: "input1_video", BitrateConf: 5000000, FramerateConf: 30},
				{SourceTrackName: "bypass_audio", NewTrackName: "input1_audio", BitrateConf: 128000},
			},
		}},
		Playlists: []ovenmedia.MultiplexPlaylist{{
			Name: "LLHLS ABR", FileName: "abr",
			Options:    &ovenmedia.MultiplexPlaylistOptions{WebrtcAutoAbr: &abr},
			Renditions: []ovenmedia.MultiplexRendition{{Name: "input1", Video: "input1_video", Audio: "input1_audio"}},
		}},
	})
	if err != nil || res.StatusCode != 201 {
		t.Fatalf("CreateMultiplexChannel = %+v, %v", res, err)
	}
	wantRequest(t, f, "POST", base, `{"outputStream":{"name":"abr"},
		"sourceStreams":[{"name":"input1","url":"stream://default/app/input1","trackMap":[
		 {"sourceTrackName":"bypass_video","newTrackName":"input1_video","bitrateConf":5000000,"framerateConf":30},
		 {"sourceTrackName":"bypass_audio","newTrackName":"input1_audio","bitrateConf":128000}]}],
		"playlists":[{"name":"LLHLS ABR","fileName":"abr","options":{"webrtcAutoAbr":true},
		 "renditions":[{"name":"input1","video":"input1_video","audio":"input1_audio"}]}]}`)

	for name, bad := range map[string]ovenmedia.MultiplexChannel{
		"no output": {SourceStreams: []ovenmedia.MultiplexSourceStream{{Name: "a", URL: "stream://default/app/a"}}},
		"no source": {OutputStream: ovenmedia.MultiplexOutputStream{Name: "abr"}},
		"no url":    {OutputStream: ovenmedia.MultiplexOutputStream{Name: "abr"}, SourceStreams: []ovenmedia.MultiplexSourceStream{{Name: "a"}}},
	} {
		if _, err := c.CreateMultiplexChannel(ctx, "default", "app", bad); !errors.Is(err, ovenmedia.ErrInvalidRequest) {
			t.Errorf("%s: err = %v", name, err)
		}
	}

	if l, err := c.GetMultiplexChannels(ctx, "default", "app"); err != nil || l.Response[0] != "abr" {
		t.Errorf("GetMultiplexChannels = %+v, %v", l, err)
	}

	info, err := c.GetMultiplexChannel(ctx, "default", "app", "abr")
	if err != nil {
		t.Fatal(err)
	}
	m := info.Response
	if m.State != "Pulling" || !strings.Contains(m.PullingMessage, "input1") || m.OutputStream.Name != "abr" {
		t.Errorf("channel = %+v", m)
	}
	if tm := m.SourceStreams[0].TrackMap[1]; tm.BitrateConf != 5000000 || tm.FramerateConf != 30 {
		t.Errorf("trackMap = %+v", tm)
	}
	if p := m.Playlists[0]; p.Options == nil || p.Options.WebrtcAutoAbr == nil || !*p.Options.WebrtcAutoAbr || p.Renditions[0].Video != "input1_video" {
		t.Errorf("playlist = %+v", p)
	}

	if _, err := c.DeleteMultiplexChannel(ctx, "default", "app", "abr"); err != nil {
		t.Error(err)
	}
	wantRequest(t, f, "DELETE", base+"/abr", "")
}

func TestHlsDump(t *testing.T) {
	f := newFake(t)
	c := f.client(t)
	const stream = "/v1/vhosts/default/apps/app/streams/stream"

	res, err := c.StartHlsDump(ctx, "default", "app", "stream", ovenmedia.RequestHlsDump{
		ID: "dump1", OutputStreamName: "stream", OutputPath: "/tmp/dump/", Playlist: []string{"llhls.m3u8"},
	})
	if err != nil || res.Response[0] != "stream" {
		t.Fatalf("StartHlsDump = %+v, %v", res, err)
	}
	wantRequest(t, f, "POST", stream+":startHlsDump",
		`{"id":"dump1","outputStreamName":"stream","outputPath":"/tmp/dump/","playlist":["llhls.m3u8"]}`)

	before := f.count()
	if _, err := c.StartHlsDump(ctx, "default", "app", "stream", ovenmedia.RequestHlsDump{OutputStreamName: "stream"}); !errors.Is(err, ovenmedia.ErrInvalidRequest) {
		t.Errorf("no id: err = %v", err)
	}
	if _, err := c.StartHlsDump(ctx, "default", "app", "stream", ovenmedia.RequestHlsDump{ID: "d"}); !errors.Is(err, ovenmedia.ErrInvalidRequest) {
		t.Errorf("no output stream: err = %v", err)
	}
	if f.count() != before {
		t.Error("an invalid request reached the server")
	}

	if _, err := c.StopHlsDump(ctx, "default", "app", "stream", ovenmedia.RequestHlsDumpStop{OutputStreamName: "stream", ID: "dump1"}); err != nil {
		t.Error(err)
	}
	wantRequest(t, f, "POST", stream+":stopHlsDump", `{"outputStreamName":"stream","id":"dump1"}`)

	// without an ID OME stops every dump of the stream
	if _, err := c.StopHlsDump(ctx, "default", "app", "stream", ovenmedia.RequestHlsDumpStop{OutputStreamName: "stream"}); err != nil {
		t.Error(err)
	}
	wantRequest(t, f, "POST", stream+":stopHlsDump", `{"outputStreamName":"stream"}`)
}

func TestTimeHourOnlyOffset(t *testing.T) {
	var got ovenmedia.Time
	if err := json.Unmarshal([]byte(`"2023-11-13T20:57:00.000+09"`), &got); err != nil {
		t.Fatal(err)
	}
	if got.UTC().Format(time.RFC3339) != "2023-11-13T11:57:00Z" {
		t.Errorf("got %s", got.UTC())
	}
}
