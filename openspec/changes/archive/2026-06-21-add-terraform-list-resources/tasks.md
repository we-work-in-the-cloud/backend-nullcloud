## 1. Setup & Infrastructure

- [x] 1.1 Create base list resource scaffolding (abstract methods and helpers)
- [x] 1.2 Add ListResources() method to provider.go with empty list
- [x] 1.3 Verify provider compiles with ListResources interface

## 2. VPC List Resource

- [x] 2.1 Create internal/provider/vpc_list_resource.go with Metadata, Schema, List methods
- [x] 2.2 Define vpc list config model with region filter parameter
- [x] 2.3 Implement List method to call client.ListVPCs with filters
- [x] 2.4 Map VPC API responses to results attribute
- [x] 2.5 Register NewVPCListResource in provider.ListResources()
- [x] 2.6 Unit test VPC list resource (all scenarios from spec)

## 3. Subnet List Resource

- [x] 3.1 Create internal/provider/subnet_list_resource.go with Metadata, Schema, List methods
- [x] 3.2 Define subnet list config model with vpc_id and zone filter parameters
- [x] 3.3 Implement List method to call client.ListSubnets with filters
- [x] 3.4 Map Subnet API responses to results attribute
- [x] 3.5 Register NewSubnetListResource in provider.ListResources()
- [x] 3.6 Unit test Subnet list resource (all scenarios from spec)

## 4. Instance List Resource

- [x] 4.1 Create internal/provider/instance_list_resource.go with Metadata, Schema, List methods
- [x] 4.2 Define instance list config model with subnet_id and status filter parameters
- [x] 4.3 Implement List method to call client.ListInstances with filters
- [x] 4.4 Map Instance API responses to results attribute
- [x] 4.5 Register NewInstanceListResource in provider.ListResources()
- [x] 4.6 Unit test Instance list resource (all scenarios from spec)

## 5. LoadBalancer List Resource

- [x] 5.1 Create internal/provider/loadbalancer_list_resource.go with Metadata, Schema, List methods
- [x] 5.2 Define LoadBalancer list config model with protocol filter parameter
- [x] 5.3 Implement List method to call client.ListLoadBalancers with filters
- [x] 5.4 Map LoadBalancer API responses to results attribute
- [x] 5.5 Register NewLoadBalancerListResource in provider.ListResources()
- [x] 5.6 Unit test LoadBalancer list resource (all scenarios from spec)

## 6. Bucket List Resource

- [x] 6.1 Create internal/provider/bucket_list_resource.go with Metadata, Schema, List methods
- [x] 6.2 Define bucket list config model with region filter parameter
- [x] 6.3 Implement List method to call client.ListBuckets with filters
- [x] 6.4 Map Bucket API responses to results attribute
- [x] 6.5 Register NewBucketListResource in provider.ListResources()
- [x] 6.6 Unit test Bucket list resource (all scenarios from spec)

## 7. Database List Resource

- [x] 7.1 Create internal/provider/database_list_resource.go with Metadata, Schema, List methods
- [x] 7.2 Define database list config model with engine filter parameter
- [x] 7.3 Implement List method to call client.ListDatabases with filters
- [x] 7.4 Map Database API responses to results attribute
- [x] 7.5 Register NewDatabaseListResource in provider.ListResources()
- [x] 7.6 Unit test Database list resource (all scenarios from spec)

## 8. KubernetesCluster List Resource

- [x] 8.1 Create internal/provider/cluster_list_resource.go with Metadata, Schema, List methods
- [x] 8.2 Define KubernetesCluster list config model with version filter parameter
- [x] 8.3 Implement List method to call client.ListKubernetesClusters with filters
- [x] 8.4 Map KubernetesCluster API responses to results attribute
- [x] 8.5 Register NewKubernetesClusterListResource in provider.ListResources()
- [x] 8.6 Unit test KubernetesCluster list resource (all scenarios from spec)

## 9. Integration & Validation

- [x] 9.1 Verify all list resources register in provider (no compilation errors)
- [x] 9.2 Test provider schema includes all list resource types
- [x] 9.3 Manual test with HCL: list block without filters (returns all results)
- [x] 9.4 Manual test with HCL: list block with filters (returns filtered results)
- [x] 9.5 Verify terraform query command works with .tfquery.hcl files
- [x] 9.6 Test edge cases: empty result sets, API errors, invalid filter values

## 10. Documentation & Release

- [ ] 10.1 Add list resource examples to provider documentation
- [ ] 10.2 Document filter parameters for each list resource
- [ ] 10.3 Update CHANGELOG with list resource additions
- [ ] 10.4 Bump provider version (minor version bump for new feature)
- [ ] 10.5 Tag release and verify Terraform Registry sync
