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

The package lives in the `ovenmedia` subdirectory.

## Usage

Enable the API server in OME's `Server.xml` (`<Bind><Managers><API>`, commonly port `8081`) and
set an `AccessToken`. OME expects `Authorization: Basic base64(AccessToken)`.

```go
import (
	"fmt"
	"log"

	"github.com/Allan-Nava/OvenMediaEngine-go-sdk/ovenmedia"
)

headers := ovenmedia.InitHeaderConfigurator()
headers.CreateBasicAuthHeader("admin", "secret") // AccessToken "admin:secret"

client, err := ovenmedia.BuildOven("http://ome.example.com:8081", false, headers)
if err != nil {
	log.Fatal(err)
}

hosts, err := client.GetAllVirtualHosts()
if err != nil {
	log.Fatal(err)
}
if hosts.StatusCode != 200 { // OME error replies are decoded, not returned as err
	log.Fatalf("OME: %d %s", hosts.StatusCode, hosts.Message)
}
fmt.Println(hosts.Response) // [default]
```

Always check `StatusCode` as well as `err`: the client doesn't yet turn HTTP errors into Go
errors. Don't enable `debug` in production, because it logs request headers, including
`Authorization`.

## What's covered

`IOvenMediaClient` covers virtual hosts, applications, streams, pushes (`:startPush`,
`:stopPush`, `:pushes`), recordings (`:startRecord`, `:stopRecord`, `:records`), current
statistics and thumbnails. The [project site](https://allan-nava.github.io/OvenMediaEngine-go-sdk/#api)
lists each method with its endpoint and status.

## Status

Pre-1.0. An audit against the current OME docs found methods that don't work yet
(`CreateVirtualHost`, the three stats methods, `GetRecordingState`) along with other issues.
They are listed, with stable ids and suggested fixes, in [AUDIT.md](AUDIT.md).

## Contributing

Contributions are welcome: open an issue or a pull request. Read [CLAUDE.md](CLAUDE.md) for the
layout, conventions and definition of done. Tests must use `httptest`, not a live server.

## Contributors

<!-- readme: contributors -start -->
<!-- readme: contributors -end -->

## License

MIT. See [LICENSE](LICENSE).
