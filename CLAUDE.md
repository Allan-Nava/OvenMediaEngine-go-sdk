# CLAUDE.md — OvenMediaEngine-go-sdk

Go client for the [OvenMediaEngine](https://ovenmedia.com/docs/ome/) (OME) REST API v1
(`/v1/...`, usually on the API port `:8081`, or `:8082` with TLS). Small library, no binary: one
package, `ovenmedia/`. Module `github.com/Allan-Nava/OvenMediaEngine-go-sdk`, Go 1.25+ (floor set
by `golang.org/x/net`, pulled in by resty). Tags `v0.4.x` are published; the reworked API is
`v0.5.0` and breaks every caller (see README "Upgrading from v0.4").

`AUDIT.md` is the single backlog: every known defect or gap has a stable id (`A-01`...), and
fixed items stay listed with the commit that fixed them.

## Working rules (ALWAYS)

- **Never `git push`** — the user pushes. Never a `Co-Authored-By` trailer or any tool footer
  in commits, PRs or docs. Commit identity for this repo is set in the local git config
  (`Allan Nava`, the gmail address); don't override it.
- **Never tag** without being asked — a tag is a release on the Go proxy and can't be taken back.
- **One backlog**: open problems go in `AUDIT.md`. No `TODO:` comments in code.
- **Align everything**: a change to the public surface updates, in the same change, the doc
  comment, an `Example`, `README.md`, `docs/index.html` (endpoint table + quickstart), `AUDIT.md`
  and this file. Grep all of them for what the change made stale.
- **Unit tests never hit the network.** They use the fake in `ovenmedia/fake_test.go`. The
  `test/` package is the only place that talks to a real OME, read-only, and only when `OME_URL`
  is set.
- Public repo: no internal hostnames, customer names, real IPs, tokens or `/Users/...` paths in
  code, fixtures, commits or examples. Use `ome.example.com`, RFC 5737 addresses
  (`203.0.113.x`, `198.51.100.x`) and obviously fake credentials (`admin:secret`).

## Commands

```bash
make build        # go build ./...  (the root has no Go files, never `go build .`)
make test         # go vet + go test -race
make lint         # golangci-lint v2
make vuln         # govulncheck
make integration OME_URL=http://localhost:8081 OME_ACCESS_TOKEN=admin:secret   # read-only, real OME
```

CI (`.github/workflows/`): `go-build.yml` (tidy check, build, vet on 1.25.x and stable),
`go-test.yml` (race tests, both versions), `go-lint.yml` (golangci-lint v2 + govulncheck),
`tag-autorelease.yml` (`gh release create --generate-notes` on `v*.*.*`), `contributors.yml`
(fills the README contributors block). The file names are what the README badges point at —
rename one and update the badge. Renovate keeps modules and actions current and leaves the `go`
directive alone; raising the floor means updating the first entry of both CI matrices.

## Layout

- `ovenmedia/` — the SDK, package `ovenmedia`.
  - `ovenmedia.go` — package doc, `IOvenMediaClient` (the whole public surface), the unexported
    `ovenMedia` struct, and `do[T]`, the one request helper.
  - `builder.go` — `New(baseURL, opts...)` and the options (`WithAccessToken`, `WithHeaders`,
    `WithHTTPClient`, `WithTimeout`, `WithDebug`); `BuildOven` is the deprecated v0.4 constructor.
  - `errors.go` — `APIError`. `ErrInvalidRequest` lives in `ovenmedia.go`.
  - `types.go` — `Time` (OME timestamps) and `FlexInt64` (numbers sent as number or string).
  - `constants.go` — unexported path builders (`pathApp`, `pathAppAction`...) and the enums.
  - One file per API area: `virtualhost.go`, `application.go` (+ output profiles), `stream.go`
    (+ `SendEvent`), `push.go`, `recording.go`, `stats.go`, `thumbnail.go`.
  - `request.go` / `response.go` — JSON payloads; request types carry their `validate()`.
  - `header_configurator.go` — `HeaderConfigurator`, for `WithHeaders`.
  - Tests (package `ovenmedia_test`): `fake_test.go` (the fake OME), `ovenmedia_test.go`,
    `example_test.go`; `doc_test.go` (package `ovenmedia`) checks every example is attached.
- `test/` — opt-in integration tests against a real OME, read-only.
- `docs/` — GitHub Pages, served from `main:/docs` (legacy build). Static HTML + `.nojekyll`,
  no build step. `index.html` documents the API by hand.

## How a call works

Every endpoint method is a thin wrapper: validate the arguments (returning an error wrapping
`ErrInvalidRequest` before anything is sent), then `do[Resp](ctx, o, method, path, body)`.
`do` sets the context and JSON headers, sends through the resty client built by `New` (base URL,
static headers, timeout), turns any non-2xx reply into an `*APIError` (status, OME's `message`,
raw body), returns the zero `Resp` for an empty 204, and otherwise unmarshals into `Resp`.

OME wraps replies in `{"statusCode", "message", "response"}`: response types embed
`BaseResponseOK` and put the payload in `Response`. Creates of vhosts, applications and output
profiles take **and return a top-level array**, so those methods return a slice.

Auth is HTTP Basic: `Authorization: Basic base64(AccessToken)`, the plaintext token from
`Server.xml` (`<Managers><API><AccessToken>`), often `user:password`. `WithAccessToken` does the
encoding. `WithDebug` turns on resty's request log with an `OnRequestLog` hook that redacts
`Authorization` in the log copy only.

Thumbnails aren't on the API server: `GetThumbnail` takes the publisher URL and uses the
underlying `http.Client` directly, so resty's headers (the API credentials) aren't sent to it.

### Adding an endpoint

1. Read the OME reference page (https://ovenmedia.com/docs/ome/rest-api/v1): method, path,
   body, reply. Actions use the `:verb` suffix and are `POST`.
2. Add its documented reply to `docRoutes()` in `fake_test.go` — the fake is the contract.
3. Write the test: it asserts method, path **and** body with `wantRequest`, and checks decoded
   fields. See it fail.
4. Path builder in `constants.go` (escape every name with `esc`), request type in `request.go`
   (`omitempty` on every optional field, a `validate()` if something is required), reply type in
   `response.go` embedding `BaseResponseOK`, method on `IOvenMediaClient` + `*ovenMedia`.
5. An `Example<RequestType>` (or `Example<Type>`) in `example_test.go` with `// Output:`.
6. Row in the `docs/index.html` endpoint table, README "What's covered", `AUDIT.md` A-20.

## Gotchas (verified against the OME docs, 2026-09)

- Timestamps: some OME versions print `+0900` without the colon. Use `Time`, never `time.Time`:
  one unparsable field fails the whole `Unmarshal`.
- Bitrates arrive as numbers in current OME and as strings in older versions: use `FlexInt64`.
- The push/record stream selector is `trackIds` + `variantNames` today; `tracks` is the older
  name, kept on `SimpleStream` for old servers.
- `:records` always answers with a list, even when filtered by ID.
- `:pushes` may answer 204; `GetAllPushes` then returns an empty, non-nil list.
- Stream creation answers **201**, not 200 — `do` accepts any 2xx.
- SRT and MPEG-TS pushes send `"streamKey": ""`; only RTMP requires a key.
- `SendEvent` defaults `eventFormat` to `id3v2`, the only format OME supports; OME answers 409
  if the stream has no media yet.
- `resty.NewWithClient(hc)` + `SetTimeout` writes `hc.Timeout`: `WithHTTPClient` without
  `WithTimeout` keeps the caller's timeout untouched on purpose.
- `TestDebugRedactsAuthorization` swaps `os.Stderr` while calling `New`, because resty's default
  logger captures it at construction. Don't run it in parallel.

## Conventions

- **TDD**: red first. The fake returns documented replies; a test that only pins existing
  behaviour must be shown to fail on a mutation (the v0.5 suite was checked by re-introducing
  each v0.4 bug).
- **Examples**: every public API change ships an `Example...` with an `// Output:` block, named
  after a function or type go/doc can attach it to (`ExampleNew_healthCheck`,
  `ExampleRequestBodyPush`). Examples named after interface methods compile but pkg.go.dev drops
  them; `TestExamplesAreDocumented` fails on those.
- Server behaviour comes from the OME docs or source (`AirenSoft/OvenMediaEngine`,
  `src/projects/api_server/controllers/v1/`), never from guesses. Quote the page in the commit.
- Keep `ovenMedia` unexported; the public surface is `IOvenMediaClient`. Adding a method to the
  interface breaks anyone who implements it (mocks) — say so in the commit.
- Errors are prefixed `ovenmedia:`; no `fmt.Println` or `log` in library code (the one
  `log.Println` in the deprecated `BuildOven` is kept for v0.4 behaviour).

## Definition of done

1. Test first, then the fix; an `Example...` for any public API change.
2. Doc comments, `README.md`, `docs/index.html`, `AUDIT.md` and this file updated in the same
   change.
3. `make test lint vuln` green locally and CI green.
4. After the user tags: the release workflow succeeded and pkg.go.dev lists the version with
   its examples. If it 404s:
   `curl -X POST https://pkg.go.dev/fetch/github.com/Allan-Nava/OvenMediaEngine-go-sdk@vX.Y.Z`.
