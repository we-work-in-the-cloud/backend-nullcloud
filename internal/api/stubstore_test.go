package api_test

import (
	"context"
	"errors"

	"github.com/we-work-in-the-cloud/nullcloud/backend/internal/model"
)

var errStub = errors.New("store error")

// funcStore is a configurable store stub for testing error paths.
// Set each field to a function; call newErrStore() to get one where every
// method returns errStub by default.
type funcStore struct {
	createVPC       func(context.Context, string, model.VPC) error
	getVPC          func(context.Context, string, string) (model.VPC, bool, error)
	listVPCs        func(context.Context, string) ([]model.VPC, error)
	deleteVPC       func(context.Context, string, string) error
	renameVPC       func(context.Context, string, string, string) error
	createSubnet    func(context.Context, string, model.Subnet) error
	getSubnet       func(context.Context, string, string) (model.Subnet, bool, error)
	listSubnets     func(context.Context, string) ([]model.Subnet, error)
	deleteSubnet    func(context.Context, string, string) error
	renameSubnet    func(context.Context, string, string, string) error
	createVSI       func(context.Context, string, model.VSI) error
	getVSI          func(context.Context, string, string) (model.VSI, bool, error)
	listVSIs        func(context.Context, string) ([]model.VSI, error)
	deleteVSI       func(context.Context, string, string) error
	updateVSIStatus func(context.Context, string, string, string) error
	renameVSI       func(context.Context, string, string, string) error

	createLoadBalancer func(context.Context, string, model.LoadBalancer) error
	getLoadBalancer    func(context.Context, string, string) (model.LoadBalancer, bool, error)
	listLoadBalancers  func(context.Context, string) ([]model.LoadBalancer, error)
	deleteLoadBalancer func(context.Context, string, string) error
	renameLoadBalancer func(context.Context, string, string, string) error

	createBucket func(context.Context, string, model.Bucket) error
	getBucket    func(context.Context, string, string) (model.Bucket, bool, error)
	listBuckets  func(context.Context, string) ([]model.Bucket, error)
	deleteBucket func(context.Context, string, string) error
	renameBucket func(context.Context, string, string, string) error

	createDatabase func(context.Context, string, model.Database) error
	getDatabase    func(context.Context, string, string) (model.Database, bool, error)
	listDatabases  func(context.Context, string) ([]model.Database, error)
	deleteDatabase func(context.Context, string, string) error
	renameDatabase func(context.Context, string, string, string) error

	createKubernetesCluster func(context.Context, string, model.KubernetesCluster) error
	getKubernetesCluster    func(context.Context, string, string) (model.KubernetesCluster, bool, error)
	listKubernetesClusters  func(context.Context, string) ([]model.KubernetesCluster, error)
	deleteKubernetesCluster func(context.Context, string, string) error
	renameKubernetesCluster func(context.Context, string, string, string) error
}

