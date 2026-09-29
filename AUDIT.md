# Audit — OvenMediaEngine-go-sdk

Audit of `main` at `cda577f` (tag `v0.4.75`), 2026-09-29. Reference: the OME REST API v1 docs
at https://ovenmedia.com/docs/ome/rest-api/v1 (the old `airensoft.gitbook.io` and
`docs.ovenmediaengine.com` links redirect there).

**How it was checked.** `go build`, `go vet`, `go test`, `golangci-lint` v2 and `govulncheck`
run locally (Go 1.27). Every *correctness* finding marked **confirmed** was reproduced with a
throwaway `httptest` server that returns the reply documented by OME. The server records what the
SDK sends, and the test prints what the SDK decodes. The probe was not committed. Findings marked
**plausible** follow from the docs but were not reproduced against a real OME.

Ids are stable: when an item is fixed, write `fixed in <sha>` next to it rather than deleting it.

## Resolution (branch `fix/audit`, for v0.5.0)

`9470355` is the SDK rewrite, `77dfe17` the CI rebuild. Items marked **open** are still out of scope.

| id | Status |
|---|---|
| A-01 | fixed in `9470355`: `CreateVirtualHost(ctx, ...VirtualHostConfig)` sends and decodes arrays |
| A-02, A-03 | fixed in `9470355`: `/streams/` path, `ResponseStats.Response` holds the stats, with `connections` and the throughput fields |
| A-04 | fixed in `9470355`: non-2xx replies return `*APIError`; `HealthCheck` calls `/v1/version` |
| A-05 | fixed in `9470355`: `GetRecordingState(ctx, vhost, app, id)` returns the list |
| A-06 | fixed in `9470355`: `ovenmedia.Time` accepts `+09:00` and `+0900` |
| A-07, A-08 | fixed in `9470355`: stream key required for RTMP only, `StopPush` takes an ID, recording paths optional, `trackIds`/`variantNames` supported; `validator.v2` removed |
| A-09 | fixed in `9470355`: `CreateOmeBasicAuthHeader*` set `Authorization`, deprecated |
| A-10 | fixed in `9470355`: `WithDebug` redacts `Authorization` in the log |
| A-11 | fixed in `9470355`: `ctx` on every method, 30s default timeout, `WithHTTPClient` |
| A-12, A-13 | fixed in `9470355`: no stdout output; typed reply; empty list on 204 |
| A-14, A-15 | fixed in `9470355`: `httptest` suite (78.6% coverage, every v0.4 bug re-introduced to check it fails) and 15 examples checked by `TestExamplesAreDocumented`; `test/` is an opt-in read-only integration suite |
| A-16 | fixed in `9470355` for the path builders (unexported functions, names escaped). The enum constants (`LIVE`, `H264`...) keep their names for compatibility |
| A-17 | fixed in `9470355`: `New` validates the URL, `baseURL` unexported, `getChangePort` removed, `CodecAudio` split from `CodecVideo` |
| A-18 | fixed in `77dfe17` |
| A-19 | fixed by the README rewrite |
| A-20 | **partly open**: version, vhost get/delete, application CRUD, output profiles, stream create/delete and `:sendEvent` added in `9470355`. Scheduled channels, multiplex channels and HLS dumps are still missing |
| A-21 | fixed in `9470355`: resty v2.17.2, `golang.org/x/net` v0.58.0, Go 1.25; govulncheck reports no vulnerabilities |
| A-22 | fixed in `77dfe17` |
| A-23 | fixed in `77dfe17`; `release.yml.old` removed in the cleanup commit |
| A-24 | fixed: Renovate config in `77dfe17`, `dependabot.yml` removed in the cleanup commit |
| A-25, A-26 | fixed by the Pages rewrite |
| A-27 | fixed: `.gitignore` in `77dfe17`, `.vscode/` untracked in the cleanup commit |
| A-28 | fixed in `9470355` (found while fixing): current OME sends `bitrate` as a number, older versions as a string; the v0.4 `string` field made `GetStreamInfo` fail on current servers. Now `FlexInt64` |

## Summary

| Area | Result |
|---|---|
| Build / vet | pass |
| Tests | pass, but they prove nothing: 4 tests, 0 assertions on SDK behaviour, 2 hit a public docs URL over the network; 13.7% of `ovenmedia/` statements executed, none asserted |
| Lint (golangci-lint v2) | 1 issue: `getChangePort` unused |
| govulncheck | 0 reachable; **16 known vulns in required modules** (`golang.org/x/net v0.8.0`, `resty v2.7.0`) |
| Endpoints | 16 methods: **6 broken** against the documented API, the other 10 hide HTTP errors |
| Docs / README | leftover chatbot text published on both; Jekyll site with a third-party analytics id |

## Correctness — broken endpoints

