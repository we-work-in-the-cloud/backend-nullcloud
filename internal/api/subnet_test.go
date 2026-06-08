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
	srv := httptest.NewServer(api.NewServer(store.NewMemoryStore(), nil))
	defer srv.Close()

	resp := doRequest(t, "GET", srv.URL+"/v1/subnets", "", "")
	mustStatus(t, resp, 401)
}

func TestSubnet_Lifecycle(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(store.NewMemoryStore(), nil))
	defer srv.Close()

	token := "test-token"

	// create VPC first (defaults to us-east)
	resp := doRequest(t, "POST", srv.URL+"/v1/vpcs", token, `{"name":"my-vpc"}`)
	mustStatus(t, resp, 201)
	var vpc model.VPC
	json.NewDecoder(resp.Body).Decode(&vpc)

	base := srv.URL + "/v1/subnets"

	// empty list
	resp = doRequest(t, "GET", base, token, "")
	mustStatus(t, resp, 200)

	// create with zone in VPC's region
	body := fmt.Sprintf(`{"name":"my-subnet","vpc":{"id":"%s"},"zone":"us-east-1","cidr_block":"10.0.0.0/24"}`, vpc.ID)
	resp = doRequest(t, "POST", base, token, body)
	mustStatus(t, resp, 201)
	var sub model.Subnet
	json.NewDecoder(resp.Body).Decode(&sub)
	if sub.ID == "" || sub.Name != "my-subnet" || sub.VPCID != vpc.ID || sub.CIDRBlock != "10.0.0.0/24" || sub.CRN == "" {
		t.Fatalf("unexpected subnet: %+v", sub)
	}
	if sub.Zone != "us-east-1" {
		t.Fatalf("expected zone us-east-1, got %q", sub.Zone)
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

func TestSubnet_Zone_AllZonesInRegion(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(store.NewMemoryStore(), nil))
	defer srv.Close()

	token := "tok"

	resp := doRequest(t, "POST", srv.URL+"/v1/vpcs", token, `{"name":"vpc","region":"eu-west"}`)
	mustStatus(t, resp, 201)
	var vpc model.VPC
	json.NewDecoder(resp.Body).Decode(&vpc)

	for i, zone := range []string{"eu-west-1", "eu-west-2", "eu-west-3"} {
		body := fmt.Sprintf(`{"name":"sub-%s","vpc":{"id":"%s"},"zone":"%s","cidr_block":"10.%d.0.0/24"}`, zone, vpc.ID, zone, i)
		resp = doRequest(t, "POST", srv.URL+"/v1/subnets", token, body)
		mustStatus(t, resp, 201)
		var sub model.Subnet
		json.NewDecoder(resp.Body).Decode(&sub)
		if sub.Zone != zone {
			t.Errorf("expected zone %q, got %q", zone, sub.Zone)
		}
	}
}

func TestSubnet_Zone_WrongRegion(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(store.NewMemoryStore(), nil))
	defer srv.Close()

	token := "tok"

	// VPC in us-west; subnet zone from a different region → bad request
	resp := doRequest(t, "POST", srv.URL+"/v1/vpcs", token, `{"name":"vpc","region":"us-west"}`)
	mustStatus(t, resp, 201)
	var vpc model.VPC
	json.NewDecoder(resp.Body).Decode(&vpc)

	for _, zone := range []string{"us-east-1", "eu-west-2", "us-central-3", "us-west"} {
		body := fmt.Sprintf(`{"name":"sub","vpc":{"id":"%s"},"zone":"%s","cidr_block":"10.0.0.0/24"}`, vpc.ID, zone)
		resp = doRequest(t, "POST", srv.URL+"/v1/subnets", token, body)
		mustStatus(t, resp, 400)
	}
}

func TestSubnet_Create_MissingZone(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(store.NewMemoryStore(), nil))
	defer srv.Close()

	token := "tok"

	resp := doRequest(t, "POST", srv.URL+"/v1/vpcs", token, `{"name":"vpc"}`)
	mustStatus(t, resp, 201)
	var vpc model.VPC
	json.NewDecoder(resp.Body).Decode(&vpc)

	body := fmt.Sprintf(`{"name":"sub","vpc":{"id":"%s"},"cidr_block":"10.0.0.0/24"}`, vpc.ID)
	resp = doRequest(t, "POST", srv.URL+"/v1/subnets", token, body)
	mustStatus(t, resp, 400)
}

func TestSubnet_Create_InvalidVPC(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(store.NewMemoryStore(), nil))
	defer srv.Close()

	resp := doRequest(t, "POST", srv.URL+"/v1/subnets", "tok",
		`{"name":"s","vpc":{"id":"nonexistent-vpc"},"zone":"us-east-1","cidr_block":"10.0.0.0/24"}`)
	mustStatus(t, resp, 404)
}

func TestSubnet_Create_BadRequest(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(store.NewMemoryStore(), nil))
	defer srv.Close()

	// missing name
	resp := doRequest(t, "POST", srv.URL+"/v1/subnets", "tok", `{"vpc":{"id":"x"},"zone":"us-east-1","cidr_block":"10.0.0.0/24"}`)
	mustStatus(t, resp, 400)

	// missing vpc
	resp = doRequest(t, "POST", srv.URL+"/v1/subnets", "tok", `{"name":"s","zone":"us-east-1","cidr_block":"10.0.0.0/24"}`)
	mustStatus(t, resp, 400)

	// missing zone
	resp = doRequest(t, "POST", srv.URL+"/v1/subnets", "tok", `{"name":"s","vpc":{"id":"x"},"cidr_block":"10.0.0.0/24"}`)
	mustStatus(t, resp, 400)

	// missing cidr_block
	resp = doRequest(t, "POST", srv.URL+"/v1/subnets", "tok", `{"name":"s","vpc":{"id":"x"},"zone":"us-east-1"}`)
	mustStatus(t, resp, 400)
}
