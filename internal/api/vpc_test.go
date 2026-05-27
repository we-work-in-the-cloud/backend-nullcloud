package api_test

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/we-work-in-the-cloud/nullcloud/backend/internal/api"
	"github.com/we-work-in-the-cloud/nullcloud/backend/internal/model"
	"github.com/we-work-in-the-cloud/nullcloud/backend/internal/store"
)

func TestVPC_MissingAuth(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(store.NewMemoryStore(), nil))
	defer srv.Close()

	resp := doRequest(t, "GET", srv.URL+"/v1/vpcs", "", "")
	mustStatus(t, resp, 401)
}

func TestVPC_Lifecycle(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(store.NewMemoryStore(), nil))
	defer srv.Close()

	token := "test-token"
	base := srv.URL + "/v1/vpcs"

	// empty list
	resp := doRequest(t, "GET", base, token, "")
	mustStatus(t, resp, 200)
	var listResp struct {
		VPCs []model.VPC `json:"vpcs"`
	}
	json.NewDecoder(resp.Body).Decode(&listResp)
	if len(listResp.VPCs) != 0 {
		t.Fatal("expected empty list")
	}

	// create
	resp = doRequest(t, "POST", base, token, `{"name":"my-vpc"}`)
	mustStatus(t, resp, 201)
	var vpc model.VPC
	json.NewDecoder(resp.Body).Decode(&vpc)
	if vpc.ID == "" || vpc.Name != "my-vpc" || vpc.Status != "available" || vpc.CRN == "" {
		t.Fatalf("unexpected vpc: %+v", vpc)
	}

	// get
	resp = doRequest(t, "GET", base+"/"+vpc.ID, token, "")
	mustStatus(t, resp, 200)

	// list has 1
	resp = doRequest(t, "GET", base, token, "")
	json.NewDecoder(resp.Body).Decode(&listResp)
	if len(listResp.VPCs) != 1 {
		t.Fatalf("expected 1 VPC, got %d", len(listResp.VPCs))
	}

	// token isolation
	resp = doRequest(t, "GET", base+"/"+vpc.ID, "other-token", "")
	mustStatus(t, resp, 404)

	// delete
	resp = doRequest(t, "DELETE", base+"/"+vpc.ID, token, "")
	mustStatus(t, resp, 204)

	// gone
	resp = doRequest(t, "GET", base+"/"+vpc.ID, token, "")
	mustStatus(t, resp, 404)
}

func TestVPC_Create_BadRequest(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(store.NewMemoryStore(), nil))
	defer srv.Close()

	token := "test-token"
	base := srv.URL + "/v1/vpcs"

	// missing name
	resp := doRequest(t, "POST", base, token, `{}`)
	mustStatus(t, resp, 400)

	// malformed JSON
	resp = doRequest(t, "POST", base, token, `{bad json`)
	mustStatus(t, resp, 400)
}

func TestVPC_Delete_NotFound(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(store.NewMemoryStore(), nil))
	defer srv.Close()

	resp := doRequest(t, "DELETE", srv.URL+"/v1/vpcs/nonexistent", "tok", "")
	mustStatus(t, resp, 404)
}
