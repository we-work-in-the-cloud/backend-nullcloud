## 1. Backend API: Database Plan Updates

- [x] 1.1 Extend database.go updateDatabase() to accept plan parameter
- [x] 1.2 Add plan validation (small, medium, large) in updateDatabase()
- [x] 1.3 Implement store.UpdateDatabasePlan() method
- [x] 1.4 Update updateDatabase() to call UpdateDatabasePlan()

## 2. Backend API: Load Balancer Targets Updates

- [x] 2.1 Extend loadbalancer.go updateLoadBalancer() to accept targets parameter
- [x] 2.2 Add targets validation (type and ID existence checks) in updateLoadBalancer()
- [x] 2.3 Implement store.UpdateLoadBalancerTargets() method
- [x] 2.4 Update updateLoadBalancer() to call UpdateLoadBalancerTargets()

## 3. Backend Tests

- [x] 3.1 Add acceptance test for database plan update
- [x] 3.2 Add acceptance test for load balancer targets update

## 4. Terraform Provider: Remove ForceNew from Name Fields

- [x] 4.1 Remove RequiresReplace from VPC name in vpc_resource.go (line 49-51)
- [x] 4.2 Remove RequiresReplace from Subnet name in subnet_resource.go (line 51-53)
- [x] 4.3 Remove RequiresReplace from Database name in database_resource.go (line 55-57)
- [x] 4.4 Remove RequiresReplace from Load Balancer name in loadbalancer_resource.go (line 66-68)

## 5. Terraform Provider: Remove ForceNew from Operational Attributes

- [x] 5.1 Remove RequiresReplace from Database plan in database_resource.go (line 76-78)
- [x] 5.2 Remove RequiresReplace from Load Balancer targets in loadbalancer_resource.go (line 87-89)

## 6. Terraform Provider: Update Methods

- [x] 6.1 Verify VPC Update method supports name changes
- [x] 6.2 Verify Subnet Update method supports name changes
- [x] 6.3 Implement Database Update method to handle name and plan changes
- [x] 6.4 Implement Load Balancer Update method for name and targets changes

## 7. Integration Tests

- [x] 7.1 Test VPC name update with Terraform apply (acceptance test added)
- [x] 7.2 Test Subnet name update with Terraform apply (acceptance test added)
- [x] 7.3 Test Database name update with Terraform apply (acceptance test added)
- [ ] 7.4 Test Database plan update with Terraform apply (ready to run)
- [x] 7.5 Test Load Balancer name update with Terraform apply (acceptance test added)
- [ ] 7.6 Test Load Balancer targets update with Terraform apply (ready to run)
- [ ] 7.7 Run acceptance tests to ensure no regressions