| id | Sev | Status | Finding |
|---|---|---|---|
| A-01 | high | confirmed | **`CreateVirtualHost` sends the wrong body and can't decode the reply.** OME expects a top-level array `[{"name":"v"}]`; the SDK sends `{"VirtualHostsName":[{"name":"v"}]}`. The reply is also an array, and unmarshalling it into the struct `ResponseVirtualHost` fails with `cannot unmarshal array`. `ovenmedia/virtualhost.go:12`, `request.go:3`, `response.go:11` |
| A-02 | high | confirmed | **Stream stats hit a non-existent path.** `V1_CURRENT_STATS_APPP_STREAMS_NAME` is `.../apps/%s/stream/%s`; OME's path is `.../streams/{stream}`. `GetStatsStreamVhosts` always 404s. `ovenmedia/constants.go:31` |
| A-03 | high | confirmed | **All three stats methods return zeroes.** OME wraps stats in `{"statusCode","message","response":{...}}`; `ResponseStats` has the fields at the top level and no `response`, so `{"response":{"totalConnections":42}}` decodes to `TotalConnections: 0` with `err == nil`. The type also lacks `connections`, the throughput fields, etc. `ovenmedia/response.go:155` |
| A-05 | high | confirmed | **`GetRecordingState` always errors.** `:records` returns `"response": [...]`; the method decodes into `ResponseRecordingStart`, whose `Response` is a single object. `ovenmedia/recording.go:55` |
| A-06 | high | confirmed (format) / plausible (OME version) | **`time.Time` fields break the decode.** `ResponseRecording.CreatedTime` (and `ResponsePushes`, `ResponseStreamInfo`) are `time.Time`, which only parses RFC 3339. OME's own recording example (quoted in `response.go:187`) uses `+0900`, which fails: `cannot parse "+0900" as "Z07:00"`, and one bad field fails the whole reply. Introduced by `d60eef8`. Newer OME builds print `+09:00` in the stats examples, so whether it bites depends on the server version — a custom time type that accepts both is the safe fix |
| A-07 | medium | confirmed | **Push validation rejects valid requests.** `RequestBodyPush.StreamKey` is `validate:"min=1"`, but OME only needs a stream key for RTMP: an SRT push fails client-side with `StreamKey: less than min`. `StopPush` takes the full `RequestBodyPush` and validates it, though OME's `:stopPush` body is just `{"id"}`. `ovenmedia/request.go:7`, `push.go:31` |
| A-08 | low | confirmed | **`RequestRecordingStart` tags are inert and wrong.** `filePath`/`infoPath` carry `validate` tags but `StartRecording` never calls the validator, and OME marks both optional; without `omitempty` they're sent as `""`. The push/record `stream` object uses `tracks`; current OME documents `variantNames`. `ovenmedia/request.go:32` |

## Correctness — error handling and safety

| id | Sev | Status | Finding |
|---|---|---|---|
| A-04 | high | confirmed | **HTTP errors are returned as success.** No method checks `resp.StatusCode()` or `resp.IsError()`. Against a 401 `GetAllVirtualHosts` returns `err == nil` and `StatusCode: 401`. `HealthCheck` returns `nil` on 401 or 500 — it only fails when the TCP connection does. Every caller has to remember to check `StatusCode`, and for stats (A-03) it isn't even decoded. `ovenmedia/ovenmedia.go:42` and every endpoint file |
| A-09 | medium | confirmed (docs) | **Two auth helpers set a header OME doesn't read.** OME authenticates with `Authorization: Basic base64(AccessToken)`. `CreateOmeBasicAuthHeader` / `CreateOmeBasicAuthHeaderEncoded` write `ome-access-token`, so a client configured with them is unauthenticated. `ovenmedia/header_configurator.go:90` |
| A-10 | medium | confirmed | **Debug mode logs credentials.** `BuildOven(..., debug=true, ...)` calls resty `SetDebug(true)`, which prints every request's headers, `Authorization` included, to the standard logger. The test suite runs with `debug=true`. `ovenmedia/builder.go:26` |
| A-11 | medium | confirmed | **No timeout, no context.** `resty.New()` with no timeout, and no method takes a `context.Context`: a hung OME blocks the caller forever and a request can't be cancelled. There's no way to pass a custom resty/http client either. `ovenmedia/builder.go:14` |
| A-12 | low | confirmed | **`StopPush` prints to stdout.** `fmt.Println("resp", resp)` in library code, next to a `// TODO:`; it also returns the raw `*resty.Response` instead of a typed reply. `ovenmedia/push.go:41` |
| A-13 | low | confirmed | `GetAllPushes` returns `nil, nil` on 204. OME documents an empty `response: []` instead; a caller that dereferences the result panics. `ovenmedia/push.go:53` |

## Tests

| id | Sev | Finding |
|---|---|---|
| A-14 | high | The suite doesn't test the SDK. `test/a_main_test.go` builds a client against `https://airensoft.gitbook.io/...` (a docs page, now a redirect) and `HealthCheck` passes because any HTTP reply is success (A-04). `Test_GetAllVirtualHost` ignores the error and prints; `Test_HeaderInit` prints. No `httptest`, no assertions, so every bug above ships green. `ovenmedia/` has no tests of its own; the `test/` package executes 13.7% of its statements without checking any result. Tests need the network, so CI is flaky by design |
| A-15 | medium | No `Example...` functions, so pkg.go.dev shows no usage and the README is the only guide — and the README is wrong (A-19) |

