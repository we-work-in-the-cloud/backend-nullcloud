package api_test

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/we-work-in-the-cloud/nullcloud/backend/internal/api"
	"github.com/we-work-in-the-cloud/nullcloud/backend/internal/model"
)

// helpers that satisfy the funcStore field signatures

func okGetVPC(_ context.Context, _, _ string) (model.VPC, bool, error) {
	return model.VPC{ID: "vpc-1", Name: "v"}, true, nil
}

func okGetSubnet(_ context.Context, _, _ string) (model.Subnet, bool, error) {
	return model.Subnet{ID: "sub-1", Name: "s", VPCID: "vpc-1"}, true, nil
}

func okGetVSI(_ context.Context, _, _ string) (model.VSI, bool, error) {
	return model.VSI{ID: "vsi-1", Name: "v", SubnetID: "sub-1", VPCID: "vpc-1", Status: "running"}, true, nil
}

func notFoundGetVSI(_ context.Context, _, _ string) (model.VSI, bool, error) {
	return model.VSI{}, false, nil
}

func okUpdateVSIStatus(_ context.Context, _, _, _ string) error { return nil }

// tokenMiddleware: token not in allowed list → 401
func TestMiddleware_TokenNotAllowed(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(newErrStore(), []string{"allowed-token"}))
	defer srv.Close()
	mustStatus(t, doRequest(t, "GET", srv.URL+"/v1/vpcs", "wrong-token", ""), 401)
}

// tokenMiddleware: correct token passes through (500 from errStore, not 401)
func TestMiddleware_TokenAllowed(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(newErrStore(), []string{"good-token"}))
	defer srv.Close()
	mustStatus(t, doRequest(t, "GET", srv.URL+"/v1/vpcs", "good-token", ""), 500)
}

// VPC error paths

func TestVPC_List_StoreError(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(newErrStore(), nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "GET", srv.URL+"/v1/vpcs", "tok", ""), 500)
}

func TestVPC_Create_StoreError(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(newErrStore(), nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "POST", srv.URL+"/v1/vpcs", "tok", `{"name":"x"}`), 500)
}

func TestVPC_Get_StoreError(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(newErrStore(), nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "GET", srv.URL+"/v1/vpcs/vpc-1", "tok", ""), 500)
}

func TestVPC_Delete_GetStoreError(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(newErrStore(), nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "DELETE", srv.URL+"/v1/vpcs/vpc-1", "tok", ""), 500)
}

func TestVPC_Delete_DeleteStoreError(t *testing.T) {
	fs := newErrStore()
	fs.getVPC = okGetVPC
	srv := httptest.NewServer(api.NewServer(fs, nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "DELETE", srv.URL+"/v1/vpcs/vpc-1", "tok", ""), 500)
}

// Subnet error paths

func TestSubnet_List_StoreError(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(newErrStore(), nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "GET", srv.URL+"/v1/subnets", "tok", ""), 500)
}

func TestSubnet_Create_GetVPCStoreError(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(newErrStore(), nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "POST", srv.URL+"/v1/subnets", "tok", `{"name":"s","vpc":{"id":"v1"}}`), 500)
}

func TestSubnet_Create_CreateStoreError(t *testing.T) {
	fs := newErrStore()
	fs.getVPC = okGetVPC
	srv := httptest.NewServer(api.NewServer(fs, nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "POST", srv.URL+"/v1/subnets", "tok", `{"name":"s","vpc":{"id":"v1"}}`), 500)
}

func TestSubnet_Get_StoreError(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(newErrStore(), nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "GET", srv.URL+"/v1/subnets/sub-1", "tok", ""), 500)
}

func TestSubnet_Delete_GetStoreError(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(newErrStore(), nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "DELETE", srv.URL+"/v1/subnets/sub-1", "tok", ""), 500)
}

