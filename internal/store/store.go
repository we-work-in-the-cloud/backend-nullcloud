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

	CreateSubnet(ctx context.Context, token string, s model.Subnet) error
	GetSubnet(ctx context.Context, token, id string) (model.Subnet, bool, error)
	ListSubnets(ctx context.Context, token string) ([]model.Subnet, error)
	DeleteSubnet(ctx context.Context, token, id string) error

	CreateVSI(ctx context.Context, token string, v model.VSI) error
	GetVSI(ctx context.Context, token, id string) (model.VSI, bool, error)
	ListVSIs(ctx context.Context, token string) ([]model.VSI, error)
	DeleteVSI(ctx context.Context, token, id string) error
}
