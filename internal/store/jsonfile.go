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
	VPCs     map[string]map[string]model.VPC              `json:"vpcs"`
	Subnets  map[string]map[string]model.Subnet           `json:"subnets"`
	VSIs     map[string]map[string]model.VSI              `json:"vsis"`
	LBs      map[string]map[string]model.LoadBalancer     `json:"load_balancers"`
	Buckets  map[string]map[string]model.Bucket           `json:"buckets"`
	DBs      map[string]map[string]model.Database         `json:"databases"`
	Clusters map[string]map[string]model.KubernetesCluster `json:"kubernetes_clusters"`
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
			VPCs:     make(map[string]map[string]model.VPC),
			Subnets:  make(map[string]map[string]model.Subnet),
			VSIs:     make(map[string]map[string]model.VSI),
			LBs:      make(map[string]map[string]model.LoadBalancer),
			Buckets:  make(map[string]map[string]model.Bucket),
			DBs:      make(map[string]map[string]model.Database),
			Clusters: make(map[string]map[string]model.KubernetesCluster),
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
		if s.data.LBs == nil {
			s.data.LBs = make(map[string]map[string]model.LoadBalancer)
		}
		if s.data.Buckets == nil {
			s.data.Buckets = make(map[string]map[string]model.Bucket)
		}
		if s.data.DBs == nil {
			s.data.DBs = make(map[string]map[string]model.Database)
		}
		if s.data.Clusters == nil {
			s.data.Clusters = make(map[string]map[string]model.KubernetesCluster)
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

func (s *JSONFileStore) RenameVPC(_ context.Context, token, id, name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.data.VPCs[token][id]
	if !ok {
		return fmt.Errorf("VPC %s not found", id)
	}
	v.Name = name
	s.data.VPCs[token][id] = v
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

func (s *JSONFileStore) RenameSubnet(_ context.Context, token, id, name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.data.Subnets[token][id]
	if !ok {
		return fmt.Errorf("Subnet %s not found", id)
	}
	v.Name = name
	s.data.Subnets[token][id] = v
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

func (s *JSONFileStore) RenameVSI(_ context.Context, token, id, name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.data.VSIs[token][id]
	if !ok {
		return fmt.Errorf("VSI %s not found", id)
	}
	v.Name = name
	s.data.VSIs[token][id] = v
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

func (s *JSONFileStore) CreateLoadBalancer(_ context.Context, token string, lb model.LoadBalancer) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.data.LBs[token] == nil {
		s.data.LBs[token] = make(map[string]model.LoadBalancer)
	}
	s.data.LBs[token][lb.ID] = lb
	return s.flush()
}

func (s *JSONFileStore) GetLoadBalancer(_ context.Context, token, id string) (model.LoadBalancer, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.data.LBs[token][id]
	return v, ok, nil
}

func (s *JSONFileStore) ListLoadBalancers(_ context.Context, token string) ([]model.LoadBalancer, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	m := s.data.LBs[token]
	result := make([]model.LoadBalancer, 0, len(m))
	for _, v := range m {
		result = append(result, v)
	}
	return result, nil
}

func (s *JSONFileStore) DeleteLoadBalancer(_ context.Context, token, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data.LBs[token], id)
	return s.flush()
}

func (s *JSONFileStore) RenameLoadBalancer(_ context.Context, token, id, name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.data.LBs[token][id]
	if !ok {
		return fmt.Errorf("LoadBalancer %s not found", id)
	}
	v.Name = name
	s.data.LBs[token][id] = v
	return s.flush()
}

func (s *JSONFileStore) CreateBucket(_ context.Context, token string, b model.Bucket) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.data.Buckets[token] == nil {
		s.data.Buckets[token] = make(map[string]model.Bucket)
	}
	s.data.Buckets[token][b.ID] = b
	return s.flush()
}

func (s *JSONFileStore) GetBucket(_ context.Context, token, id string) (model.Bucket, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.data.Buckets[token][id]
	return v, ok, nil
}

func (s *JSONFileStore) ListBuckets(_ context.Context, token string) ([]model.Bucket, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	m := s.data.Buckets[token]
	result := make([]model.Bucket, 0, len(m))
	for _, v := range m {
		result = append(result, v)
	}
	return result, nil
}

func (s *JSONFileStore) DeleteBucket(_ context.Context, token, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data.Buckets[token], id)
	return s.flush()
}

func (s *JSONFileStore) RenameBucket(_ context.Context, token, id, name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.data.Buckets[token][id]
	if !ok {
		return fmt.Errorf("Bucket %s not found", id)
	}
	v.Name = name
	s.data.Buckets[token][id] = v
	return s.flush()
}

func (s *JSONFileStore) CreateDatabase(_ context.Context, token string, db model.Database) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.data.DBs[token] == nil {
		s.data.DBs[token] = make(map[string]model.Database)
	}
	s.data.DBs[token][db.ID] = db
	return s.flush()
}

func (s *JSONFileStore) GetDatabase(_ context.Context, token, id string) (model.Database, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.data.DBs[token][id]
	return v, ok, nil
}

func (s *JSONFileStore) ListDatabases(_ context.Context, token string) ([]model.Database, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	m := s.data.DBs[token]
	result := make([]model.Database, 0, len(m))
	for _, v := range m {
		result = append(result, v)
	}
	return result, nil
}

func (s *JSONFileStore) DeleteDatabase(_ context.Context, token, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data.DBs[token], id)
	return s.flush()
}

func (s *JSONFileStore) RenameDatabase(_ context.Context, token, id, name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.data.DBs[token][id]
	if !ok {
		return fmt.Errorf("Database %s not found", id)
	}
	v.Name = name
	s.data.DBs[token][id] = v
	return s.flush()
}

func (s *JSONFileStore) CreateKubernetesCluster(_ context.Context, token string, c model.KubernetesCluster) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.data.Clusters[token] == nil {
		s.data.Clusters[token] = make(map[string]model.KubernetesCluster)
	}
	s.data.Clusters[token][c.ID] = c
	return s.flush()
}

func (s *JSONFileStore) GetKubernetesCluster(_ context.Context, token, id string) (model.KubernetesCluster, bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.data.Clusters[token][id]
	return v, ok, nil
}

func (s *JSONFileStore) ListKubernetesClusters(_ context.Context, token string) ([]model.KubernetesCluster, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	m := s.data.Clusters[token]
	result := make([]model.KubernetesCluster, 0, len(m))
	for _, v := range m {
		result = append(result, v)
	}
	return result, nil
}

func (s *JSONFileStore) DeleteKubernetesCluster(_ context.Context, token, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data.Clusters[token], id)
	return s.flush()
}

func (s *JSONFileStore) RenameKubernetesCluster(_ context.Context, token, id, name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.data.Clusters[token][id]
	if !ok {
		return fmt.Errorf("KubernetesCluster %s not found", id)
	}
	v.Name = name
	s.data.Clusters[token][id] = v
	return s.flush()
}
