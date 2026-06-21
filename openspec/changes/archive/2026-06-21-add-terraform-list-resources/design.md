## Context

The Terraform framework (v1.9+) introduces the `list.ListResource` interface, enabling providers to expose resource discovery and filtering via the `terraform query` command. The nullcloud backend already implements robust list endpoints for all resource types (`GET /v1/vpcs`, `/v1/subnets`, etc.), returning paginated filtered results. The Terraform provider needs to implement list resources that wrap these backend capabilities. This design outlines the architecture for implementing list resources across all seven resource types (VPC, Subnet, Instance, LoadBalancer, Bucket, Database, KubernetesCluster).

**Current State:**
- Backend: List endpoints operational and tested (all resource types)
- Provider: No list resource implementations
- Framework: `list.ListResource` interface available in Go SDK

**Constraints:**
- List resources must implement three methods: `Metadata`, `Schema`, and `List`
- Schema must define config block attributes (filters) and results structure
- Results reuse existing resource attribute definitions for consistency
- Provider registration requires `ProviderWithListResources.ListResources()` method

## Goals / Non-Goals

**Goals:**
- Implement `list.ListResource` for all seven resource types in the Terraform provider
- Define intuitive, resource-specific filter parameters (region, status, engine, etc.)
- Reuse existing resource model definitions and attribute schemas
- Register all list resources in the provider for discovery via `terraform query`
- Enable end-to-end workflows: users can discover and filter resources without external tools

**Non-Goals:**
- Modify backend API endpoints or contracts (use existing list endpoints only)
- Alter existing resource or data source implementations
- Implement pagination in Terraform (return all results in a single call)
- Add complex nested filtering or boolean filter combinations
- Create list resource management UI (scope limited to HCL/terraform query support)

## Decisions

**Decision 1: Implementation Location (Provider vs Backend)**
- **Choice**: Implement list resources in the Terraform provider repository (`terraform-provider-nullcloud`), not the backend
- **Rationale**: List resources are Terraform abstractions that wrap backend list endpoints; they belong in the provider where resource/datasource implementations live. Backend remains unchanged.
- **Alternatives Considered**: Adding list-specific backend types (unnecessary—backend already exposes needed data), monolithic backend change (increases coupling)

**Decision 2: File Organization**
- **Choice**: Create separate files per list resource (vpc_list_resource.go, subnet_list_resource.go, etc.), following existing provider conventions
- **Rationale**: Mirrors resource/datasource file structure for consistency and maintainability; each file is ~80-150 lines
- **Alternatives Considered**: Single file for all list resources (harder to navigate), shared base class (Go doesn't have inheritance)

**Decision 3: Schema and Model Reuse**
- **Choice**: Reuse existing `vpcModel`, `subnetModel`, etc. struct types for results items; define config attributes separately
- **Rationale**: Avoids duplication, keeps model definitions authoritative in one place, ensures results match resource attribute structure
- **Alternatives Considered**: Define list-specific models (duplication), flatten results (loses structure)

**Decision 4: Filter Parameter Scope**
- **Choice**: Support 2-4 filters per resource type; focus on most common queries (region, status, engine, protocol, etc.)
- **Rationale**: Covers 80% of discovery use cases without overwhelming schema; aligns with backend query capabilities
- **Alternatives Considered**: All possible filters (schema bloat), hardcoded filters (inflexible), no filters (limits discoverability)

**Decision 5: Backend API Integration**
- **Choice**: Call existing `client.List*()` methods with filter parameters; no API changes to backend
- **Rationale**: Backend list endpoints already exist and are tested; minimal code change in provider
- **Alternatives Considered**: New API methods (unnecessary expansion), paginated iteration (increases token spend)

**Decision 6: Results Attribute Structure**
- **Choice**: Single `results` attribute containing a list of resource objects with all computed attributes
- **Rationale**: Mirrors Terraform datasource patterns; users can iterate with `for_each` idiom
- **Alternatives Considered**: Separate attributes per field (verbose), flattened results (loses structure), separate pages (pagination complexity)

## Risks / Trade-offs

**[Risk] Incomplete Filter Coverage**
- Users may want to filter by attributes not exposed in config (tags, labels, subnet_ids, etc.)
- **Mitigation**: Document filter capabilities clearly; design schema to be extensible (add filters in future versions without breaking changes)

**[Risk] Large Result Set Performance**
- If result set exceeds backend query timeout or memory limits, query fails
- **Mitigation**: Backend enforces reasonable defaults; start with conservative limits; backend team can tune if needed

**[Risk] Schema Maintenance**
- If existing resource schema changes, list resource schema must stay in sync
- **Mitigation**: Reuse model structs where possible; apply consistent naming; schema changes are infrequent

**[Trade-off] No Pagination**
- `terraform query` returns all results in one call; large result sets may be slow
- **Benefit**: Simpler Terraform user experience (no pagination logic in HCL)
- **Mitigation**: Document limitations; recommend filtering to narrow result sets

**[Trade-off] Read-Only**
- List resources don't support `terraform apply` (queries only)
- **Benefit**: Prevents accidental resource creation via list syntax; clear separation of concerns
- **Intentional**: Aligns with Terraform list resource design philosophy

## Migration Plan

**Phase 1: Provider Implementation**
1. Create list resource implementations for all seven resource types (parallel tasks)
2. Add `ListResources()` method to provider struct
3. Unit test each list resource independently (mock client responses)

**Phase 2: Integration & Testing**
1. Manual test with Terraform config containing list blocks
2. Verify `terraform query` command works end-to-end
3. Test edge cases (empty results, filter combinations, API errors)

**Phase 3: Release**
1. Bump provider version (minor version: new feature)
2. Update provider registry docs with list resource examples
3. Release on Terraform Registry

**Rollback Strategy:**
- If critical bug discovered: revert `provider.go` change (removes list resources from provider); list implementation files remain in repo but unused
- No state or data mutations, so safe to disable

## Open Questions

- Should we support tag-based filtering for resources that will have tags in future?
  - **Defer to**: Tasks/implementation phase; start without tag filters, add in follow-up if requested
- Any preferred naming conventions for filter parameters (e.g., `region_eq` vs `region`)?
  - **Defer to**: Spec validation; use simple parameter names for common filters
- Should list resources be added to the registry docs immediately or after v1.0 release?
  - **Defer to**: Release planning; recommend documentation update for beta release
