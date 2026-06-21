# terraform-list-resources Specification

## Purpose

Define the Terraform provider list resources that allow users to query existing nullcloud infrastructure (VPCs, subnets, instances, load balancers, buckets, databases, and Kubernetes clusters) with optional filtering via the `terraform query` command.

## Requirements

### Requirement: VPC List Resource
The Terraform provider SHALL implement a list resource for VPCs that allows users to query existing VPCs with optional filtering.

#### Scenario: Query all VPCs
- **WHEN** user writes a list block for nullcloud_vpcs without filter conditions
- **THEN** the system returns all VPCs in the account as items in the results list

#### Scenario: Filter VPCs by region
- **WHEN** user includes a region filter in the config block
- **THEN** the system returns only VPCs in the specified region

#### Scenario: Access VPC attributes in results
- **WHEN** user iterates over results with for_each in HCL
- **THEN** each item exposes id, name, region, status, crn, and created_at attributes

### Requirement: Subnet List Resource
The Terraform provider SHALL implement a list resource for Subnets that allows users to query existing subnets with optional filtering.

#### Scenario: Query all subnets
- **WHEN** user writes a list block for nullcloud_subnets without filter conditions
- **THEN** the system returns all subnets in the account as items in the results list

#### Scenario: Filter subnets by VPC
- **WHEN** user includes a vpc_id filter in the config block
- **THEN** the system returns only subnets that belong to the specified VPC

#### Scenario: Filter subnets by zone
- **WHEN** user includes a zone filter in the config block
- **THEN** the system returns only subnets in the specified zone

#### Scenario: Access subnet attributes in results
- **WHEN** user iterates over results with for_each in HCL
- **THEN** each item exposes id, name, vpc_id, zone, status, crn, cidr_block, and created_at attributes

### Requirement: Instance List Resource
The Terraform provider SHALL implement a list resource for Instances that allows users to query existing compute instances with optional filtering.

#### Scenario: Query all instances
- **WHEN** user writes a list block for nullcloud_instances without filter conditions
- **THEN** the system returns all instances in the account as items in the results list

#### Scenario: Filter instances by subnet
- **WHEN** user includes a subnet_id filter in the config block
- **THEN** the system returns only instances deployed in the specified subnet

#### Scenario: Filter instances by status
- **WHEN** user includes a status filter in the config block (e.g., "running", "stopped")
- **THEN** the system returns only instances with the specified status

#### Scenario: Access instance attributes in results
- **WHEN** user iterates over results with for_each in HCL
- **THEN** each item exposes id, name, subnet_id, profile, image, status, crn, primary_ip, and created_at attributes

### Requirement: LoadBalancer List Resource
The Terraform provider SHALL implement a list resource for LoadBalancers that allows users to query existing load balancers with optional filtering.

#### Scenario: Query all load balancers
- **WHEN** user writes a list block for nullcloud_loadbalancers without filter conditions
- **THEN** the system returns all load balancers in the account as items in the results list

#### Scenario: Filter load balancers by protocol
- **WHEN** user includes a protocol filter in the config block (e.g., "HTTP", "HTTPS")
- **THEN** the system returns only load balancers using the specified protocol

#### Scenario: Access load balancer attributes in results
- **WHEN** user iterates over results with for_each in HCL
- **THEN** each item exposes id, name, status, crn, protocol, port, targets, and created_at attributes

### Requirement: Bucket List Resource
The Terraform provider SHALL implement a list resource for Buckets that allows users to query existing storage buckets with optional filtering.

#### Scenario: Query all buckets
- **WHEN** user writes a list block for nullcloud_buckets without filter conditions
- **THEN** the system returns all buckets in the account as items in the results list

#### Scenario: Filter buckets by region
- **WHEN** user includes a region filter in the config block
- **THEN** the system returns only buckets in the specified region

#### Scenario: Access bucket attributes in results
- **WHEN** user iterates over results with for_each in HCL
- **THEN** each item exposes id, name, status, crn, region, and created_at attributes

### Requirement: Database List Resource
The Terraform provider SHALL implement a list resource for Databases that allows users to query existing database instances with optional filtering.

#### Scenario: Query all databases
- **WHEN** user writes a list block for nullcloud_databases without filter conditions
- **THEN** the system returns all database instances in the account as items in the results list

#### Scenario: Filter databases by engine
- **WHEN** user includes an engine filter in the config block (e.g., "postgres", "mysql")
- **THEN** the system returns only databases using the specified engine

#### Scenario: Access database attributes in results
- **WHEN** user iterates over results with for_each in HCL
- **THEN** each item exposes id, name, engine, version, status, crn, plan, subnet_ids, endpoint, and created_at attributes

### Requirement: KubernetesCluster List Resource
The Terraform provider SHALL implement a list resource for Kubernetes Clusters that allows users to query existing clusters with optional filtering.

#### Scenario: Query all clusters
- **WHEN** user writes a list block for nullcloud_kubernetesclusters without filter conditions
- **THEN** the system returns all Kubernetes clusters in the account as items in the results list

#### Scenario: Filter clusters by version
- **WHEN** user includes a version filter in the config block
- **THEN** the system returns only clusters running the specified version

#### Scenario: Access cluster attributes in results
- **WHEN** user iterates over results with for_each in HCL
- **THEN** each item exposes id, name, version, status, crn, node_count, subnet_ids, and created_at attributes

### Requirement: Provider List Resource Registration
The Terraform provider SHALL register all list resources so they are discoverable and available via terraform query command.

#### Scenario: List resources appear in terraform schema
- **WHEN** user runs terraform init
- **THEN** all list resource types are available in the provider schema

#### Scenario: List resources work with terraform query
- **WHEN** user runs terraform query with a .tfquery.hcl file containing list blocks
- **THEN** the system executes the query and returns results in terraform JSON output format
