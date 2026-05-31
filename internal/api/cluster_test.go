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

func okGetKubernetesCluster(_ context.Context, _, _ string) (model.KubernetesCluster, bool, error) {
	return model.KubernetesCluster{ID: "k8s-1", Name: "cluster", Status: "running", Version: "1.30", NodeCount: 2}, true, nil
}

func TestCluster_MissingAuth(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(store.NewMemoryStore(), nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "GET", srv.URL+"/v1/clusters", "", ""), 401)
}

func TestCluster_Lifecycle(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(store.NewMemoryStore(), nil))
	defer srv.Close()

	token := "test-token"

	// create VPC + Subnet
	resp := doRequest(t, "POST", srv.URL+"/v1/vpcs", token, `{"name":"vpc"}`)
	mustStatus(t, resp, 201)
	var vpc model.VPC
	json.NewDecoder(resp.Body).Decode(&vpc)

	resp = doRequest(t, "POST", srv.URL+"/v1/subnets", token,
		fmt.Sprintf(`{"name":"sub","vpc":{"id":"%s"}}`, vpc.ID))
	mustStatus(t, resp, 201)
	var sub model.Subnet
	json.NewDecoder(resp.Body).Decode(&sub)

	base := srv.URL + "/v1/clusters"

	// empty list
	resp = doRequest(t, "GET", base, token, "")
	mustStatus(t, resp, 200)
	var listResp struct {
		Clusters []model.KubernetesCluster `json:"clusters"`
	}
	json.NewDecoder(resp.Body).Decode(&listResp)
	if len(listResp.Clusters) != 0 {
		t.Fatal("expected empty list")
	}

	// create
	body := fmt.Sprintf(`{"name":"my-cluster","version":"1.30","node_count":3,"subnet_ids":["%s"]}`, sub.ID)
	resp = doRequest(t, "POST", base, token, body)
	mustStatus(t, resp, 201)
	var cluster model.KubernetesCluster
	json.NewDecoder(resp.Body).Decode(&cluster)
	if cluster.ID == "" || cluster.Name != "my-cluster" || cluster.Version != "1.30" || cluster.NodeCount != 3 {
		t.Fatalf("unexpected cluster: %+v", cluster)
	}
	if cluster.Status != "running" || cluster.CRN == "" {
		t.Fatalf("unexpected cluster fields: %+v", cluster)
	}
	if len(cluster.SubnetIDs) != 1 || cluster.SubnetIDs[0] != sub.ID {
		t.Fatalf("unexpected subnet_ids: %+v", cluster.SubnetIDs)
	}

	// get
	resp = doRequest(t, "GET", base+"/"+cluster.ID, token, "")
	mustStatus(t, resp, 200)

	// list has 1
	resp = doRequest(t, "GET", base, token, "")
	json.NewDecoder(resp.Body).Decode(&listResp)
	if len(listResp.Clusters) != 1 {
		t.Fatalf("expected 1, got %d", len(listResp.Clusters))
	}

	// token isolation
	resp = doRequest(t, "GET", base+"/"+cluster.ID, "other-token", "")
	mustStatus(t, resp, 404)

	// rename
	resp = doRequest(t, "PATCH", base+"/"+cluster.ID, token, `{"name":"renamed-cluster"}`)
	mustStatus(t, resp, 200)
	var renamed model.KubernetesCluster
	json.NewDecoder(resp.Body).Decode(&renamed)
	if renamed.Name != "renamed-cluster" {
		t.Fatalf("expected renamed-cluster, got %q", renamed.Name)
	}

	// delete
	resp = doRequest(t, "DELETE", base+"/"+cluster.ID, token, "")
	mustStatus(t, resp, 204)

	// gone
	resp = doRequest(t, "GET", base+"/"+cluster.ID, token, "")
	mustStatus(t, resp, 404)
}