func newErrStore() *funcStore {
	return &funcStore{
		createVPC:    func(context.Context, string, model.VPC) error { return errStub },
		getVPC:       func(context.Context, string, string) (model.VPC, bool, error) { return model.VPC{}, false, errStub },
		listVPCs:     func(context.Context, string) ([]model.VPC, error) { return nil, errStub },
		deleteVPC:    func(context.Context, string, string) error { return errStub },
		renameVPC:    func(context.Context, string, string, string) error { return errStub },
		createSubnet: func(context.Context, string, model.Subnet) error { return errStub },
		getSubnet: func(context.Context, string, string) (model.Subnet, bool, error) {
			return model.Subnet{}, false, errStub
		},
		listSubnets:     func(context.Context, string) ([]model.Subnet, error) { return nil, errStub },
		deleteSubnet:    func(context.Context, string, string) error { return errStub },
		renameSubnet:    func(context.Context, string, string, string) error { return errStub },
		createVSI:       func(context.Context, string, model.VSI) error { return errStub },
		getVSI:          func(context.Context, string, string) (model.VSI, bool, error) { return model.VSI{}, false, errStub },
		listVSIs:        func(context.Context, string) ([]model.VSI, error) { return nil, errStub },
		deleteVSI:       func(context.Context, string, string) error { return errStub },
		updateVSIStatus: func(context.Context, string, string, string) error { return errStub },
		renameVSI:       func(context.Context, string, string, string) error { return errStub },

		createLoadBalancer: func(context.Context, string, model.LoadBalancer) error { return errStub },
		getLoadBalancer: func(context.Context, string, string) (model.LoadBalancer, bool, error) {
			return model.LoadBalancer{}, false, errStub
		},
		listLoadBalancers:  func(context.Context, string) ([]model.LoadBalancer, error) { return nil, errStub },
		deleteLoadBalancer: func(context.Context, string, string) error { return errStub },
		renameLoadBalancer: func(context.Context, string, string, string) error { return errStub },

		createBucket: func(context.Context, string, model.Bucket) error { return errStub },
		getBucket: func(context.Context, string, string) (model.Bucket, bool, error) {
			return model.Bucket{}, false, errStub
		},
		listBuckets:  func(context.Context, string) ([]model.Bucket, error) { return nil, errStub },
		deleteBucket: func(context.Context, string, string) error { return errStub },
		renameBucket: func(context.Context, string, string, string) error { return errStub },

		createDatabase: func(context.Context, string, model.Database) error { return errStub },
		getDatabase: func(context.Context, string, string) (model.Database, bool, error) {
			return model.Database{}, false, errStub
		},
		listDatabases:  func(context.Context, string) ([]model.Database, error) { return nil, errStub },
		deleteDatabase: func(context.Context, string, string) error { return errStub },
		renameDatabase: func(context.Context, string, string, string) error { return errStub },

		createKubernetesCluster: func(context.Context, string, model.KubernetesCluster) error { return errStub },
		getKubernetesCluster: func(context.Context, string, string) (model.KubernetesCluster, bool, error) {
			return model.KubernetesCluster{}, false, errStub
		},
		listKubernetesClusters:  func(context.Context, string) ([]model.KubernetesCluster, error) { return nil, errStub },
		deleteKubernetesCluster: func(context.Context, string, string) error { return errStub },
		renameKubernetesCluster: func(context.Context, string, string, string) error { return errStub },
	}
}

func (s *funcStore) CreateVPC(ctx context.Context, token string, v model.VPC) error {
	return s.createVPC(ctx, token, v)
}
func (s *funcStore) GetVPC(ctx context.Context, token, id string) (model.VPC, bool, error) {
	return s.getVPC(ctx, token, id)
}
func (s *funcStore) ListVPCs(ctx context.Context, token string) ([]model.VPC, error) {
	return s.listVPCs(ctx, token)
}
func (s *funcStore) DeleteVPC(ctx context.Context, token, id string) error {
	return s.deleteVPC(ctx, token, id)
}
func (s *funcStore) RenameVPC(ctx context.Context, token, id, name string) error {
	return s.renameVPC(ctx, token, id, name)
}
func (s *funcStore) CreateSubnet(ctx context.Context, token string, v model.Subnet) error {
	return s.createSubnet(ctx, token, v)
}
func (s *funcStore) GetSubnet(ctx context.Context, token, id string) (model.Subnet, bool, error) {
	return s.getSubnet(ctx, token, id)
}
func (s *funcStore) ListSubnets(ctx context.Context, token string) ([]model.Subnet, error) {
	return s.listSubnets(ctx, token)
}
func (s *funcStore) DeleteSubnet(ctx context.Context, token, id string) error {
	return s.deleteSubnet(ctx, token, id)
}
func (s *funcStore) RenameSubnet(ctx context.Context, token, id, name string) error {
	return s.renameSubnet(ctx, token, id, name)
}
func (s *funcStore) CreateVSI(ctx context.Context, token string, v model.VSI) error {
	return s.createVSI(ctx, token, v)
}
func (s *funcStore) GetVSI(ctx context.Context, token, id string) (model.VSI, bool, error) {
	return s.getVSI(ctx, token, id)
}
func (s *funcStore) ListVSIs(ctx context.Context, token string) ([]model.VSI, error) {
	return s.listVSIs(ctx, token)
}
func (s *funcStore) DeleteVSI(ctx context.Context, token, id string) error {
	return s.deleteVSI(ctx, token, id)
}
func (s *funcStore) UpdateVSIStatus(ctx context.Context, token, id, status string) error {
	return s.updateVSIStatus(ctx, token, id, status)
}
func (s *funcStore) RenameVSI(ctx context.Context, token, id, name string) error {
	return s.renameVSI(ctx, token, id, name)
}

