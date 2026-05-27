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

func TestSubnet_MissingAuth(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(store.NewMemoryStore()))
	defer srv.Close()

	resp := doRequest(t, "GET", srv.URL+"/v1/subnets", "", "")
	mustStatus(t, resp, 401)
}

func TestSubnet_Lifecycle(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(store.NewMemoryStore()))
	defer srv.Close()

	token := "test-token"

	// create VPC first
	resp := doRequest(t, "POST", srv.URL+"/v1/vpcs", token, `{"name":"my-vpc"}`)
	mustStatus(t, resp, 201)
	var vpc model.VPC
	json.NewDecoder(resp.Body).Decode(&vpc)

	base := srv.URL + "/v1/subnets"

	// empty list
	resp = doRequest(t, "GET", base, token, "")
	mustStatus(t, resp, 200)

	// create
	body := fmt.Sprintf(`{"name":"my-subnet","vpc":{"id":"%s"}}`, vpc.ID)
	resp = doRequest(t, "POST", base, token, body)
	mustStatus(t, resp, 201)
	var sub model.Subnet
	json.NewDecoder(resp.Body).Decode(&sub)
	if sub.ID == "" || sub.Name != "my-subnet" || sub.VPCID != vpc.ID || sub.CIDRBlock == "" || sub.CRN == "" {
		t.Fatalf("unexpected subnet: %+v", sub)
	}

	// get
	resp = doRequest(t, "GET", base+"/"+sub.ID, token, "")
	mustStatus(t, resp, 200)

	// token isolation
	resp = doRequest(t, "GET", base+"/"+sub.ID, "other-token", "")
	mustStatus(t, resp, 404)

	// delete
	resp = doRequest(t, "DELETE", base+"/"+sub.ID, token, "")
	mustStatus(t, resp, 204)

	// gone
	resp = doRequest(t, "GET", base+"/"+sub.ID, token, "")
	mustStatus(t, resp, 404)
}

func TestSubnet_Create_InvalidVPC(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(store.NewMemoryStore()))
	defer srv.Close()

	resp := doRequest(t, "POST", srv.URL+"/v1/subnets", "tok",
		`{"name":"s","vpc":{"id":"nonexistent-vpc"}}`)
	mustStatus(t, resp, 404)
}

func TestSubnet_Create_BadRequest(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(store.NewMemoryStore()))
	defer srv.Close()

	// missing name
	resp := doRequest(t, "POST", srv.URL+"/v1/subnets", "tok", `{"vpc":{"id":"x"}}`)
	mustStatus(t, resp, 400)

	// missing vpc
	resp = doRequest(t, "POST", srv.URL+"/v1/subnets", "tok", `{"name":"s"}`)
	mustStatus(t, resp, 400)
}
