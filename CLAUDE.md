# CLAUDE.md — OvenMediaEngine-go-sdk

Go client for the [OvenMediaEngine](https://ovenmedia.com/docs/ome/) (OME) REST API v1
(`/v1/...`, usually on the API port `:8081`, or `:8082` with TLS). Small library, no binary: one
package, `ovenmedia/`. Module `github.com/Allan-Nava/OvenMediaEngine-go-sdk`, `go 1.18`.
Tags `v0.4.x` are published; the API is pre-1.0, but every exported identifier has importers.

Known defects, missing endpoints and CI debt live in **`AUDIT.md`**: one list, stable ids
(`A-01`...). Read it before touching an endpoint — several of them are broken today.

## Working rules (ALWAYS)

- **Never `git push`** — the user pushes. Never a `Co-Authored-By` trailer or any tool footer
  in commits, PRs or docs.
- **Never tag** without being asked — a tag is a release on the Go proxy and can't be taken back.
- **One backlog**: open problems go in `AUDIT.md` (keep the id, mark it `fixed in <sha>` instead
  of deleting it). No scattered `TODO:` comments in code.
- **Align everything**: a change that touches the public surface updates, in the same commit,
  the doc comment, `README.md`, `docs/index.html` (endpoint table + quickstart), `AUDIT.md`
  and this file. Grep all of them for what the change made stale.
- **Tests never hit the network.** A test that needs OME uses an `httptest.Server` that mirrors
  the documented reply. The current `test/` suite calls a public docs URL and asserts nothing
  (`AUDIT.md` A-14) — don't copy that pattern, replace it.
- Public repo: no internal hostnames, customer names, IPs, tokens or `/Users/...` paths in code,
  fixtures, commit messages or examples. Use `ome.example.com` and obviously fake credentials
  (`admin:secret`).

## Commands

```bash
go build ./...                      # the root has no Go files, never `go build .` (the Makefile does — A-18)
go vet ./...
go test -race -count=1 ./...
golangci-lint run ./...             # v2; today: 1 issue (unused getChangePort)
go run golang.org/x/vuln/cmd/govulncheck@latest ./...
```

CI (`.github/workflows/`): `go-build.yml`, `go-test.yml`, `go-lint.yml`, `tag-autorelease.yml`
(release on `v*.*.*`), `contributors.yml` (rewrites the README contributors block on every push to
main). All pinned to old actions and Go 1.18/1.19 — see `AUDIT.md` §CI.

## Layout

- `ovenmedia/` — the SDK, package `ovenmedia`.
  - `ovenmedia.go` — `IOvenMediaClient` interface (the whole public surface), unexported
    `ovenMedia` struct, resty helpers `get` / `post` / `postNoBody`.
  - `builder.go` — `BuildOven(url, debug, header)`: the only constructor.
  - `header_configurator.go` — `HeaderConfigurator`, a map of headers applied to every request;
    auth helpers live here.
  - `constants.go` — path templates (`V1_*`) and the `GET_*` path builders, plus enums
    (`ApplicationType`, `CodecVideo`, `SessionState`...).
  - One file per API area: `virtualhost.go`, `application.go`, `stream.go`, `push.go`,
    `recording.go`, `stats.go`, `thumbnail.go`.
  - `request.go` / `response.go` — JSON payloads.
- `test/` — external package `test`, integration-style tests (see Working rules).
- `docs/` — GitHub Pages, served from `main:/docs` (legacy build). Static HTML + `.nojekyll`,
  no Jekyll, no build step. `index.html` documents the API by hand.

## How a call works

`BuildOven` creates a `resty.Client` with the base URL and the static headers. Each method builds
the path with a `GET_*` helper, calls `get`/`post`, then `json.Unmarshal`s `resp.Body()` into a
`Response*` type. OME wraps every reply in `{"statusCode", "message", "response"}`; the Go types
embed `BaseResponseOK` for the first two.

**The HTTP status is never checked.** A 401/404/500 with a JSON body decodes fine and returns
`err == nil`: callers must look at `StatusCode`. `HealthCheck` only fails on a network error.
This is `AUDIT.md` A-04 — fixing it changes behaviour for every method, do it as one deliberate
change with a new error type, not piecemeal.

### Auth

OME uses HTTP Basic: `Authorization: Basic base64(<AccessToken>)`, where `AccessToken` is the
plaintext string from `Server.xml` (`<Managers><API><AccessToken>`), often `user:password`.
In `HeaderConfigurator`, the right helpers are `CreateBasicAuthHeader(user, pass)` and
`CreateOmeBasicAuthHeaderWord(token)`. `CreateOmeBasicAuthHeader*` set an `ome-access-token`
header OME doesn't read (A-09).

### Adding an endpoint

1. Read the OME reference page for it (https://ovenmedia.com/docs/ome/rest-api/v1): method,
   path, body, reply. Actions use the `:verb` suffix (`apps/app:startPush`) and are all `POST`.
2. Path template + builder in `constants.go`, request type in `request.go` (with `omitempty` on
   every field OME marks optional), reply type in `response.go` embedding `BaseResponseOK`.
3. Method on `IOvenMediaClient` + `*ovenMedia`, in the file of its area.
4. `httptest` test that asserts method, path **and** body sent, and decodes the documented reply.
5. Row in the `docs/index.html` endpoint table, README if the quickstart changes.

## Gotchas (verified against the OME docs, 2026-09)

- `POST /v1/vhosts` takes and returns a **top-level JSON array**. `RequestCreateVirtualHost` and
  `ResponseVirtualHost` are structs, so the body is wrong and the reply never decodes (A-01).
- Stats: the stream path is `/streams/{stream}` (plural); the SDK sends `/stream/` (A-02).
  Stats replies are wrapped in `response`, `ResponseStats` isn't — every field decodes as 0 (A-03).
- `:records` returns an **array**; `GetRecordingState` decodes into a single object and always
  errors (A-05).
- OME timestamps like `2021-08-31T23:44:44.789+0900` (no colon in the offset) do not parse into
  `time.Time`; one bad field fails the whole `Unmarshal` (A-06). Prefer `string` or a custom type.
- `validator.v2` `min=1` on `StreamKey` rejects SRT/MPEG-TS pushes, which have no stream key;
  `StopPush` validates the whole push body although OME only needs `{"id"}` (A-07).
- `SetDebug(true)` makes resty log every request **including the `Authorization` header** (A-10).
- The resty client has no timeout: a hung OME blocks the caller forever (A-11).
- `GET_*` builders are exported `var`s, so an importer can reassign them (A-16). Don't add more.
- Thumbnails aren't on the API port: they are served by the publisher
  (`http(s)://<host>:<port>/<app>/<stream>/thumb.<jpg|png>`), so `GetThumbnail` takes the full
  host, which resty joins with the base URL only when `host` is absolute.

## Conventions

- **TDD**: red first — a test that fails for the right reason, then the fix. A test that only
  pins existing behaviour must be shown to fail on a mutation.
- **Examples**: every public method gets a runnable `Example...` (package `ovenmedia_test`,
  `// Output:` block, `httptest` server) so pkg.go.dev shows it. Name it after something go/doc
  can attach it to (`ExampleBuildOven`, `ExampleRequestBodyPush`) — examples named after
  interface methods compile but pkg.go.dev drops them.
- Server behaviour comes from the OME docs or source (`AirenSoft/OvenMediaEngine`,
  `src/projects/api_server/controllers/v1/`), never from guesses. Quote the page in the PR.
- Keep `ovenMedia` unexported; the public surface is `IOvenMediaClient`. Adding a method to the
  interface is a breaking change for anyone who implements it (mocks) — say so in the PR.
- Match the surrounding style (`fmt.Sprintf` path builders, `json.Unmarshal` into a local `obj`),
  but no `fmt.Println` in library code and no empty `//` comment lines in new code.

## Definition of done

1. Test first, then the fix; an `Example...` for any public API change.
2. Doc comments, `README.md`, `docs/index.html`, `AUDIT.md` (mark the item fixed) and this file
   updated in the same change.
3. `go vet`, `go test -race`, `golangci-lint` green locally and in CI.
4. After the user tags: the release workflow succeeded and pkg.go.dev lists the version. If it
   404s: `curl -X POST https://pkg.go.dev/fetch/github.com/Allan-Nava/OvenMediaEngine-go-sdk@vX.Y.Z`.