func TestSubnet_Delete_DeleteStoreError(t *testing.T) {
	fs := newErrStore()
	fs.getSubnet = okGetSubnet
	srv := httptest.NewServer(api.NewServer(fs, nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "DELETE", srv.URL+"/v1/subnets/sub-1", "tok", ""), 500)
}

// VSI error paths

func TestVSI_List_StoreError(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(newErrStore(), nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "GET", srv.URL+"/v1/instances", "tok", ""), 500)
}

func TestVSI_Create_GetSubnetStoreError(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(newErrStore(), nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "POST", srv.URL+"/v1/instances", "tok",
		`{"name":"v","subnet":{"id":"s1"}}`), 500)
}

func TestVSI_Create_CreateStoreError(t *testing.T) {
	fs := newErrStore()
	fs.getSubnet = okGetSubnet
	srv := httptest.NewServer(api.NewServer(fs, nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "POST", srv.URL+"/v1/instances", "tok",
		`{"name":"v","subnet":{"id":"s1"}}`), 500)
}

func TestVSI_Get_StoreError(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(newErrStore(), nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "GET", srv.URL+"/v1/instances/vsi-1", "tok", ""), 500)
}

func TestVSI_Delete_GetStoreError(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(newErrStore(), nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "DELETE", srv.URL+"/v1/instances/vsi-1", "tok", ""), 500)
}

func TestVSI_Delete_DeleteStoreError(t *testing.T) {
	fs := newErrStore()
	fs.getVSI = okGetVSI
	srv := httptest.NewServer(api.NewServer(fs, nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "DELETE", srv.URL+"/v1/instances/vsi-1", "tok", ""), 500)
}

// vsiAction paths

func TestVSIAction_BadBody(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(newErrStore(), nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "POST", srv.URL+"/v1/instances/vsi-1/actions", "tok", `{bad`), 400)
}

func TestVSIAction_BadType(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(newErrStore(), nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "POST", srv.URL+"/v1/instances/vsi-1/actions", "tok", `{"type":"explode"}`), 400)
}

func TestVSIAction_NotFound(t *testing.T) {
	fs := newErrStore()
	fs.getVSI = notFoundGetVSI
	srv := httptest.NewServer(api.NewServer(fs, nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "POST", srv.URL+"/v1/instances/vsi-1/actions", "tok", `{"type":"start"}`), 404)
}

func TestVSIAction_GetStoreError(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(newErrStore(), nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "POST", srv.URL+"/v1/instances/vsi-1/actions", "tok", `{"type":"stop"}`), 500)
}

func TestVSIAction_UpdateStoreError(t *testing.T) {
	fs := newErrStore()
	fs.getVSI = okGetVSI
	srv := httptest.NewServer(api.NewServer(fs, nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "POST", srv.URL+"/v1/instances/vsi-1/actions", "tok", `{"type":"restart"}`), 500)
}

func TestVSIAction_Start(t *testing.T) {
	fs := newErrStore()
	fs.getVSI = okGetVSI
	fs.updateVSIStatus = okUpdateVSIStatus
	srv := httptest.NewServer(api.NewServer(fs, nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "POST", srv.URL+"/v1/instances/vsi-1/actions", "tok", `{"type":"start"}`), 200)
}

func TestVSIAction_Stop(t *testing.T) {
	fs := newErrStore()
	fs.getVSI = okGetVSI
	fs.updateVSIStatus = okUpdateVSIStatus
	srv := httptest.NewServer(api.NewServer(fs, nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "POST", srv.URL+"/v1/instances/vsi-1/actions", "tok", `{"type":"stop"}`), 200)
}

func TestVSIAction_Restart(t *testing.T) {
	fs := newErrStore()
	fs.getVSI = okGetVSI
	fs.updateVSIStatus = okUpdateVSIStatus
	srv := httptest.NewServer(api.NewServer(fs, nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "POST", srv.URL+"/v1/instances/vsi-1/actions", "tok", `{"type":"restart"}`), 200)
}
