package store

import (
	"context"

	"github.com/we-work-in-the-cloud/nullcloud/backend/internal/model"
)

// Store is the single interface the API layer depends on.
// All methods are token-scoped: two different tokens are fully isolated.
type Store interface {
	CreateVPC(ctx context.Context, token string, vpc model.VPC) error
	GetVPC(ctx context.Context, token, id string) (model.VPC, bool, error)
	ListVPCs(ctx context.Context, token string) ([]model.VPC, error)
	DeleteVPC(ctx context.Context, token, id string) error
	RenameVPC(ctx context.Context, token, id, name string) error

	CreateSubnet(ctx context.Context, token string, s model.Subnet) error
	GetSubnet(ctx context.Context, token, id string) (model.Subnet, bool, error)
	ListSubnets(ctx context.Context, token string) ([]model.Subnet, error)
	DeleteSubnet(ctx context.Context, token, id string) error
	RenameSubnet(ctx context.Context, token, id, name string) error

	CreateVSI(ctx context.Context, token string, v model.VSI) error
	GetVSI(ctx context.Context, token, id string) (model.VSI, bool, error)
	ListVSIs(ctx context.Context, token string) ([]model.VSI, error)
	DeleteVSI(ctx context.Context, token, id string) error
	UpdateVSIStatus(ctx context.Context, token, id, status string) error
	RenameVSI(ctx context.Context, token, id, name string) error

	CreateLoadBalancer(ctx context.Context, token string, lb model.LoadBalancer) error
	GetLoadBalancer(ctx context.Context, token, id string) (model.LoadBalancer, bool, error)
	ListLoadBalancers(ctx context.Context, token string) ([]model.LoadBalancer, error)
	DeleteLoadBalancer(ctx context.Context, token, id string) error
	RenameLoadBalancer(ctx context.Context, token, id, name string) error

	CreateBucket(ctx context.Context, token string, b model.Bucket) error
	GetBucket(ctx context.Context, token, id string) (model.Bucket, bool, error)
	ListBuckets(ctx context.Context, token string) ([]model.Bucket, error)
	DeleteBucket(ctx context.Context, token, id string) error
	RenameBucket(ctx context.Context, token, id, name string) error

	CreateDatabase(ctx context.Context, token string, db model.Database) error
	GetDatabase(ctx context.Context, token, id string) (model.Database, bool, error)
	ListDatabases(ctx context.Context, token string) ([]model.Database, error)
	DeleteDatabase(ctx context.Context, token, id string) error
	RenameDatabase(ctx context.Context, token, id, name string) error

	CreateKubernetesCluster(ctx context.Context, token string, c model.KubernetesCluster) error
	GetKubernetesCluster(ctx context.Context, token, id string) (model.KubernetesCluster, bool, error)
	ListKubernetesClusters(ctx context.Context, token string) ([]model.KubernetesCluster, error)
	DeleteKubernetesCluster(ctx context.Context, token, id string) error
	RenameKubernetesCluster(ctx context.Context, token, id, name string) error
}
