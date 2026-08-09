# Controller SDK Go - Agent Instructions

> **When you change code patterns, conventions, or architecture, update this file
> in the same commit.** AGENTS.md is the single source of truth for project
> conventions - if it's stale, the next agent will follow wrong patterns.

## Project Overview

Go SDK for interacting with the Drycc controller API. Part of a three-repo stack:

```
controller (Python/Django API server)  -->  controller-sdk-go (this repo, Go SDK)  -->  workflow-cli (Go CLI)
```

The controller defines the REST API. This SDK provides typed Go structs and HTTP client methods for that API. The workflow-cli consumes the SDK to build the user-facing CLI.

- **Module**: `github.com/drycc/controller-sdk-go`
- **Go version**: 1.26 (see `go.mod`)
- **Go API version**: `2.3` (constant `APIVersion` in `drycc.go`)

## Directory Structure

```
controller-sdk-go/
  drycc.go              # Client struct, New(), APIVersion constant
  http.go               # Client.Do(), Request(), LimitedRequest() - HTTP transport layer
  errors.go             # Error types (ErrNotFound, ErrUnprocessable, ErrConflict, etc.)
  utils.go              # Helper functions
  api/                  # Type definitions with json tags (one file per domain)
    apps.go             # App, Apps, AppCreateRequest, ...
    config.go           # Config, Lifecycle, Healthcheck, ContainerProbe, ...
    gateways.go         # Gateway, Gateways, GatewayUpdateRequest, GatewayInfo, ...
    routes.go           # Route, Routes, RouteUpdateRequest, RouteInfo, ...
    addons.go           # AddonClass, AddonInstance, AddonInstanceUpsertRequest, ...
    workspaces.go       # Workspace, WorkspaceMember, WorkspaceInvitation, ...
    ps.go               # Pods, PodsList, ContainerState, ...
    ...
  addons/               # Domain package: HTTP methods for addons (List, Get, Upsert, Delete, ...)
  apps/                 # Domain package: HTTP methods for apps
  config/               # Domain package: HTTP methods for config
  gateways/             # Domain package: HTTP methods for gateways (List, Apply, Info, Delete)
  routes/               # Domain package: HTTP methods for routes (List, Apply, Info, Delete)
  workspaces/           # Domain package: HTTP methods for workspaces
    members/            # Sub-package: workspace members
    invitations/        # Sub-package: workspace invitations
  auth/                 # Authentication (Login)
  tokens/               # Token management
  ...                   # Other domain packages: builds, certs, domains, events, etc.
  pkg/                  # Shared utilities (e.g. pkg/time)
  Makefile              # Build/test commands (uses podman dev image)
  vendor/               # Vendored dependencies
```

### Layer Separation

- **`api/` package**: Pure type definitions with `json` struct tags. No HTTP logic. One file per domain (e.g. `api/gateways.go`, `api/routes.go`).
- **Domain packages** (e.g. `gateways/`, `routes/`, `addons/`): HTTP client methods that use `drycc.Client` to call controller endpoints and unmarshal into `api/` types. One file per domain (e.g. `gateways/gateways.go`).
- **Root package** (`drycc.go`, `http.go`, `errors.go`): `Client` struct, HTTP transport, error handling.

## Naming Convention

### JSON Struct Tags - snake_case ALWAYS

All `json` struct tags use **snake_case**. K8s camelCase field names are NOT used in json tags. The controller API converts K8s camelCase to snake_case via `CamelCaseJSONField` on the server side.

```go
// CORRECT - snake_case json tags
type ContainerProbe struct {
    InitialDelaySeconds int              `json:"initial_delay_seconds"`
    TimeoutSeconds      int              `json:"timeout_seconds"`
    HTTPGet             *HTTPGetAction   `json:"http_get,omitempty"`
    TCPSocket           *TCPSocketAction `json:"tcp_socket,omitempty"`
}

type Lifecycle struct {
    PostStart  **LifecycleHandler `json:"post_start,omitempty"`
    PreStop    **LifecycleHandler `json:"pre_stop,omitempty"`
    StopSignal string             `json:"stop_signal,omitempty"`
}

// WRONG - never use K8s camelCase in json tags
type ContainerProbe struct {
    InitialDelaySeconds int `json:"initialDelaySeconds"`  // NO
    HTTPGet             *HTTPGetAction `json:"httpGet"`    // NO
}
```

K8s camelCase names and their SDK snake_case equivalents:

| K8s (camelCase) | SDK json tag (snake_case) |
|---|---|
| `livenessProbe` | `liveness_probe` |
| `httpGet` | `http_get` |
| `postStart` | `post_start` |
| `preStop` | `pre_stop` |
| `tcpSocket` | `tcp_socket` |
| `initialDelaySeconds` | `initial_delay_seconds` |
| `parentRefs` | `parent_refs` |
| `backendRefs` | `backend_refs` |
| `terminationGracePeriod` | `termination_grace_period` |

