package api_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/we-work-in-the-cloud/nullcloud/backend/internal/api"
	"github.com/we-work-in-the-cloud/nullcloud/backend/internal/model"
	"github.com/we-work-in-the-cloud/nullcloud/backend/internal/store"
)

func okGetLoadBalancer(_ context.Context, _, _ string) (model.LoadBalancer, bool, error) {
	return model.LoadBalancer{ID: "lb-1", Name: "lb", Status: "active", Protocol: "http", Port: 80}, true, nil
}

func TestLoadBalancer_MissingAuth(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(store.NewMemoryStore(), nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "GET", srv.URL+"/v1/loadbalancers", "", ""), 401)
}

func TestLoadBalancer_Lifecycle(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(store.NewMemoryStore(), nil))
	defer srv.Close()

	token := "test-token"
	base := srv.URL + "/v1/loadbalancers"

	// empty list
	resp := doRequest(t, "GET", base, token, "")
	mustStatus(t, resp, 200)
	var listResp struct {
		LBs []model.LoadBalancer `json:"load_balancers"`
	}
	json.NewDecoder(resp.Body).Decode(&listResp)
	if len(listResp.LBs) != 0 {
		t.Fatal("expected empty list")
	}

	// create
	resp = doRequest(t, "POST", base, token, `{"name":"my-lb","protocol":"https","port":443}`)
	mustStatus(t, resp, 201)
	var lb model.LoadBalancer
	json.NewDecoder(resp.Body).Decode(&lb)
	if lb.ID == "" || lb.Name != "my-lb" || lb.Protocol != "https" || lb.Port != 443 || lb.Status != "active" || lb.CRN == "" {
		t.Fatalf("unexpected lb: %+v", lb)
	}
	if lb.Targets != nil {
		t.Fatalf("expected nil targets, got %+v", lb.Targets)
	}

	// get
	resp = doRequest(t, "GET", base+"/"+lb.ID, token, "")
	mustStatus(t, resp, 200)

	// list has 1
	resp = doRequest(t, "GET", base, token, "")
	json.NewDecoder(resp.Body).Decode(&listResp)
	if len(listResp.LBs) != 1 {
		t.Fatalf("expected 1, got %d", len(listResp.LBs))
	}

	// token isolation
	resp = doRequest(t, "GET", base+"/"+lb.ID, "other-token", "")
	mustStatus(t, resp, 404)

	// rename
	resp = doRequest(t, "PATCH", base+"/"+lb.ID, token, `{"name":"renamed-lb"}`)
	mustStatus(t, resp, 200)
	var renamed model.LoadBalancer
	json.NewDecoder(resp.Body).Decode(&renamed)
	if renamed.Name != "renamed-lb" {
		t.Fatalf("expected renamed-lb, got %q", renamed.Name)
	}

	// delete
	resp = doRequest(t, "DELETE", base+"/"+lb.ID, token, "")
	mustStatus(t, resp, 204)

	// gone
	resp = doRequest(t, "GET", base+"/"+lb.ID, token, "")
	mustStatus(t, resp, 404)
}

func TestLoadBalancer_WithVSITarget(t *testing.T) {
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

	resp = doRequest(t, "POST", srv.URL+"/v1/instances", token,
		fmt.Sprintf(`{"name":"vsi","subnet":{"id":"%s"}}`, sub.ID))
	mustStatus(t, resp, 201)
	var vsi model.VSI
	json.NewDecoder(resp.Body).Decode(&vsi)

	body := fmt.Sprintf(`{"name":"my-lb","protocol":"tcp","port":80,"targets":[{"type":"vsi","id":"%s"}]}`, vsi.ID)
	resp = doRequest(t, "POST", srv.URL+"/v1/loadbalancers", token, body)
	mustStatus(t, resp, 201)
	var lb model.LoadBalancer
	json.NewDecoder(resp.Body).Decode(&lb)
	if len(lb.Targets) != 1 || lb.Targets[0].Type != "vsi" || lb.Targets[0].ID != vsi.ID {
		t.Fatalf("unexpected targets: %+v", lb.Targets)
	}
}