## API design (breaking to fix — batch them into one `v0.5.0`)

| id | Finding |
|---|---|
| A-16 | Path builders are exported mutable `var`s (`GET_VHOSTS_BY_NAME = func...`): any importer can reassign them for everyone. Screaming-case exported names (`V1_HOSTS`, `GET_THUMBNAIL`) go against Go style, and golint-style linters flag them; typo `APPP` in two constant names |
| A-17 | `BuildOven` returns an `error` it can never produce; the `ovenMedia.Url` field is exported on an unexported type; `getChangePort` is dead code (the only lint issue); `CodecVideo` contains the audio codecs `OPUS` and `AAC` |
| A-20 | Missing endpoints: vhost get/delete, application create/get/update/delete, stream create (pull)/delete, `outputProfiles`, `:sendEvent`, HLS dump, multiplex and scheduled channels. The README's "Control stream playback" feature doesn't exist |

## Build, CI, dependencies

| id | Sev | Finding |
|---|---|---|
| A-18 | medium | `Makefile` runs `go build .` — the root has no Go files, so `make build` and `make test` fail with `no Go files`. `make` targets also `echo` noise |
| A-21 | medium | **Dependencies are three years old.** `resty v2.7.0` pulls `golang.org/x/net v0.8.0`, with 16 known vulnerabilities (none reachable from this code today, per govulncheck). `go 1.18`, and `gopkg.in/validator.v2` is unmaintained. Raising the floor to a supported Go and bumping resty clears them |
| A-22 | medium | CI tests Go 1.18/1.19, both years out of support. `cache-dependency-path: subdir/go.sum` points to a path that doesn't exist. `go-test.yml` runs `go mod tidy` (which rewrites) rather than `go mod tidy -diff` (which checks), and only on `push`, not PRs. `go-lint.yml` uses `golangci-lint-action@v2` with `version: latest`. No `-race`, no govulncheck. No `permissions:` block on any workflow |
| A-23 | low | `tag-autorelease.yml` installs **ffmpeg** for a library with no ffmpeg dependency, uses the archived `actions/create-release@v1` with `permissions: write-all`, and publishes no notes. `release.yml.old` is a dead copy of it |
| A-24 | low | Dependabot **and** Renovate are both configured and will open duplicate PRs; `dependabot.yml` also watches `/tests`, which doesn't exist. MistServer-go-sdk settled on Renovate only (`config:recommended`, grouped actions, `gomodTidy`, the `go` directive left for a human) |

## Docs and repository hygiene

| id | Sev | Finding |
|---|---|---|
| A-19 | high (reputation) | **README and the Pages home page publish chatbot output**: "I'm sorry, I'm unable to access external resources such as Github…" and "Please note that this is a sample". The install block is fenced as `go` and the usage shows `import "github.com/Allan-Nava/OvenMediaEngine-go-sdk"` — the root package doesn't exist; the import path is `.../ovenmedia`. **README fixed by the README rewrite**; the old Pages copy went with the Pages rewrite. The README also lacked the `readme: contributors` markers, so `contributors.yml` had nothing to update — added |
| A-25 | medium | The Jekyll `_config.yml` (just-the-docs template) carried `ga_tracking: UA-2709176-10`, which appears to be the value copied from the template rather than a property of this project, and is a Universal Analytics id — a product Google has shut down. The footer HTML has an unclosed attribute and links to `LICENSE.md`, which doesn't exist. **Fixed by the Pages rewrite** (static site, no analytics) |
| A-26 | low | The old `docs/*.md` pages copied OME request examples without saying which SDK method covers them (and `output-profile.md` documents an endpoint the SDK doesn't implement). Folded into the endpoint table of the new site |
| A-27 | low | `.vscode/launch.json` committed; `.gitignore` lists four individual `.idea/` files instead of `.idea/`; no `.DS_Store` rule |

History check (`git log -p --all`) for credentials, private hosts and IPs: nothing found beyond
the commit authors' own addresses, which are already public on the forge.

## Suggested order

1. **`v0.4.76` (non-breaking fixes):** A-02, A-03, A-05, A-13 (return an empty list), A-12, A-18,
   A-19. Each with an `httptest` test that fails first. Replace `test/` with such tests (A-14).
2. **CI:** one `ci.yml` modelled on MistServer-go-sdk (tidy check, vet, race tests on the floor
   and on stable, golangci-lint v2, govulncheck), `gh release create --generate-notes` on tags,
   Renovate only (A-21 to A-24).
3. **`v0.5.0` (breaking, one release):** A-01, A-04 (typed `*APIError` carrying `StatusCode` and
   `Message`), A-06, A-07, A-09, A-10, A-11 (`context.Context` and options such as `WithHTTPClient`,
   `WithTimeout`), A-16, A-17. Then the missing endpoints (A-20), each with an `Example`.