### Go Field Names - PascalCase

Go struct field names use standard PascalCase (e.g. `HTTPGet`, `TCPSocket`, `InitialDelaySeconds`). Only the json tags are snake_case.

## API Contract

The **controller** (Python/Django) is the source of truth for the API. SDK types must match the controller's API response format exactly.

- When the controller changes a field name, the SDK `json` tag must be updated to match.
- When the controller adds a new field, add a corresponding struct field with the matching snake_case json tag.
- When the controller removes a field, remove it from the SDK struct.
- The SDK never invents field names - it mirrors the controller's serializer output.

## Type Conventions

### Response Types: `XxxInfo` or `Xxx`

Single-resource GET responses use `XxxInfo` for resources that have a distinct "info" shape, or just `Xxx` when the list item and single item share the same type:

```go
type GatewayInfo struct { ... }   // GET /v2/apps/<app>/gateways/<name>/
type RouteInfo struct { ... }     // GET /v2/apps/<app>/routes/<name>/
type AddonInstance struct { ... } // GET /v2/apps/<app>/addons/<name>/ (same as list item)
type Workspace struct { ... }     // GET /v2/workspaces/<name>
```

### Request Types: `XxxUpdateRequest` / `XxxUpsertRequest` / `XxxCreateRequest`

```go
type GatewayUpdateRequest struct { ... }             // PUT body
type RouteUpdateRequest struct { ... }               // PUT body
type AddonInstanceUpsertRequest struct { ... }       // PUT body (create or update)
type AppCreateRequest struct { ... }                 // POST body
type WorkspaceCreateRequest struct { ... }           // POST body
type WorkspaceUpdateRequest struct { ... }           // PATCH body
```

### Collection Types: `Xxxs` with sort methods

Collection types are named `Xxxs` (plural of the item type) and implement `sort.Interface` (`Len()`, `Swap()`, `Less()`):

```go
type Gateways []Gateway

func (g Gateways) Len() int           { return len(g) }
func (g Gateways) Swap(i, j int)      { g[i], g[j] = g[j], g[i] }
func (g Gateways) Less(i, j int) bool { return g[i].Name < g[j].Name }
```

Same pattern for `Routes`, `Apps`, `AddonClasses`, `AddonInstances`, `Workspaces`, `WorkspaceMembers`, `PodsList`, etc.

### Map Types for Flexible K8s Data

When a field carries arbitrary K8s rule/spec data that varies by version, use `map[string]any`:

```go
type RouteRule map[string]any           // Flexible K8s HTTPRoute rule data
type AddonConnection map[string]any     // Decoded Secret connection data
type ConfigTags map[string]any          // Tag key-values
```

### Field Tags

- `omitempty` on optional fields (most fields in response types)
- Explicit json tags on ALL fields - never rely on Go's default field name
- Required fields in request types omit `omitempty` (e.g. `Kind`, `Plan` in `AddonInstanceUpsertRequest`)
- Use `*bool` for optional boolean fields that need to distinguish zero-value from unset (e.g. `Routable *bool`)

### Pointer Types for Optional Nested Structs

Nested structs that may be absent use double pointers when they represent "set or unset" semantics (common for K8s lifecycle/probe fields):

```go
type Healthcheck struct {
    StartupProbe   **ContainerProbe `json:"startup_probe,omitempty"`
    LivenessProbe  **ContainerProbe `json:"liveness_probe,omitempty"`
    ReadinessProbe **ContainerProbe `json:"readiness_probe,omitempty"`
}
```

## Domain Package Pattern

Each domain package follows the same structure. Example from `gateways/gateways.go`:

```go
package gateways

import (
    "encoding/json"
    "fmt"

    drycc "github.com/drycc/controller-sdk-go"
    "github.com/drycc/controller-sdk-go/api"
)

// List returns gateways for an app.
func List(c *drycc.Client, appID string, results int) (api.Gateways, int, error) {
    u := fmt.Sprintf("/v2/apps/%s/gateways/", appID)
    body, count, reqErr := c.LimitedRequest(u, results)
    if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
        return []api.Gateway{}, -1, reqErr
    }
    var gateways []api.Gateway
    if err := json.Unmarshal([]byte(body), &gateways); err != nil {
        return []api.Gateway{}, -1, err
    }
    return gateways, count, reqErr
}

// Apply creates or updates a gateway (PUT).
func Apply(c *drycc.Client, appID string, req api.GatewayUpdateRequest) (api.GatewayInfo, error) { ... }

// Info retrieves a single gateway (GET).
func Info(c *drycc.Client, appID string, name string) (api.GatewayInfo, error) { ... }

// Delete removes a gateway (DELETE).
func Delete(c *drycc.Client, appID string, name string) error { ... }
```

