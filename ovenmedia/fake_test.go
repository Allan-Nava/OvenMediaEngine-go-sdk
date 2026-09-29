package ovenmedia_test

import (
	"encoding/base64"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Allan-Nava/OvenMediaEngine-go-sdk/ovenmedia"
)

// The fake answers with the replies documented in the OME REST API v1
// reference (https://ovenmedia.com/docs/ome/rest-api/v1), trimmed. When OME
// changes a reply, change it here first and let the tests say what broke.

const token = "admin:secret"

var wantAuth = "Basic " + base64.StdEncoding.EncodeToString([]byte(token))

type reply struct {
	status int
	body   string
	delay  time.Duration
}

type recorded struct {
	Method string
	Path   string // escaped, as sent on the wire
	Auth   string
	Body   string
}

type fakeOME struct {
	*httptest.Server
	mu     sync.Mutex
	routes map[string]reply
	reqs   []recorded
}

const ok = `{"statusCode":200,"message":"OK"}`

const pushJSON = `{"id":"push1","state":"ready","vhost":"default","app":"app",
 "stream":{"name":"stream","trackIds":[],"variantNames":[]},
 "protocol":"rtmp","url":"rtmp://ingest.example.com/live2","streamKey":"KEY",
 "sentBytes":0,"sentTime":0,"sequence":0,"totalsentBytes":0,"totalsentTime":0,
 "createdTime":"2023-03-15T23:02:34.371+09:00",
 "startTime":"1970-01-01T09:00:00.000+09:00","finishTime":"1970-01-01T09:00:00.000+09:00"}`

const recordJSON = `{"id":"rec1","state":"recording","vhost":"default","app":"app",
 "stream":{"name":"stream","trackIds":[],"variantNames":[]},
 "interval":60000,"segmentationRule":"discontinuity",
 "createdTime":"2023-03-15T21:15:20.113+09:00"}`

const statsJSON = `{"statusCode":200,"message":"OK","response":{
 "connections":{"file":0,"hlsv3":0,"llhls":12,"ovt":0,"push":1,"srt":0,"thumbnail":0,"webrtc":30},
 "createdTime":"2023-03-15T19:46:13.728+09:00","lastRecvTime":"2023-03-15T19:46:13.728+09:00",
 "lastSentTime":"2023-03-15T19:46:13.728+09:00","lastUpdatedTime":"2023-03-15T19:46:13.728+09:00",
 "maxTotalConnectionTime":"2023-03-15T19:46:13.728+09:00",
 "lastThroughputIn":2500000,"lastThroughputOut":75000000,"maxTotalConnections":50,
 "totalBytesIn":494713880,"totalBytesOut":0,"totalConnections":43,
 "avgThroughputIn":0,"avgThroughputOut":0,"maxThroughputIn":0,"maxThroughputOut":0}}`

const streamInfoJSON = `{"statusCode":200,"message":"OK","response":{"name":"stream",
 "input":{"createdTime":"2026-07-29T15:04:21.879+09:00","sourceType":"Rtmp","sourceUrl":"TCP://203.0.113.10:41008",
  "connection":{"transport":"TCP","protocol":"TCP","localAddress":"198.51.100.1","localPort":1935,"remoteAddress":"203.0.113.10","remotePort":41008},
  "tracks":[{"id":0,"name":"Video","type":"Video","video":{"bypass":false,"codec":"H264","width":1280,"height":720,
   "bitrate":2500000,"framerate":30.0,"keyFrameInterval":30.0}},
   {"id":1,"name":"Audio","type":"Audio","audio":{"bypass":false,"codec":"AAC","bitrate":"128000","samplerate":48000,"channel":2}}]},
 "outputs":[{"name":"stream","tracks":[],"playlists":[]}]}}`

