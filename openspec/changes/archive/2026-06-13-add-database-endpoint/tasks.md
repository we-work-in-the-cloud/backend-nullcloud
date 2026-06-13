## 1. Model Layer

- [x] 1.1 Add `Endpoint string` field to the `Database` struct in `internal/model/model.go`
- [x] 1.2 Add JSON tag `json:"endpoint"` to the new Endpoint field

## 2. API Layer

- [x] 2.1 Create a utility function `computeDatabaseEndpoint(id, engine string) string` in `internal/api/database.go`
- [x] 2.2 Map engine types to ports: postgres→5432, mysql→3306, mariadb→3306
- [x] 2.3 Update `createDatabase()` handler to compute and assign the endpoint before storing the Database
- [x] 2.4 Verify endpoint is included in the create response
- [x] 2.5 Verify endpoint is included in get and list responses (no changes needed to getDatabase/listDatabases handlers, as they serialize the struct directly)

## 3. Terraform Provider

- [x] 3.1 Update the Database resource schema in the Terraform provider to include an `endpoint` attribute
- [x] 3.2 Mark endpoint as `Computed: true` and `Required: false` (read-only)
- [x] 3.3 Test that endpoint appears in Terraform state after database creation
- [x] 3.4 Verify that attempting to set endpoint in Terraform config results in an error

## 4. UI

- [x] 4.1 Identify the database resource details component in `internal/ui/src`
- [x] 4.2 Add endpoint field to the database display (e.g., connection details section)
- [x] 4.3 Style/format the endpoint for readability (e.g., monospace font)

## 5. Examples

- [x] 5.1 Update database component example in `examples/stacks/components/database/` to reference the endpoint in outputs or locals
- [x] 5.2 Update network component example (if applicable) to show how endpoint relates to networking

## 6. Documentation

- [x] 6.1 Update API documentation to document the endpoint field in Database responses
- [x] 6.2 Update Terraform provider documentation to document the endpoint attribute
- [x] 6.3 Add usage example showing how to use the endpoint to connect to a database

## 7. Testing

- [x] 7.1 Add/update tests for `computeDatabaseEndpoint()` utility function (test all engine types)
- [x] 7.2 Add/update API tests to verify endpoint is returned in create, get, and list responses
- [x] 7.3 Test that different engine types produce correct port numbers
- [ ] 7.4 Manual test: Create a database via UI/API/Terraform and verify endpoint is displayed correctly
