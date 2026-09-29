<p align="center">
  <a href="https://allan-nava.github.io/OvenMediaEngine-go-sdk/"><img src="docs/logo.svg" width="112" height="112" alt="OvenMediaEngine Go SDK logo"></a>
</p>

# OvenMediaEngine Go SDK

[![Go build](https://github.com/Allan-Nava/OvenMediaEngine-go-sdk/actions/workflows/go-build.yml/badge.svg)](https://github.com/Allan-Nava/OvenMediaEngine-go-sdk/actions/workflows/go-build.yml)
[![Go test](https://github.com/Allan-Nava/OvenMediaEngine-go-sdk/actions/workflows/go-test.yml/badge.svg)](https://github.com/Allan-Nava/OvenMediaEngine-go-sdk/actions/workflows/go-test.yml)
[![Go Lint](https://github.com/Allan-Nava/OvenMediaEngine-go-sdk/actions/workflows/go-lint.yml/badge.svg)](https://github.com/Allan-Nava/OvenMediaEngine-go-sdk/actions/workflows/go-lint.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/Allan-Nava/OvenMediaEngine-go-sdk/ovenmedia.svg)](https://pkg.go.dev/github.com/Allan-Nava/OvenMediaEngine-go-sdk/ovenmedia)

A Go client for the [OvenMediaEngine](https://ovenmedia.com/docs/ome/) (OME) REST API v1.
OME is an open-source, sub-second latency live streaming server: it ingests WebRTC, SRT, RTMP,
MPEG-TS and RTSP and delivers LL-HLS and WebRTC.

[OME REST API reference](https://ovenmedia.com/docs/ome/rest-api/v1) ·
[project site](https://allan-nava.github.io/OvenMediaEngine-go-sdk/)

## Installation

```bash
go get github.com/Allan-Nava/OvenMediaEngine-go-sdk
```

The package lives in the `ovenmedia` subdirectory. Requires Go 1.25 or newer.

## Usage

Enable the API server in OME's `Server.xml` (`<Bind><Managers><API>`, commonly port `8081`) and
set an `AccessToken`. The client sends it as `Authorization: Basic base64(AccessToken)`.

```go
import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/Allan-Nava/OvenMediaEngine-go-sdk/ovenmedia"
)

client, err := ovenmedia.New("http://ome.example.com:8081",
	ovenmedia.WithAccessToken("admin:secret"), // the AccessToken from Server.xml
)
if err != nil {
	log.Fatal(err)
}

hosts, err := client.GetAllVirtualHosts(context.Background())
var apiErr *ovenmedia.APIError
if errors.As(err, &apiErr) {
	log.Fatalf("OME said %d: %s", apiErr.StatusCode, apiErr.Message)
} else if err != nil {
	log.Fatal(err) // network, timeout, or ovenmedia.ErrInvalidRequest
}
fmt.Println(hosts.Response) // [default]
```

- Every method takes a `context.Context`; requests time out after 30s unless you pass
  `WithTimeout` or your own client with `WithHTTPClient`.
- A non-2xx reply from OME is an `*ovenmedia.APIError`. A request the SDK can tell is invalid
  (a push without a URL, say) fails with `ovenmedia.ErrInvalidRequest` before it is sent.
- `WithDebug(true)` logs requests and replies with the `Authorization` header redacted.
- Thumbnails are served by OME's publisher port, so `GetThumbnail` takes that URL separately and
  doesn't send the API credentials to it.

## What's covered

Virtual hosts, applications and output profiles (create, list, get, update, delete), streams
(pull, list, info, delete, `:sendEvent`), pushes (`:startPush`, `:stopPush`, `:pushes`),
recordings (`:startRecord`, `:stopRecord`, `:records`), current statistics, the version and
thumbnails. The [project site](https://allan-nava.github.io/OvenMediaEngine-go-sdk/#api) maps
each method to its endpoint, and there is a runnable example for each area on
[pkg.go.dev](https://pkg.go.dev/github.com/Allan-Nava/OvenMediaEngine-go-sdk/ovenmedia#pkg-examples).

Not covered yet: scheduled channels, multiplex channels and HLS dumps.

## Upgrading from v0.4

v0.5.0 is a breaking release: see [AUDIT.md](AUDIT.md) for what changed and why. In short, add
a `ctx` as the first argument everywhere, replace `BuildOven` with `New` (it still works, but is
deprecated), pass an ID to `StopPush`, `StopRecording` and `GetRecordingState`, and handle
`*APIError` instead of checking `StatusCode`.

## Testing

```bash
make test lint
```

Unit tests run against an `httptest` fake that serves the replies from the OME docs, so they
need no server. `make integration OME_URL=... OME_ACCESS_TOKEN=...` runs read-only checks against
a real OME.

## Contributing

Contributions are welcome: open an issue or a pull request. Read [CLAUDE.md](CLAUDE.md) for the
layout, conventions and definition of done: tests first, against the fake, and an example for
any public change.

## Contributors

<!-- readme: contributors -start -->
<table>
	<tbody>
		<tr>
            <td align="center">
                <a href="https://github.com/Allan-Nava">
                    <img src="https://avatars.githubusercontent.com/u/22498435?v=4" width="100;" alt="Allan-Nava"/>
                    <br />
                    <sub><b>Allan Nava</b></sub>
                </a>
            </td>
		</tr>
	<tbody>
</table>
<!-- readme: contributors -end -->

## License

MIT. See [LICENSE](LICENSE).
