package store

import (
	"context"
	"fmt"
	"sync"

	"github.com/we-work-in-the-cloud/nullcloud/backend/internal/model"
)

type MemoryStore struct {
	mu       sync.RWMutex
	vpcs     map[string]map[string]model.VPC
	subnets  map[string]map[string]model.Subnet
	vsis     map[string]map[string]model.VSI
	lbs      map[string]map[string]model.LoadBalancer
	buckets  map[string]map[string]model.Bucket
	dbs      map[string]map[string]model.Database
	clusters map[string]map[string]model.KubernetesCluster
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		vpcs:     make(map[string]map[string]model.VPC),
		subnets:  make(map[string]map[string]model.Subnet),
		vsis:     make(map[string]map[string]model.VSI),
		lbs:      make(map[string]map[string]model.LoadBalancer),
		buckets:  make(map[string]map[string]model.Bucket),
		dbs:      make(map[string]map[string]model.Database),
		clusters: make(map[string]map[string]model.KubernetesCluster),
	}
}

func (s *MemoryStore) CreateVPC(_ context.Context, token string, vpc model.VPC) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.vpcs[token] == nil {
		s.vpcs[token] = make(map[string]model.VPC)
	}
	s.vpcs[token][vpc.ID] = vpc
	return nil
}

func (s *MemoryStore) GetVPC(_ context.Context, token, id string) (model.VPC, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.vpcs[token][id]
	return v, ok, nil
}

func (s *MemoryStore) ListVPCs(_ context.Context, token string) ([]model.VPC, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	m := s.vpcs[token]
	result := make([]model.VPC, 0, len(m))
	for _, v := range m {
		result = append(result, v)
	}
	return result, nil
}

func (s *MemoryStore) DeleteVPC(_ context.Context, token, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.vpcs[token], id)
	return nil
}

func (s *MemoryStore) RenameVPC(_ context.Context, token, id, name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.vpcs[token][id]
	if !ok {
		return fmt.Errorf("VPC %s not found", id)
	}
	v.Name = name
	s.vpcs[token][id] = v
	return nil
}

func (s *MemoryStore) CreateSubnet(_ context.Context, token string, sub model.Subnet) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.subnets[token] == nil {
		s.subnets[token] = make(map[string]model.Subnet)
	}
	s.subnets[token][sub.ID] = sub
	return nil
}

func (s *MemoryStore) GetSubnet(_ context.Context, token, id string) (model.Subnet, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.subnets[token][id]
	return v, ok, nil
}

func (s *MemoryStore) ListSubnets(_ context.Context, token string) ([]model.Subnet, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	m := s.subnets[token]
	result := make([]model.Subnet, 0, len(m))
	for _, v := range m {
		result = append(result, v)
	}
	return result, nil
}

func (s *MemoryStore) DeleteSubnet(_ context.Context, token, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.subnets[token], id)
	return nil
}

func (s *MemoryStore) RenameSubnet(_ context.Context, token, id, name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.subnets[token][id]
	if !ok {
		return fmt.Errorf("Subnet %s not found", id)
	}
	v.Name = name
	s.subnets[token][id] = v
	return nil
}

func (s *MemoryStore) CreateVSI(_ context.Context, token string, v model.VSI) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.vsis[token] == nil {
		s.vsis[token] = make(map[string]model.VSI)
	}
	s.vsis[token][v.ID] = v
	return nil
}

func (s *MemoryStore) GetVSI(_ context.Context, token, id string) (model.VSI, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.vsis[token][id]
	return v, ok, nil
}

func (s *MemoryStore) ListVSIs(_ context.Context, token string) ([]model.VSI, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	m := s.vsis[token]
	result := make([]model.VSI, 0, len(m))
	for _, v := range m {
		result = append(result, v)
	}
	return result, nil
}

func (s *MemoryStore) DeleteVSI(_ context.Context, token, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.vsis[token], id)
	return nil
}

func (s *MemoryStore) RenameVSI(_ context.Context, token, id, name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.vsis[token][id]
	if !ok {
		return fmt.Errorf("VSI %s not found", id)
	}
	v.Name = name
	s.vsis[token][id] = v
	return nil
}

func (s *MemoryStore) UpdateVSIStatus(_ context.Context, token, id, status string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.vsis[token][id]
	if !ok {
		return fmt.Errorf("VSI %s not found", id)
	}
	v.Status = status
	s.vsis[token][id] = v
	return nil
}

