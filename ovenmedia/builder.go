package ovenmedia

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

// DefaultTimeout bounds every request unless [WithTimeout] or
// [WithHTTPClient] says otherwise.
const DefaultTimeout = 30 * time.Second

type config struct {
	httpClient *http.Client
	timeout    time.Duration
	headers    map[string]string
	debug      bool
}

// Option configures the client built by [New].
type Option func(*config)

// WithAccessToken sends Authorization: Basic base64(token) on every request.
// token is the plaintext AccessToken from Server.xml, often "user:password".
func WithAccessToken(token string) Option {
	return func(c *config) {
		c.headers["Authorization"] = "Basic " + base64.StdEncoding.EncodeToString([]byte(token))
	}
}

// WithHeaders adds every header of h to each request. Later options win.
func WithHeaders(h *HeaderConfigurator) Option {
	return func(c *config) {
		if h == nil {
			return
		}
		for k, v := range h.GetHeaders() {
			c.headers[k] = v
		}
	}
}

// WithHTTPClient uses hc for transport, TLS and redirects. Its Timeout is
// kept unless WithTimeout is also given.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *config) { c.httpClient = hc }
}

// WithTimeout bounds each request, including reading the reply. Zero means
// no limit; prefer a context deadline for that.
func WithTimeout(d time.Duration) Option {
	return func(c *config) { c.timeout = d }
}

// WithDebug logs every request and response. The Authorization header is
// redacted, but bodies are logged as they are.
func WithDebug(debug bool) Option {
	return func(c *config) { c.debug = debug }
}

// New returns a client for the OME API server at baseURL, for example
// "http://ome.example.com:8081". It fails if baseURL isn't an absolute
// http(s) URL.
func New(baseURL string, opts ...Option) (IOvenMediaClient, error) {
	u, err := url.Parse(baseURL)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return nil, fmt.Errorf("ovenmedia: base URL must be an absolute http(s) URL, got %q", baseURL)
	}
	cfg := config{timeout: -1, headers: map[string]string{}}
	for _, opt := range opts {
		opt(&cfg)
	}

	var rc *resty.Client
	if cfg.httpClient != nil {
		rc = resty.NewWithClient(cfg.httpClient)
	} else {
		rc = resty.New()
		if cfg.timeout < 0 {
			cfg.timeout = DefaultTimeout
		}
	}
	if cfg.timeout >= 0 {
		rc.SetTimeout(cfg.timeout)
	}
	rc.SetBaseURL(strings.TrimRight(baseURL, "/"))
	rc.SetHeaders(cfg.headers)
	if cfg.debug {
		rc.SetDebug(true)
		rc.OnRequestLog(func(l *resty.RequestLog) error {
			if l.Header.Get("Authorization") != "" {
				l.Header.Set("Authorization", "[redacted]")
			}
			return nil
		})
	}
	return &ovenMedia{baseURL: baseURL, debug: cfg.debug, restClient: rc}, nil
}

// BuildOven is the pre-0.5 constructor.
//
// Deprecated: use [New] with [WithDebug] and [WithHeaders].
func BuildOven(url string, debug bool, header *HeaderConfigurator) (IOvenMediaClient, error) {
	c, err := New(url, WithDebug(debug), WithHeaders(header))
	if err == nil && debug {
		log.Println("ovenmedia: debug mode enabled")
	}
	return c, err
}
