## Context

The Terraform provider currently marks mutable metadata and operational attributes as `RequiresReplace()` (equivalent to ForceNew), causing resource destruction and recreation when they should be updated in place. Additionally, the backend API does not yet support in-place updates for plan and targets attributes.

**Terraform Provider** issues:
- Resource schema files under `internal/provider/` mark these with `RequiresReplace()`:
  - `vpc_resource.go`: VPC name
  - `subnet_resource.go`: Subnet name
  - `database_resource.go`: Database name and plan
  - `loadbalancer_resource.go`: Load balancer name and targets

**Backend API** gaps:
- `database.go`: updateDatabase() only accepts `name`, not `plan`
- `loadbalancer.go`: updateLoadBalancer() only accepts `name`, not `targets`

The plan modifiers use Terraform's schema directives (`stringplanmodifier.RequiresReplace()`, `listplanmodifier.RequiresReplace()`) to control this behavior.

## Goals / Non-Goals

**Goals:**
- Remove `RequiresReplace()` from metadata fields (all name attributes) across resources
- Remove `RequiresReplace()` from operational attributes (targets, plan)
- Enable in-place updates via Terraform's Update methods for these attributes
- Ensure cloud API client supports these in-place updates

**Non-Goals:**
- Change behavior for structural attributes (region, zone, CIDR, protocol, port, engine) — these remain immutable
- Add new Update logic beyond what the provider already supports
- Change how computed or read-only attributes are handled

## Decisions

### Backend API Changes

**Extend database PATCH endpoint** (database.go, updateDatabase function):
- Add `plan` to the request struct (currently only accepts `name`)
- Validate plan is one of: small, medium, large
- Call a new store method `UpdateDatabasePlan()` to persist the change
- Return updated database in response

**Extend load balancer PATCH endpoint** (loadbalancer.go, updateLoadBalancer function):
- Add `targets` to the request struct (currently only accepts `name`)
- Validate targets list: each target must have valid type (cluster or vsi) and existing ID
- Call a new store method `UpdateLoadBalancerTargets()` to persist the change
- Return updated load balancer in response

**Backend store layer:**
- Implement `UpdateDatabasePlan(ctx, token, id, plan)` in store
- Implement `UpdateLoadBalancerTargets(ctx, token, id, targets)` in store

### Terraform Provider Changes

**Remove plan modifiers from name fields:**
- VPC name (vpc_resource.go, line 47-51)
- Subnet name (subnet_resource.go, line 49-53)
- Load balancer name (loadbalancer_resource.go, line 63-68)
- Database name (database_resource.go, line 52-57)

Simply delete the `PlanModifiers` array containing `stringplanmodifier.RequiresReplace()` from these attributes.

**Remove plan modifiers from mutable operational attributes:**
- Load balancer targets (loadbalancer_resource.go, line 84-102): remove `listplanmodifier.RequiresReplace()`
- Database plan (database_resource.go, line 73-78): remove `stringplanmodifier.RequiresReplace()`

**Update the Update() methods:**
- VPC and subnet: existing Update implementation should handle name changes
- Database: Update method needs to handle both name and plan changes
- Load balancer: currently a no-op; implement to handle name and targets changes

**Rationale:**
- Backend must support the operations before provider can use them
- Schema changes are straightforward — just remove RequiresReplace modifiers
- Terraform will automatically use Update instead of Replace once modifiers are removed
- This follows Terraform provider best practices: only use ForceNew for truly immutable attributes

## Risks / Trade-offs

**[Risk] Incomplete Update implementations** → Verify that existing Update methods in each resource handle the newly-mutable attributes correctly. Load balancer's Update method is currently a no-op and will need implementation.

**[Risk] Cloud API support** → Assumes the nullcloud API actually supports these in-place operations. This was validated during analysis, but verify the API client methods (UpdateVPC, UpdateDatabase, UpdateLoadBalancer, etc.) exist and work correctly.

**[Risk] State management** → Ensure the Terraform state is properly synced after Update operations. Some resources may need additional Read logic to refresh state.

**[Mitigation]** Write acceptance tests for each resource to verify update behavior works as expected.