func (s *MemoryStore) CreateLoadBalancer(_ context.Context, token string, lb model.LoadBalancer) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.lbs[token] == nil {
		s.lbs[token] = make(map[string]model.LoadBalancer)
	}
	s.lbs[token][lb.ID] = lb
	return nil
}

func (s *MemoryStore) GetLoadBalancer(_ context.Context, token, id string) (model.LoadBalancer, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.lbs[token][id]
	return v, ok, nil
}

func (s *MemoryStore) ListLoadBalancers(_ context.Context, token string) ([]model.LoadBalancer, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	m := s.lbs[token]
	result := make([]model.LoadBalancer, 0, len(m))
	for _, v := range m {
		result = append(result, v)
	}
	return result, nil
}

func (s *MemoryStore) DeleteLoadBalancer(_ context.Context, token, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.lbs[token], id)
	return nil
}

func (s *MemoryStore) RenameLoadBalancer(_ context.Context, token, id, name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.lbs[token][id]
	if !ok {
		return fmt.Errorf("LoadBalancer %s not found", id)
	}
	v.Name = name
	s.lbs[token][id] = v
	return nil
}

func (s *MemoryStore) CreateBucket(_ context.Context, token string, b model.Bucket) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.buckets[token] == nil {
		s.buckets[token] = make(map[string]model.Bucket)
	}
	s.buckets[token][b.ID] = b
	return nil
}

func (s *MemoryStore) GetBucket(_ context.Context, token, id string) (model.Bucket, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.buckets[token][id]
	return v, ok, nil
}

func (s *MemoryStore) ListBuckets(_ context.Context, token string) ([]model.Bucket, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	m := s.buckets[token]
	result := make([]model.Bucket, 0, len(m))
	for _, v := range m {
		result = append(result, v)
	}
	return result, nil
}

func (s *MemoryStore) DeleteBucket(_ context.Context, token, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.buckets[token], id)
	return nil
}

func (s *MemoryStore) RenameBucket(_ context.Context, token, id, name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.buckets[token][id]
	if !ok {
		return fmt.Errorf("Bucket %s not found", id)
	}
	v.Name = name
	s.buckets[token][id] = v
	return nil
}

func (s *MemoryStore) CreateDatabase(_ context.Context, token string, db model.Database) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.dbs[token] == nil {
		s.dbs[token] = make(map[string]model.Database)
	}
	s.dbs[token][db.ID] = db
	return nil
}

func (s *MemoryStore) GetDatabase(_ context.Context, token, id string) (model.Database, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.dbs[token][id]
	return v, ok, nil
}

func (s *MemoryStore) ListDatabases(_ context.Context, token string) ([]model.Database, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	m := s.dbs[token]
	result := make([]model.Database, 0, len(m))
	for _, v := range m {
		result = append(result, v)
	}
	return result, nil
}

func (s *MemoryStore) DeleteDatabase(_ context.Context, token, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.dbs[token], id)
	return nil
}

func (s *MemoryStore) RenameDatabase(_ context.Context, token, id, name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.dbs[token][id]
	if !ok {
		return fmt.Errorf("Database %s not found", id)
	}
	v.Name = name
	s.dbs[token][id] = v
	return nil
}

func (s *MemoryStore) CreateKubernetesCluster(_ context.Context, token string, c model.KubernetesCluster) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.clusters[token] == nil {
		s.clusters[token] = make(map[string]model.KubernetesCluster)
	}
	s.clusters[token][c.ID] = c
	return nil
}

func (s *MemoryStore) GetKubernetesCluster(_ context.Context, token, id string) (model.KubernetesCluster, bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.clusters[token][id]
	return v, ok, nil
}

func (s *MemoryStore) ListKubernetesClusters(_ context.Context, token string) ([]model.KubernetesCluster, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	m := s.clusters[token]
	result := make([]model.KubernetesCluster, 0, len(m))
	for _, v := range m {
		result = append(result, v)
	}
	return result, nil
}

func (s *MemoryStore) DeleteKubernetesCluster(_ context.Context, token, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.clusters[token], id)
	return nil
}

func (s *MemoryStore) RenameKubernetesCluster(_ context.Context, token, id, name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.clusters[token][id]
	if !ok {
		return fmt.Errorf("KubernetesCluster %s not found", id)
	}
	v.Name = name
	s.clusters[token][id] = v
	return nil
}
