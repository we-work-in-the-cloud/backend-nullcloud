## Why

Users need a reliable, deterministic way to connect to provisioned databases. Currently, the Database model has no connection endpoint, forcing users to construct connection strings manually or guess at the infrastructure domain. By exposing a computed `endpoint` attribute (e.g., `db-{id}.db.nullcloud.internal:5432`), we make database connections transparent and consistent across Terraform, the UI, and documentation.

## What Changes

- Add `Endpoint` field to the `Database` model struct
- Compute endpoint in the API layer when a database is created (format: `db-{id}.db.nullcloud.internal:{port}`)
- Port is derived from the Engine: postgres → 5432, mysql/mariadb → 3306
- Expose endpoint in all Database API responses (list, get, create)
- Surface endpoint as a computed output in the Terraform provider
- Display endpoint in the UI's database details view
- Document the endpoint in examples and reference docs

## Capabilities

### New Capabilities
- `database-endpoint`: Computed connection string for provisioned databases. Provides a deterministic, user-friendly way to connect via the endpoint attribute (e.g., `db-{id}.db.nullcloud.internal:5432`).

### Modified Capabilities
- `database-management`: The Database resource now includes a read-only `endpoint` field that is set upon creation and available in all subsequent queries.

## Impact

- **Model**: Internal Go structs (internal/model/model.go)
- **API**: Database creation, retrieval, and listing handlers (internal/api/database.go)
- **Terraform Provider**: Database resource schema and computed attributes
- **UI**: Database resource details page (internal/ui)
- **Examples**: Database component examples in Terraform stacks
- **Documentation**: API docs, provider docs, and usage guides
