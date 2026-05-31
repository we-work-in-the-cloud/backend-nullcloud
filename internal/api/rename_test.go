package api_test

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/we-work-in-the-cloud/nullcloud/backend/internal/api"
	"github.com/we-work-in-the-cloud/nullcloud/backend/internal/model"
	"github.com/we-work-in-the-cloud/nullcloud/backend/internal/store"
)

func TestVPC_Rename(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(store.NewMemoryStore(), nil))
	defer srv.Close()

	token := "tok"
	resp := doRequest(t, "POST", srv.URL+"/v1/vpcs", token, `{"name":"vpc"}`)
	mustStatus(t, resp, 201)
	var vpc model.VPC
	json.NewDecoder(resp.Body).Decode(&vpc)

	resp = doRequest(t, "PATCH", srv.URL+"/v1/vpcs/"+vpc.ID, token, `{"name":"renamed"}`)
	mustStatus(t, resp, 200)
	var updated model.VPC
	json.NewDecoder(resp.Body).Decode(&updated)
	if updated.Name != "renamed" || updated.ID != vpc.ID {
		t.Fatalf("unexpected: %+v", updated)
	}
}

func TestVPC_Rename_NotFound(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(store.NewMemoryStore(), nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "PATCH", srv.URL+"/v1/vpcs/nonexistent", "tok", `{"name":"x"}`), 404)
}

func TestVPC_Rename_BadRequest(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(store.NewMemoryStore(), nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "PATCH", srv.URL+"/v1/vpcs/x", "tok", `{}`), 400)
}

func TestVPC_Rename_GetStoreError(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(newErrStore(), nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "PATCH", srv.URL+"/v1/vpcs/vpc-1", "tok", `{"name":"x"}`), 500)
}

func TestVPC_Rename_RenameStoreError(t *testing.T) {
	fs := newErrStore()
	fs.getVPC = okGetVPC
	srv := httptest.NewServer(api.NewServer(fs, nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "PATCH", srv.URL+"/v1/vpcs/vpc-1", "tok", `{"name":"x"}`), 500)
}

func TestSubnet_Rename(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(store.NewMemoryStore(), nil))
	defer srv.Close()

	token := "tok"
	resp := doRequest(t, "POST", srv.URL+"/v1/vpcs", token, `{"name":"vpc"}`)
	mustStatus(t, resp, 201)
	var vpc model.VPC
	json.NewDecoder(resp.Body).Decode(&vpc)

	resp = doRequest(t, "POST", srv.URL+"/v1/subnets", token, `{"name":"sub","vpc":{"id":"`+vpc.ID+`"}}`)
	mustStatus(t, resp, 201)
	var sub model.Subnet
	json.NewDecoder(resp.Body).Decode(&sub)

	resp = doRequest(t, "PATCH", srv.URL+"/v1/subnets/"+sub.ID, token, `{"name":"renamed-sub"}`)
	mustStatus(t, resp, 200)
	var updated model.Subnet
	json.NewDecoder(resp.Body).Decode(&updated)
	if updated.Name != "renamed-sub" {
		t.Fatalf("unexpected: %+v", updated)
	}
}

func TestSubnet_Rename_NotFound(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(store.NewMemoryStore(), nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "PATCH", srv.URL+"/v1/subnets/nonexistent", "tok", `{"name":"x"}`), 404)
}

func TestSubnet_Rename_BadRequest(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(store.NewMemoryStore(), nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "PATCH", srv.URL+"/v1/subnets/x", "tok", `{}`), 400)
}

func TestSubnet_Rename_GetStoreError(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(newErrStore(), nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "PATCH", srv.URL+"/v1/subnets/sub-1", "tok", `{"name":"x"}`), 500)
}

func TestSubnet_Rename_RenameStoreError(t *testing.T) {
	fs := newErrStore()
	fs.getSubnet = okGetSubnet
	srv := httptest.NewServer(api.NewServer(fs, nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "PATCH", srv.URL+"/v1/subnets/sub-1", "tok", `{"name":"x"}`), 500)
}

func TestVSI_Rename(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(store.NewMemoryStore(), nil))
	defer srv.Close()

	token := "tok"
	resp := doRequest(t, "POST", srv.URL+"/v1/vpcs", token, `{"name":"vpc"}`)
	mustStatus(t, resp, 201)
	var vpc model.VPC
	json.NewDecoder(resp.Body).Decode(&vpc)

	resp = doRequest(t, "POST", srv.URL+"/v1/subnets", token, `{"name":"sub","vpc":{"id":"`+vpc.ID+`"}}`)
	mustStatus(t, resp, 201)
	var sub model.Subnet
	json.NewDecoder(resp.Body).Decode(&sub)

	resp = doRequest(t, "POST", srv.URL+"/v1/instances", token,
		`{"name":"vsi","subnet":{"id":"`+sub.ID+`"}}`)
	mustStatus(t, resp, 201)
	var vsi model.VSI
	json.NewDecoder(resp.Body).Decode(&vsi)

	resp = doRequest(t, "PATCH", srv.URL+"/v1/instances/"+vsi.ID, token, `{"name":"renamed-vsi"}`)
	mustStatus(t, resp, 200)
	var updated model.VSI
	json.NewDecoder(resp.Body).Decode(&updated)
	if updated.Name != "renamed-vsi" {
		t.Fatalf("unexpected: %+v", updated)
	}
}

func TestVSI_Rename_NotFound(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(store.NewMemoryStore(), nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "PATCH", srv.URL+"/v1/instances/nonexistent", "tok", `{"name":"x"}`), 404)
}

func TestVSI_Rename_BadRequest(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(store.NewMemoryStore(), nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "PATCH", srv.URL+"/v1/instances/x", "tok", `{}`), 400)
}

func TestVSI_Rename_GetStoreError(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(newErrStore(), nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "PATCH", srv.URL+"/v1/instances/vsi-1", "tok", `{"name":"x"}`), 500)
}

func TestVSI_Rename_RenameStoreError(t *testing.T) {
	fs := newErrStore()
	fs.getVSI = okGetVSI
	srv := httptest.NewServer(api.NewServer(fs, nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "PATCH", srv.URL+"/v1/instances/vsi-1", "tok", `{"name":"x"}`), 500)
}
