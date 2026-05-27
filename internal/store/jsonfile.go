package store

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/we-work-in-the-cloud/nullcloud/backend/internal/model"
)

type fileData struct {
	VPCs    map[string]map[string]model.VPC    `json:"vpcs"`
	Subnets map[string]map[string]model.Subnet `json:"subnets"`
	VSIs    map[string]map[string]model.VSI    `json:"vsis"`
}

type JSONFileStore struct {
	mu   sync.Mutex
	path string
	data fileData
}

func NewJSONFileStore(path string) (*JSONFileStore, error) {
	s := &JSONFileStore{
		path: path,
		data: fileData{
			VPCs:    make(map[string]map[string]model.VPC),
			Subnets: make(map[string]map[string]model.Subnet),
			VSIs:    make(map[string]map[string]model.VSI),
		},
	}
	b, err := os.ReadFile(path)
	if err == nil {
		if err := json.Unmarshal(b, &s.data); err != nil {
			return nil, fmt.Errorf("parsing store file: %w", err)
		}
		if s.data.VPCs == nil {
			s.data.VPCs = make(map[string]map[string]model.VPC)
		}
		if s.data.Subnets == nil {
			s.data.Subnets = make(map[string]map[string]model.Subnet)
		}
		if s.data.VSIs == nil {
			s.data.VSIs = make(map[string]map[string]model.VSI)
		}
	} else if !os.IsNotExist(err) {
		return nil, fmt.Errorf("opening store file: %w", err)
	}
	return s, nil
}

func (s *JSONFileStore) flush() error {
	b, err := json.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return err
	}
	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, b, 0644); err != nil {
		return err
	}
	return os.Rename(tmp, s.path)
}

func (s *JSONFileStore) CreateVPC(_ context.Context, token string, vpc model.VPC) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.data.VPCs[token] == nil {
		s.data.VPCs[token] = make(map[string]model.VPC)
	}
	s.data.VPCs[token][vpc.ID] = vpc
	return s.flush()
}

func (s *JSONFileStore) GetVPC(_ context.Context, token, id string) (model.VPC, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.data.VPCs[token][id]
	return v, ok, nil
}

func (s *JSONFileStore) ListVPCs(_ context.Context, token string) ([]model.VPC, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	m := s.data.VPCs[token]
	result := make([]model.VPC, 0, len(m))
	for _, v := range m {
		result = append(result, v)
	}
	return result, nil
}

func (s *JSONFileStore) DeleteVPC(_ context.Context, token, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data.VPCs[token], id)
	return s.flush()
}

func (s *JSONFileStore) CreateSubnet(_ context.Context, token string, sub model.Subnet) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.data.Subnets[token] == nil {
		s.data.Subnets[token] = make(map[string]model.Subnet)
	}
	s.data.Subnets[token][sub.ID] = sub
	return s.flush()
}

func (s *JSONFileStore) GetSubnet(_ context.Context, token, id string) (model.Subnet, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.data.Subnets[token][id]
	return v, ok, nil
}

func (s *JSONFileStore) ListSubnets(_ context.Context, token string) ([]model.Subnet, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	m := s.data.Subnets[token]
	result := make([]model.Subnet, 0, len(m))
	for _, v := range m {
		result = append(result, v)
	}
	return result, nil
}

func (s *JSONFileStore) DeleteSubnet(_ context.Context, token, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data.Subnets[token], id)
	return s.flush()
}

func (s *JSONFileStore) CreateVSI(_ context.Context, token string, v model.VSI) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.data.VSIs[token] == nil {
		s.data.VSIs[token] = make(map[string]model.VSI)
	}
	s.data.VSIs[token][v.ID] = v
	return s.flush()
}

func (s *JSONFileStore) GetVSI(_ context.Context, token, id string) (model.VSI, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.data.VSIs[token][id]
	return v, ok, nil
}

func (s *JSONFileStore) ListVSIs(_ context.Context, token string) ([]model.VSI, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	m := s.data.VSIs[token]
	result := make([]model.VSI, 0, len(m))
	for _, v := range m {
		result = append(result, v)
	}
	return result, nil
}

func (s *JSONFileStore) DeleteVSI(_ context.Context, token, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data.VSIs[token], id)
	return s.flush()
}

func (s *JSONFileStore) UpdateVSIStatus(_ context.Context, token, id, status string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.data.VSIs[token][id]
	if !ok {
		return fmt.Errorf("VSI %s not found", id)
	}
	v.Status = status
	s.data.VSIs[token][id] = v
	return s.flush()
}
