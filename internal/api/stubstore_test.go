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
	createSubnet    func(context.Context, string, model.Subnet) error
	getSubnet       func(context.Context, string, string) (model.Subnet, bool, error)
	listSubnets     func(context.Context, string) ([]model.Subnet, error)
	deleteSubnet    func(context.Context, string, string) error
	createVSI       func(context.Context, string, model.VSI) error
	getVSI          func(context.Context, string, string) (model.VSI, bool, error)
	listVSIs        func(context.Context, string) ([]model.VSI, error)
	deleteVSI       func(context.Context, string, string) error
	updateVSIStatus func(context.Context, string, string, string) error
}

func newErrStore() *funcStore {
	return &funcStore{
		createVPC:       func(context.Context, string, model.VPC) error { return errStub },
		getVPC:          func(context.Context, string, string) (model.VPC, bool, error) { return model.VPC{}, false, errStub },
		listVPCs:        func(context.Context, string) ([]model.VPC, error) { return nil, errStub },
		deleteVPC:       func(context.Context, string, string) error { return errStub },
		createSubnet:    func(context.Context, string, model.Subnet) error { return errStub },
		getSubnet:       func(context.Context, string, string) (model.Subnet, bool, error) { return model.Subnet{}, false, errStub },
		listSubnets:     func(context.Context, string) ([]model.Subnet, error) { return nil, errStub },
		deleteSubnet:    func(context.Context, string, string) error { return errStub },
		createVSI:       func(context.Context, string, model.VSI) error { return errStub },
		getVSI:          func(context.Context, string, string) (model.VSI, bool, error) { return model.VSI{}, false, errStub },
		listVSIs:        func(context.Context, string) ([]model.VSI, error) { return nil, errStub },
		deleteVSI:       func(context.Context, string, string) error { return errStub },
		updateVSIStatus: func(context.Context, string, string, string) error { return errStub },
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
