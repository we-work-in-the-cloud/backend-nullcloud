package store

import (
	"context"
	"sync"

	"github.com/we-work-in-the-cloud/nullcloud/backend/internal/model"
)

type MemoryStore struct {
	mu      sync.RWMutex
	vpcs    map[string]map[string]model.VPC
	subnets map[string]map[string]model.Subnet
	vsis    map[string]map[string]model.VSI
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		vpcs:    make(map[string]map[string]model.VPC),
		subnets: make(map[string]map[string]model.Subnet),
		vsis:    make(map[string]map[string]model.VSI),
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