func (s *funcStore) CreateLoadBalancer(ctx context.Context, token string, lb model.LoadBalancer) error {
	return s.createLoadBalancer(ctx, token, lb)
}
func (s *funcStore) GetLoadBalancer(ctx context.Context, token, id string) (model.LoadBalancer, bool, error) {
	return s.getLoadBalancer(ctx, token, id)
}
func (s *funcStore) ListLoadBalancers(ctx context.Context, token string) ([]model.LoadBalancer, error) {
	return s.listLoadBalancers(ctx, token)
}
func (s *funcStore) DeleteLoadBalancer(ctx context.Context, token, id string) error {
	return s.deleteLoadBalancer(ctx, token, id)
}
func (s *funcStore) RenameLoadBalancer(ctx context.Context, token, id, name string) error {
	return s.renameLoadBalancer(ctx, token, id, name)
}

func (s *funcStore) CreateBucket(ctx context.Context, token string, b model.Bucket) error {
	return s.createBucket(ctx, token, b)
}
func (s *funcStore) GetBucket(ctx context.Context, token, id string) (model.Bucket, bool, error) {
	return s.getBucket(ctx, token, id)
}
func (s *funcStore) ListBuckets(ctx context.Context, token string) ([]model.Bucket, error) {
	return s.listBuckets(ctx, token)
}
func (s *funcStore) DeleteBucket(ctx context.Context, token, id string) error {
	return s.deleteBucket(ctx, token, id)
}
func (s *funcStore) RenameBucket(ctx context.Context, token, id, name string) error {
	return s.renameBucket(ctx, token, id, name)
}

func (s *funcStore) CreateDatabase(ctx context.Context, token string, db model.Database) error {
	return s.createDatabase(ctx, token, db)
}
func (s *funcStore) GetDatabase(ctx context.Context, token, id string) (model.Database, bool, error) {
	return s.getDatabase(ctx, token, id)
}
func (s *funcStore) ListDatabases(ctx context.Context, token string) ([]model.Database, error) {
	return s.listDatabases(ctx, token)
}
func (s *funcStore) DeleteDatabase(ctx context.Context, token, id string) error {
	return s.deleteDatabase(ctx, token, id)
}
func (s *funcStore) RenameDatabase(ctx context.Context, token, id, name string) error {
	return s.renameDatabase(ctx, token, id, name)
}

func (s *funcStore) CreateKubernetesCluster(ctx context.Context, token string, c model.KubernetesCluster) error {
	return s.createKubernetesCluster(ctx, token, c)
}
func (s *funcStore) GetKubernetesCluster(ctx context.Context, token, id string) (model.KubernetesCluster, bool, error) {
	return s.getKubernetesCluster(ctx, token, id)
}
func (s *funcStore) ListKubernetesClusters(ctx context.Context, token string) ([]model.KubernetesCluster, error) {
	return s.listKubernetesClusters(ctx, token)
}
func (s *funcStore) DeleteKubernetesCluster(ctx context.Context, token, id string) error {
	return s.deleteKubernetesCluster(ctx, token, id)
}
func (s *funcStore) RenameKubernetesCluster(ctx context.Context, token, id, name string) error {
	return s.renameKubernetesCluster(ctx, token, id, name)
}
