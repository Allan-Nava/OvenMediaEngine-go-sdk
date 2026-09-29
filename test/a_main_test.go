// Package test holds integration tests against a real OvenMediaEngine.
// They only read state and are skipped unless OME_URL is set:
//
//	OME_URL=http://localhost:8081 OME_ACCESS_TOKEN=admin:secret go test ./test/
//
// The unit tests, which need no server, live next to the code in ovenmedia/.
package test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/Allan-Nava/OvenMediaEngine-go-sdk/ovenmedia"
)

func omeClient(t *testing.T) ovenmedia.IOvenMediaClient {
	t.Helper()
	url := os.Getenv("OME_URL")
	if url == "" {
		t.Skip("OME_URL not set; skipping integration test")
	}
	c, err := ovenmedia.New(url,
		ovenmedia.WithAccessToken(os.Getenv("OME_ACCESS_TOKEN")),
		ovenmedia.WithTimeout(10*time.Second),
	)
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func Test_HealthCheck(t *testing.T) {
	c := omeClient(t)
	if err := c.HealthCheck(context.Background()); err != nil {
		t.Fatal(err)
	}
}
