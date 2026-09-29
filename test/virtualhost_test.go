package test

import (
	"context"
	"testing"
)

// Test_ReadOnlyWalk lists every virtual host, application and stream and
// reads their stats: it changes nothing on the server.
func Test_ReadOnlyWalk(t *testing.T) {
	c := omeClient(t)
	ctx := context.Background()

	vhosts, err := c.GetAllVirtualHosts(ctx)
	if err != nil {
		t.Fatal(err)
	}
	for _, vh := range vhosts.Response {
		if _, err := c.GetStatsVhosts(ctx, vh); err != nil {
			t.Errorf("stats %s: %v", vh, err)
		}
		apps, err := c.GetApplications(ctx, vh)
		if err != nil {
			t.Errorf("apps %s: %v", vh, err)
			continue
		}
		for _, app := range apps.Response {
			if _, err := c.GetAllPushes(ctx, vh, app); err != nil {
				t.Errorf("pushes %s/%s: %v", vh, app, err)
			}
			if _, err := c.ListRecordingState(ctx, vh, app); err != nil {
				t.Errorf("records %s/%s: %v", vh, app, err)
			}
			streams, err := c.GetStreams(ctx, vh, app)
			if err != nil {
				t.Errorf("streams %s/%s: %v", vh, app, err)
				continue
			}
			for _, s := range streams.Response {
				if _, err := c.GetStreamInfo(ctx, vh, app, s); err != nil {
					t.Errorf("stream %s/%s/%s: %v", vh, app, s, err)
				}
				if _, err := c.GetStatsStreamVhosts(ctx, vh, app, s); err != nil {
					t.Errorf("stats %s/%s/%s: %v", vh, app, s, err)
				}
			}
		}
	}
}
