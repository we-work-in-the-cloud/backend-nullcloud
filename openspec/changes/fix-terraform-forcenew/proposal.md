## Why

Resource updates in the Terraform provider unnecessarily destroy and recreate resources when they should be updated in place. Mutable operational attributes like names, targets, and plans are incorrectly marked as `ForceNew`, causing unnecessary resource destruction, downtime, and cost overhead. This reduces provider usability and conflicts with cloud API capabilities that support in-place updates for these attributes.

## What Changes

**Backend API (nullcloud/backend-nullcloud):**
- Add plan parameter support to database PATCH endpoint to allow in-place plan scaling
- Add targets parameter support to load balancer PATCH endpoint to allow target list updates

**Terraform Provider (nullcloud/terraform-provider-nullcloud):**
- Remove `ForceNew` from name fields across all resources (VPC, subnet, load balancer, database, etc.)
- Remove `ForceNew` from load balancer targets attribute
- Remove `ForceNew` from database plan attribute
- Implement/update Update methods to handle the newly-mutable attributes
- Establish a principle: metadata and operational attributes update in place; only structural attributes trigger destruction

## Capabilities

### Modified Capabilities
- `terraform-resource-updates`: Update behavior for resource attributes — mutable operational attributes (names, targets, plans) should update in place rather than trigger destruction

## Impact

- **Affected resources**: VPC, subnet, load balancer, database resources
- **User experience**: Terraform plans will now show updates instead of destruction for these attributes
- **Terraform provider**: Schema definitions for affected resources
