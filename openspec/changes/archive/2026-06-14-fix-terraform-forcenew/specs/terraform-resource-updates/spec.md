## ADDED Requirements

### Requirement: Database plan can be updated in place
The API SHALL support updating a database instance's plan (size tier) without destroying the database. A PATCH request with the plan parameter SHALL scale the database up or down.

#### Scenario: Scale database to larger plan
- **WHEN** user sends PATCH /database/{id} with plan=large
- **THEN** system updates the database plan and returns the updated database

#### Scenario: Scale database to smaller plan
- **WHEN** user sends PATCH /database/{id} with plan=small
- **THEN** system updates the database plan and returns the updated database

### Requirement: Load balancer targets can be updated in place
The API SHALL support updating a load balancer's target list without destroying the load balancer. A PATCH request with the targets parameter SHALL add or remove targets.

#### Scenario: Add targets to load balancer
- **WHEN** user sends PATCH /loadbalancer/{id} with targets list
- **THEN** system updates the load balancer targets and returns the updated load balancer

#### Scenario: Remove targets from load balancer
- **WHEN** user sends PATCH /loadbalancer/{id} with empty targets list
- **THEN** system updates the load balancer targets (removing all targets) and returns the updated load balancer

## MODIFIED Requirements

### Requirement: Name attribute updates do not trigger resource destruction
Mutable metadata attributes like names SHALL be updatable in place without triggering resource destruction. Changing a resource name SHALL result in an in-place update, not a destroy/recreate cycle.

#### Scenario: VPC name change
- **WHEN** user changes a VPC name in Terraform configuration
- **THEN** Terraform plan shows an in-place update, not a destroy/recreate

#### Scenario: Subnet name change
- **WHEN** user changes a subnet name in Terraform configuration
- **THEN** Terraform plan shows an in-place update, not a destroy/recreate

#### Scenario: Load balancer name change
- **WHEN** user changes a load balancer name in Terraform configuration
- **THEN** Terraform plan shows an in-place update, not a destroy/recreate

### Requirement: Load balancer targets attribute updates do not trigger resource destruction
The targets attribute on a load balancer SHALL be updatable in place. Adding or removing targets SHALL result in an in-place membership update, not a destroy/recreate of the load balancer.

#### Scenario: Add target to load balancer
- **WHEN** user adds a target to a load balancer's targets list
- **THEN** Terraform plan shows an in-place update to the load balancer

#### Scenario: Remove target from load balancer
- **WHEN** user removes a target from a load balancer's targets list
- **THEN** Terraform plan shows an in-place update to the load balancer

### Requirement: Database plan attribute updates do not trigger resource destruction
The plan attribute on a database resource (representing instance size/tier) SHALL be updatable in place. Scaling a database plan up or down SHALL result in an in-place update, not a destroy/recreate.

#### Scenario: Scale database plan up
- **WHEN** user changes a database plan to a larger tier
- **THEN** Terraform plan shows an in-place update to the database plan

#### Scenario: Scale database plan down
- **WHEN** user changes a database plan to a smaller tier
- **THEN** Terraform plan shows an in-place update to the database plan
