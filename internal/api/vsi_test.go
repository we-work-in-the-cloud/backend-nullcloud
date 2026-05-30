package api_test

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http/httptest"
	"testing"

	"github.com/we-work-in-the-cloud/nullcloud/backend/internal/api"
	"github.com/we-work-in-the-cloud/nullcloud/backend/internal/model"
	"github.com/we-work-in-the-cloud/nullcloud/backend/internal/store"
)

func TestVSI_MissingAuth(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(store.NewMemoryStore(), nil))
	defer srv.Close()

	resp := doRequest(t, "GET", srv.URL+"/v1/instances", "", "")
	mustStatus(t, resp, 401)
}

func TestVSI_Lifecycle(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(store.NewMemoryStore(), nil))
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
	srv := httptest.NewServer(api.NewServer(store.NewMemoryStore(), nil))
	defer srv.Close()

	resp := doRequest(t, "POST", srv.URL+"/v1/instances", "tok",
		`{"name":"v","subnet":{"id":"nonexistent"}}`)
	mustStatus(t, resp, 404)
}

func TestVSI_Create_BadRequest(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(store.NewMemoryStore(), nil))
	defer srv.Close()

	// missing name
	resp := doRequest(t, "POST", srv.URL+"/v1/instances", "tok", `{"subnet":{"id":"x"}}`)
	mustStatus(t, resp, 400)

	// missing subnet
	resp = doRequest(t, "POST", srv.URL+"/v1/instances", "tok", `{"name":"v"}`)
	mustStatus(t, resp, 400)
}

func TestVSI_IPAllocation(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(store.NewMemoryStore(), nil))
	defer srv.Close()

	token := "test-token"

	resp := doRequest(t, "POST", srv.URL+"/v1/vpcs", token, `{"name":"vpc"}`)
	mustStatus(t, resp, 201)
	var vpc model.VPC
	json.NewDecoder(resp.Body).Decode(&vpc)

	resp = doRequest(t, "POST", srv.URL+"/v1/subnets", token,
		fmt.Sprintf(`{"name":"sub","vpc":{"id":"%s"}}`, vpc.ID))
	mustStatus(t, resp, 201)
	var sub model.Subnet
	json.NewDecoder(resp.Body).Decode(&sub)

	_, network, err := net.ParseCIDR(sub.CIDRBlock)
	if err != nil {
		t.Fatalf("subnet has invalid CIDR %s: %v", sub.CIDRBlock, err)
	}

	mkVSI := func(name string) model.VSI {
		body := fmt.Sprintf(`{"name":"%s","subnet":{"id":"%s"}}`, name, sub.ID)
		r := doRequest(t, "POST", srv.URL+"/v1/instances", token, body)
		mustStatus(t, r, 201)
		var v model.VSI
		json.NewDecoder(r.Body).Decode(&v)
		return v
	}

	vsi1 := mkVSI("v1")
	vsi2 := mkVSI("v2")

	// each IP must be within the subnet CIDR
	for _, vsi := range []model.VSI{vsi1, vsi2} {
		ip := net.ParseIP(vsi.PrimaryIP)
		if ip == nil {
			t.Fatalf("VSI %s has unparseable IP %q", vsi.ID, vsi.PrimaryIP)
		}
		if !network.Contains(ip) {
			t.Errorf("VSI %s IP %s is not in subnet CIDR %s", vsi.ID, vsi.PrimaryIP, sub.CIDRBlock)
		}
	}

	// IPs must be unique within the subnet
	if vsi1.PrimaryIP == vsi2.PrimaryIP {
		t.Errorf("two VSIs in the same subnet got the same IP %s", vsi1.PrimaryIP)
	}
}