func docRoutes() map[string]reply {
	r := func(body string) reply { return reply{status: 200, body: body} }
	const app = "/v1/vhosts/default/apps/app"
	return map[string]reply{
		"GET /v1/version":         r(`{"statusCode":200,"message":"OK","response":{"version":"0.18.0","gitVersion":"v0.18.0-0-g1234567"}}`),
		"GET /v1/vhosts":          r(`{"statusCode":200,"message":"OK","response":["default"]}`),
		"POST /v1/vhosts":         r(`[{"statusCode":200,"message":"OK","response":{"name":"vhost","host":{"names":["ome.example.com"]}}}]`),
		"GET /v1/vhosts/default":  r(`{"statusCode":200,"message":"OK","response":{"name":"default","host":{"names":["*"]}}}`),
		"DELETE /v1/vhosts/vhost": r(ok),

		"GET /v1/vhosts/default/apps":         r(`{"statusCode":200,"message":"OK","response":["app"]}`),
		"POST /v1/vhosts/default/apps":        r(`[{"statusCode":200,"message":"OK","response":{"name":"app2","type":"live"}}]`),
		"GET " + app:                          r(`{"statusCode":200,"message":"OK","response":{"name":"app","type":"live","dynamic":false,"providers":{"rtmp":{}}}}`),
		"PATCH " + app:                        r(`{"statusCode":200,"message":"OK","response":{"name":"app","type":"live","providers":{"webrtc":{"timeout":60000}}}}`),
		"DELETE /v1/vhosts/default/apps/app2": r(ok),

		"GET " + app + "/outputProfiles":           r(`{"statusCode":200,"message":"OK","response":["bypass"]}`),
		"POST " + app + "/outputProfiles":          r(`[{"statusCode":200,"message":"OK","response":{"name":"bypass","outputStreamName":"${OriginStreamName}"}}]`),
		"GET " + app + "/outputProfiles/bypass":    r(`{"statusCode":200,"message":"OK","response":{"name":"bypass","outputStreamName":"${OriginStreamName}","encodes":{"videos":[{"bypass":true}]}}}`),
		"DELETE " + app + "/outputProfiles/bypass": r(ok),

		"GET " + app + "/streams":                   r(`{"statusCode":200,"message":"OK","response":["stream"]}`),
		"POST " + app + "/streams":                  {status: 201, body: `{"statusCode":201,"message":"Created"}`},
		"GET " + app + "/streams/stream":            r(streamInfoJSON),
		"DELETE " + app + "/streams/stream":         r(ok),
		"POST " + app + "/streams/stream:sendEvent": r(ok),

		"POST " + app + ":startPush": r(`{"statusCode":200,"message":"OK","response":` + pushJSON + `}`),
		"POST " + app + ":stopPush":  r(ok),
		"POST " + app + ":pushes":    r(`{"statusCode":200,"message":"OK","response":[` + pushJSON + `]}`),

		"POST " + app + ":startRecord": r(`{"statusCode":200,"message":"OK","response":` + strings.Replace(recordJSON, "recording", "ready", 1) + `}`),
		"POST " + app + ":stopRecord":  r(`{"statusCode":200,"message":"OK","response":` + strings.Replace(recordJSON, "recording", "stopping", 1) + `}`),
		"POST " + app + ":records":     r(`{"statusCode":200,"message":"OK","response":[` + recordJSON + `]}`),

		"GET /v1/stats/current/vhosts/default":                         r(statsJSON),
		"GET /v1/stats/current/vhosts/default/apps/app":                r(statsJSON),
		"GET /v1/stats/current/vhosts/default/apps/app/streams/stream": r(statsJSON),

		// served by the publisher, without auth
		"GET /app/stream/thumb.png": r("\x89PNG\r\n\x1a\n"),
	}
}

// startFake starts a server answering with docRoutes. Close it when done.
func startFake() *fakeOME {
	f := &fakeOME{routes: docRoutes()}
	f.Server = httptest.NewServer(http.HandlerFunc(f.serve))
	return f
}

func newFake(t testing.TB) *fakeOME {
	t.Helper()
	f := startFake()
	t.Cleanup(f.Close)
	return f
}

func (f *fakeOME) serve(w http.ResponseWriter, r *http.Request) {
	body, _ := io.ReadAll(r.Body)
	path := r.URL.EscapedPath()
	f.mu.Lock()
	f.reqs = append(f.reqs, recorded{Method: r.Method, Path: path, Auth: r.Header.Get("Authorization"), Body: string(body)})
	rep, found := f.routes[r.Method+" "+path]
	f.mu.Unlock()

	if strings.HasPrefix(path, "/v1/") && r.Header.Get("Authorization") != wantAuth {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		_, _ = io.WriteString(w, `{"statusCode":401,"message":"[HTTP] Authorization header is required to call API"}`)
		return
	}
	if !found {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = io.WriteString(w, `{"statusCode":404,"message":"[HTTP] Could not find the resource"}`)
		return
	}
	if rep.delay > 0 {
		select {
		case <-time.After(rep.delay):
		case <-r.Context().Done():
			return
		}
	}
	if rep.status == http.StatusNoContent || rep.body == "" {
		w.WriteHeader(rep.status)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(rep.status)
	_, _ = io.WriteString(w, rep.body)
}

// on replaces the reply for one route.
func (f *fakeOME) on(method, path string, rep reply) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.routes[method+" "+path] = rep
}

func (f *fakeOME) last(t testing.TB) recorded {
	t.Helper()
	f.mu.Lock()
	defer f.mu.Unlock()
	if len(f.reqs) == 0 {
		t.Fatal("no request reached the server")
	}
	return f.reqs[len(f.reqs)-1]
}

func (f *fakeOME) count() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.reqs)
}

func (f *fakeOME) client(t testing.TB, opts ...ovenmedia.Option) ovenmedia.IOvenMediaClient {
	t.Helper()
	c, err := ovenmedia.New(f.URL, append([]ovenmedia.Option{ovenmedia.WithAccessToken(token)}, opts...)...)
	if err != nil {
		t.Fatal(err)
	}
	return c
}
