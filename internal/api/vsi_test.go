package api_test

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/we-work-in-the-cloud/nullcloud/backend/internal/api"
	"github.com/we-work-in-the-cloud/nullcloud/backend/internal/model"
	"github.com/we-work-in-the-cloud/nullcloud/backend/internal/store"
)

func TestVSI_MissingAuth(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(store.NewMemoryStore()))
	defer srv.Close()

	resp := doRequest(t, "GET", srv.URL+"/v1/instances", "", "")
	mustStatus(t, resp, 401)
}

func TestVSI_Lifecycle(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(store.NewMemoryStore()))
	defer srv.Close()

	token := "test-token"

	// create VPC
	resp := doRequest(t, "POST", srv.URL+"/v1/vpcs", token, `{"name":"my-vpc"}`)
	mustStatus(t, resp, 201)
	var vpc model.VPC
	json.NewDecoder(resp.Body).Decode(&vpc)

	// create Subnet
	resp = doRequest(t, "POST", srv.URL+"/v1/subnets", token,
		fmt.Sprintf(`{"name":"my-subnet","vpc":{"id":"%s"}}`, vpc.ID))
	mustStatus(t, resp, 201)
	var sub model.Subnet
	json.NewDecoder(resp.Body).Decode(&sub)

	base := srv.URL + "/v1/instances"

	// empty list
	resp = doRequest(t, "GET", base, token, "")
	mustStatus(t, resp, 200)

	// create
	body := fmt.Sprintf(`{"name":"my-vsi","subnet":{"id":"%s"},"profile":{"name":"cx2-2x4"},"image":{"id":"ubuntu-22-04"}}`, sub.ID)
	resp = doRequest(t, "POST", base, token, body)
	mustStatus(t, resp, 201)
	var vsi model.VSI
	json.NewDecoder(resp.Body).Decode(&vsi)
	if vsi.ID == "" || vsi.Name != "my-vsi" || vsi.SubnetID != sub.ID || vsi.VPCID != vpc.ID {
		t.Fatalf("unexpected vsi: %+v", vsi)
	}
	if vsi.Status != "running" || vsi.PrimaryIP == "" || vsi.CRN == "" {
		t.Fatalf("unexpected vsi fields: %+v", vsi)
	}

	// get
	resp = doRequest(t, "GET", base+"/"+vsi.ID, token, "")
	mustStatus(t, resp, 200)

	// token isolation
	resp = doRequest(t, "GET", base+"/"+vsi.ID, "other-token", "")
	mustStatus(t, resp, 404)

	// delete
	resp = doRequest(t, "DELETE", base+"/"+vsi.ID, token, "")
	mustStatus(t, resp, 204)

	// gone
	resp = doRequest(t, "GET", base+"/"+vsi.ID, token, "")
	mustStatus(t, resp, 404)
}

func TestVSI_Create_InvalidSubnet(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(store.NewMemoryStore()))
	defer srv.Close()

	resp := doRequest(t, "POST", srv.URL+"/v1/instances", "tok",
		`{"name":"v","subnet":{"id":"nonexistent"}}`)
	mustStatus(t, resp, 404)
}

func TestVSI_Create_BadRequest(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(store.NewMemoryStore()))
	defer srv.Close()

	// missing name
	resp := doRequest(t, "POST", srv.URL+"/v1/instances", "tok", `{"subnet":{"id":"x"}}`)
	mustStatus(t, resp, 400)

	// missing subnet
	resp = doRequest(t, "POST", srv.URL+"/v1/instances", "tok", `{"name":"v"}`)
	mustStatus(t, resp, 400)
}
