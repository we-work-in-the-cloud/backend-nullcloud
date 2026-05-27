# Agent Guidelines

## Repository layout

This is the NullCloud backend API server. It has a sibling repository that must be kept in sync:

- `backend-nullcloud/` — this repo (Go HTTP API)
- `terraform-provider-nullcloud/` — Terraform provider that wraps this API

Both live under the same parent directory. When making changes to one, always check whether the other needs a corresponding update.

## Keeping backend and Terraform provider in sync

The provider mirrors every model type defined here. Any structural change in the backend requires matching changes in the provider:

| Backend file | Provider equivalent |
|---|---|
| `internal/model/model.go` | `internal/client/client.go` (duplicate type definitions) |
| `internal/api/vpc.go` | `internal/provider/vpc_resource.go` |
| `internal/api/subnet.go` | `internal/provider/subnet_resource.go` |
| `internal/api/vsi.go` | `internal/provider/instance_resource.go` |

### Adding a field to a model

1. Add the field to the struct in `internal/model/model.go`.
2. Populate it in the relevant `internal/api/*.go` create handler.
3. Add the matching field to the client struct in `terraform-provider-nullcloud/internal/client/client.go`.
4. Add the attribute to the provider resource schema (`Computed: true` + `UseStateForUnknown` for server-generated fields).
5. Set the field in both `Create` and `Read` in the provider resource.
6. Update any tests that assert on the struct fields (backend `*_test.go` and provider `client_test.go`).

### Adding a new resource type

1. Create `internal/api/<name>.go` and wire routes in `internal/api/server.go`.
2. Write `internal/api/<name>_test.go` covering the full lifecycle.
3. Create matching files in the provider: `internal/client/` CRUD methods and `internal/provider/<name>_resource.go`.
4. Register the new resource in `internal/provider/provider.go`.

## CRN convention

Every resource must have a `CRN` field. The format is:

```
crn:nullcloud:<resource-type>:<id>
```

Resource type tokens: `vpc`, `subnet`, `instance`.

## Build and test

```sh
# backend
make test          # runs go test -v -cover ./...
make build         # cross-compiles for all platforms into dist/

# provider (from terraform-provider-nullcloud/)
go build ./...
go test ./...
```

Always run both before considering a change complete.

## JSON field naming

Go struct fields use `json:"snake_case"` tags. Match the existing style exactly — the provider client deserializes the same JSON the backend emits.
