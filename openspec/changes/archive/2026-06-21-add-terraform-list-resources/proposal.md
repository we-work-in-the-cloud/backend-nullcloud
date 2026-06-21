## Why

The Terraform plugin framework (v1.9+) introduces the `list` block type (via `terraform query` command), enabling providers to expose resource discovery and filtering capabilities. The nullcloud backend already implements list endpoints for all resource types (`GET /v1/vpcs`, `/v1/subnets`, etc.), returning paginated filtered results. The Terraform provider currently has resource and data source implementations but no list resources. This change enables the provider to surface these backend list capabilities through the new Terraform list resource abstraction, allowing users to discover resources via `terraform query` without external tooling.

## What Changes

- Terraform provider adds list resource implementations for all seven resource types (VPC, Subnet, Instance, LoadBalancer, Bucket, Database, KubernetesCluster)
- Each list resource accepts a `config` block with filter parameters (region, status, engine, etc.) and returns matching resources via a `results` attribute
- Provider extends implementation to register all list resources via the `ProviderWithListResources` interface
- No changes to backend API contracts, resource definitions, or data source implementations—list resources wrap existing list endpoints

## Capabilities

### New Capabilities
- `terraform-list-resources`: Add `list` block support for querying existing nullcloud resources with filters (region, status, and resource-specific criteria). Enables users to discover resources via `terraform query` without external tooling.

### Modified Capabilities
<!-- None - backend list endpoints already exist; no API changes required -->

## Impact

**Terraform Provider Changes:**
- New files: `vpc_list_resource.go`, `subnet_list_resource.go`, `instance_list_resource.go`, `loadbalancer_list_resource.go`, `bucket_list_resource.go`, `database_list_resource.go`, `cluster_list_resource.go` in `internal/provider/`
- Modified: `internal/provider/provider.go` (add `ListResources()` method to register all list resources)
- No changes to existing resource or data source implementations

**Backend API Usage:**
- Leverages existing list endpoints (`GET /v1/vpcs`, `/v1/subnets`, `/v1/instances`, etc.)
- No new API endpoints required; list endpoints already fully implemented and tested

**Dependencies:**
- Requires Terraform 1.9+ for `terraform query` command support
- Uses existing `client.List*` methods already available in the provider's API client