### Standard Method Signatures

```go
// List - paginated collection
func List(c *drycc.Client, appID string, results int) (api.Xxxs, int, error)

// Get/Info - single resource
func Get(c *drycc.Client, appID string, name string) (api.Xxx, error)

// Apply/Upsert - create or update (PUT)
func Apply(c *drycc.Client, appID string, req api.XxxUpdateRequest) (api.XxxInfo, error)

// Delete - remove resource
func Delete(c *drycc.Client, appID string, name string) error
```

### Error Handling Pattern

All domain methods follow this error handling convention:

```go
if reqErr != nil && !drycc.IsErrAPIMismatch(reqErr) {
    return zeroValue, -1, reqErr
}
// API mismatch errors are non-fatal - return data along with the error
```

`ErrAPIMismatch` is returned when the controller's API version differs from the SDK's. It's non-fatal (data may still be usable), so domain methods return it alongside the data rather than discarding results.

## Testing

Each domain package has `*_test.go` files (e.g. `gateways/gateways_test.go`, `addons/addons_test.go`).

### Test Pattern

Tests use a `fakeHTTPServer` that mocks HTTP responses:

```go
type fakeHTTPServer struct{}

func (fakeHTTPServer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
    res.Header().Add("DRYCC_API_VERSION", drycc.APIVersion)

    if req.URL.Path == "/v2/apps/example-go/gateways/" && req.Method == "GET" {
        res.Write([]byte(gatewaysFixture))
        return
    }
    // ... more route matching
}

func TestGatewaysList(t *testing.T) {
    t.Parallel()

    handler := fakeHTTPServer{}
    server := httptest.NewServer(handler)
    defer server.Close()

    client, err := drycc.New(false, server.URL, "abc")
    if err != nil {
        t.Fatal(err)
    }

    actual, _, err := List(client, "example-go", 100)
    if err != nil {
        t.Fatal(err)
    }

    if !reflect.DeepEqual(expected, actual) {
        t.Error(fmt.Errorf("Expected %v, Got %v", expected, actual))
    }
}
```

### Mock JSON Conventions

- Mock JSON fixtures must use **snake_case** keys matching the API format (same as json struct tags)
- Fixtures are defined as `const` strings at the top of the test file
- For PUT/POST tests, define both the expected request body and the mock response
- Use `t.Parallel()` on all tests
- Compare with `reflect.DeepEqual` for struct comparison

## Build and Test

### Local Development

```bash
go build ./...          # Build all packages
go test ./...           # Run all tests
go vet ./...            # Run go vet
```

### Makefile (CI - uses podman dev image)

```bash
make bootstrap      # go mod vendor
make build          # go build ./...
make test-style     # lint
make test           # build + test-style + go test -race -cover
```

### CI Pipeline

Woodpecker CI (`.woodpecker/test-linux.yml`): runs `make bootstrap test` then uploads coverage to codecov.

## Dependency Management

- No `replace` directives in `go.mod`.
- Dependencies are vendored into `vendor/`.
- Published by commit hash - downstream consumers (workflow-cli) pin to a specific commit.
- **Update workflow**: When the controller API changes:
  1. Update types in `api/` package (add/modify json tags to match controller response)
  2. Update domain packages if new endpoints are added
  3. Update/add tests with matching mock JSON
  4. Run `go test ./...` to verify
  5. Commit and push
  6. Update workflow-cli's `go.mod` to pin the new commit hash

### Key Dependencies

| Dependency | Purpose |
|---|---|
| `golang.org/x/net` | HTTP utilities |
| `gopkg.in/yaml.v3` | YAML parsing |
| `github.com/stretchr/testify` | Test assertions |

## Git Conventions

- Commit message format: `type(scope): subject` (same as workflow-cli)
- One commit per feature per repo
- Example: `feat(gateways): add gateway list and apply methods`

## Key Patterns Summary

1. **Types in `api/`, methods in domain packages** - never mix HTTP logic into `api/` types.
2. **snake_case json tags always** - the controller converts K8s camelCase to snake_case.
3. **Controller is source of truth** - SDK mirrors the controller's API format exactly.
4. **`ErrAPIMismatch` is non-fatal** - return data alongside the error, don't discard results.
5. **Collection types implement `sort.Interface`** - `Len()`, `Swap()`, `Less()`.
6. **`LimitedRequest` for list endpoints** - handles pagination via `?limit=N` query param.
7. **`fakeHTTPServer` for tests** - mock server returns canned JSON with snake_case keys.