func TestLoadBalancer_WithClusterTarget(t *testing.T) {
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

	resp = doRequest(t, "POST", srv.URL+"/v1/clusters", token,
		fmt.Sprintf(`{"name":"cluster","version":"1.30","node_count":2,"subnet_ids":["%s"]}`, sub.ID))
	mustStatus(t, resp, 201)
	var cluster model.KubernetesCluster
	json.NewDecoder(resp.Body).Decode(&cluster)

	body := fmt.Sprintf(`{"name":"my-lb","protocol":"http","port":80,"targets":[{"type":"cluster","id":"%s"}]}`, cluster.ID)
	resp = doRequest(t, "POST", srv.URL+"/v1/loadbalancers", token, body)
	mustStatus(t, resp, 201)
	var lb model.LoadBalancer
	json.NewDecoder(resp.Body).Decode(&lb)
	if len(lb.Targets) != 1 || lb.Targets[0].Type != "cluster" || lb.Targets[0].ID != cluster.ID {
		t.Fatalf("unexpected targets: %+v", lb.Targets)
	}
}

func TestLoadBalancer_Create_BadRequest(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(store.NewMemoryStore(), nil))
	defer srv.Close()
	token := "tok"
	base := srv.URL + "/v1/loadbalancers"

	cases := []string{
		`{}`,
		`{"name":"lb"}`,
		`{"name":"lb","protocol":"tcp"}`,
		`{"name":"lb","protocol":"tcp","port":0}`,
		`{"name":"lb","protocol":"tcp","port":70000}`,
		`{"name":"lb","protocol":"ftp","port":80}`,
		`{bad json`,
		`{"name":"lb","protocol":"http","port":80,"targets":[{"type":"bad","id":"x"}]}`,
		`{"name":"lb","protocol":"http","port":80,"targets":[{"type":"vsi","id":"nonexistent"}]}`,
		`{"name":"lb","protocol":"http","port":80,"targets":[{"type":"cluster","id":"nonexistent"}]}`,
	}
	for _, body := range cases {
		resp := doRequest(t, "POST", base, token, body)
		if resp.StatusCode < 400 {
			t.Errorf("body %q: expected error, got %d", body, resp.StatusCode)
		}
	}
}

func TestLoadBalancer_Delete_NotFound(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(store.NewMemoryStore(), nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "DELETE", srv.URL+"/v1/loadbalancers/nonexistent", "tok", ""), 404)
}

func TestLoadBalancer_Rename_NotFound(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(store.NewMemoryStore(), nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "PATCH", srv.URL+"/v1/loadbalancers/nonexistent", "tok", `{"name":"x"}`), 404)
}

func TestLoadBalancer_Rename_BadRequest(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(store.NewMemoryStore(), nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "PATCH", srv.URL+"/v1/loadbalancers/x", "tok", `{}`), 400)
}

// --- error path tests using funcStore ---

func TestLB_List_StoreError(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(newErrStore(), nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "GET", srv.URL+"/v1/loadbalancers", "tok", ""), 500)
}

func TestLB_Create_StoreError(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(newErrStore(), nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "POST", srv.URL+"/v1/loadbalancers", "tok",
		`{"name":"lb","protocol":"tcp","port":80}`), 500)
}

func TestLB_Get_StoreError(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(newErrStore(), nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "GET", srv.URL+"/v1/loadbalancers/lb-1", "tok", ""), 500)
}

func TestLB_Delete_GetStoreError(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(newErrStore(), nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "DELETE", srv.URL+"/v1/loadbalancers/lb-1", "tok", ""), 500)
}

func TestLB_Delete_DeleteStoreError(t *testing.T) {
	fs := newErrStore()
	fs.getLoadBalancer = okGetLoadBalancer
	srv := httptest.NewServer(api.NewServer(fs, nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "DELETE", srv.URL+"/v1/loadbalancers/lb-1", "tok", ""), 500)
}

func TestLB_Rename_GetStoreError(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(newErrStore(), nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "PATCH", srv.URL+"/v1/loadbalancers/lb-1", "tok", `{"name":"x"}`), 500)
}

func TestLB_Rename_RenameStoreError(t *testing.T) {
	fs := newErrStore()
	fs.getLoadBalancer = okGetLoadBalancer
	srv := httptest.NewServer(api.NewServer(fs, nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "PATCH", srv.URL+"/v1/loadbalancers/lb-1", "tok", `{"name":"x"}`), 500)
}

func TestLB_Create_GetVSIStoreError(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(newErrStore(), nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "POST", srv.URL+"/v1/loadbalancers", "tok",
		`{"name":"lb","protocol":"tcp","port":80,"targets":[{"type":"vsi","id":"vsi-1"}]}`), 500)
}

func TestLB_Create_GetClusterStoreError(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(newErrStore(), nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "POST", srv.URL+"/v1/loadbalancers", "tok",
		`{"name":"lb","protocol":"tcp","port":80,"targets":[{"type":"cluster","id":"k8s-1"}]}`), 500)
}