func TestCluster_MultipleSubnets(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(store.NewMemoryStore(), nil))
	defer srv.Close()

	token := "test-token"

	resp := doRequest(t, "POST", srv.URL+"/v1/vpcs", token, `{"name":"vpc"}`)
	mustStatus(t, resp, 201)
	var vpc model.VPC
	json.NewDecoder(resp.Body).Decode(&vpc)

	var subIDs []string
	for i := 0; i < 2; i++ {
		resp = doRequest(t, "POST", srv.URL+"/v1/subnets", token,
			fmt.Sprintf(`{"name":"sub-%d","vpc":{"id":"%s"}}`, i, vpc.ID))
		mustStatus(t, resp, 201)
		var sub model.Subnet
		json.NewDecoder(resp.Body).Decode(&sub)
		subIDs = append(subIDs, sub.ID)
	}

	body := fmt.Sprintf(`{"name":"cluster","version":"1.28","node_count":1,"subnet_ids":["%s","%s"]}`, subIDs[0], subIDs[1])
	resp = doRequest(t, "POST", srv.URL+"/v1/clusters", token, body)
	mustStatus(t, resp, 201)
	var cluster model.KubernetesCluster
	json.NewDecoder(resp.Body).Decode(&cluster)
	if len(cluster.SubnetIDs) != 2 {
		t.Fatalf("expected 2 subnet_ids, got %d", len(cluster.SubnetIDs))
	}
}

func TestCluster_Create_BadRequest(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(store.NewMemoryStore(), nil))
	defer srv.Close()
	token := "tok"
	base := srv.URL + "/v1/clusters"

	cases := []string{
		`{}`,
		`{bad json`,
		`{"name":"c","version":"1.30","node_count":0,"subnet_ids":["x"]}`,
		`{"name":"c","node_count":1,"subnet_ids":["x"]}`,
		`{"name":"c","version":"1.30","node_count":1}`,
		`{"name":"c","version":"1.30","node_count":1,"subnet_ids":[]}`,
		`{"name":"c","version":"1.30","node_count":1,"subnet_ids":["nonexistent"]}`,
	}
	for _, body := range cases {
		resp := doRequest(t, "POST", base, token, body)
		if resp.StatusCode < 400 {
			t.Errorf("body %q: expected error, got %d", body, resp.StatusCode)
		}
	}
}

func TestCluster_Delete_NotFound(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(store.NewMemoryStore(), nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "DELETE", srv.URL+"/v1/clusters/nonexistent", "tok", ""), 404)
}

func TestCluster_Rename_NotFound(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(store.NewMemoryStore(), nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "PATCH", srv.URL+"/v1/clusters/nonexistent", "tok", `{"name":"x"}`), 404)
}

func TestCluster_Rename_BadRequest(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(store.NewMemoryStore(), nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "PATCH", srv.URL+"/v1/clusters/x", "tok", `{}`), 400)
}

// --- error path tests ---

func TestCluster_List_StoreError(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(newErrStore(), nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "GET", srv.URL+"/v1/clusters", "tok", ""), 500)
}

func TestCluster_Create_GetSubnetStoreError(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(newErrStore(), nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "POST", srv.URL+"/v1/clusters", "tok",
		`{"name":"c","version":"1.30","node_count":1,"subnet_ids":["sub-1"]}`), 500)
}

func TestCluster_Create_CreateStoreError(t *testing.T) {
	fs := newErrStore()
	fs.getSubnet = okGetSubnet
	srv := httptest.NewServer(api.NewServer(fs, nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "POST", srv.URL+"/v1/clusters", "tok",
		`{"name":"c","version":"1.30","node_count":1,"subnet_ids":["sub-1"]}`), 500)
}

func TestCluster_Get_StoreError(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(newErrStore(), nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "GET", srv.URL+"/v1/clusters/k8s-1", "tok", ""), 500)
}

func TestCluster_Delete_GetStoreError(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(newErrStore(), nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "DELETE", srv.URL+"/v1/clusters/k8s-1", "tok", ""), 500)
}

func TestCluster_Delete_DeleteStoreError(t *testing.T) {
	fs := newErrStore()
	fs.getKubernetesCluster = okGetKubernetesCluster
	srv := httptest.NewServer(api.NewServer(fs, nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "DELETE", srv.URL+"/v1/clusters/k8s-1", "tok", ""), 500)
}

func TestCluster_Rename_GetStoreError(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(newErrStore(), nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "PATCH", srv.URL+"/v1/clusters/k8s-1", "tok", `{"name":"x"}`), 500)
}

func TestCluster_Rename_RenameStoreError(t *testing.T) {
	fs := newErrStore()
	fs.getKubernetesCluster = okGetKubernetesCluster
	srv := httptest.NewServer(api.NewServer(fs, nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "PATCH", srv.URL+"/v1/clusters/k8s-1", "tok", `{"name":"x"}`), 500)
}
