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
| `internal/api/vpc.go` | `internal/provider/vpc_resource.go`, `internal/provider/vpc_data_source.go` |
| `internal/api/subnet.go` | `internal/provider/subnet_resource.go`, `internal/provider/subnet_data_source.go` |
| `internal/api/vsi.go` | `internal/provider/instance_resource.go`, `internal/provider/instance_data_source.go`, `internal/provider/instance_action.go` |

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
5. Update `README.md` in this repo — add the new resource to the Resources list.
6. Update `terraform-provider-nullcloud/README.md` — add rows to the Resources and Data Sources tables.
7. Update the sync table in this `AGENTS.md` — add a row mapping the new `internal/api/<name>.go` to its provider equivalents.
8. Update `terraform-provider-nullcloud/AGENTS.md` — add the new files to the internal structure diagram and the sync table.

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

## Keeping docs current

Whenever you make a structural change, update these files before closing the task:

| Change | Files to update |
|---|---|
| New resource type | `README.md` (Resources list), `AGENTS.md` sync table, `terraform-provider-nullcloud/README.md` (Resources + Data Sources tables), `terraform-provider-nullcloud/AGENTS.md` (internal structure + sync table) |
| New field on a model | No README change needed; update `AGENTS.md` sync table only if the mapping changes |
| Renamed / removed resource | Same files as "new resource type" — remove or rename the old entries |

The goal: a reader of either `README.md` or `AGENTS.md` should be able to understand the current state of the project without reading the code.

## JSON field naming

Go struct fields use `json:"snake_case"` tags. Match the existing style exactly — the provider client deserializes the same JSON the backend emits.
