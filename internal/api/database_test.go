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

func okGetDatabase(_ context.Context, _, _ string) (model.Database, bool, error) {
	return model.Database{ID: "db-1", Name: "db", Status: "available", Engine: "postgres", Version: "15", Plan: "small"}, true, nil
}

func TestDatabase_MissingAuth(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(store.NewMemoryStore(), nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "GET", srv.URL+"/v1/databases", "", ""), 401)
}

func TestDatabase_Lifecycle(t *testing.T) {
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

	base := srv.URL + "/v1/databases"

	// empty list
	resp = doRequest(t, "GET", base, token, "")
	mustStatus(t, resp, 200)
	var listResp struct {
		DBs []model.Database `json:"databases"`
	}
	json.NewDecoder(resp.Body).Decode(&listResp)
	if len(listResp.DBs) != 0 {
		t.Fatal("expected empty list")
	}

	// create
	body := fmt.Sprintf(`{"name":"my-db","engine":"postgres","version":"15","plan":"medium","subnet_ids":["%s"]}`, sub.ID)
	resp = doRequest(t, "POST", base, token, body)
	mustStatus(t, resp, 201)
	var db model.Database
	json.NewDecoder(resp.Body).Decode(&db)
	if db.ID == "" || db.Name != "my-db" || db.Engine != "postgres" || db.Version != "15" || db.Plan != "medium" {
		t.Fatalf("unexpected db: %+v", db)
	}
	if db.Status != "available" || db.CRN == "" {
		t.Fatalf("unexpected db fields: %+v", db)
	}
	if len(db.SubnetIDs) != 1 || db.SubnetIDs[0] != sub.ID {
		t.Fatalf("unexpected subnet_ids: %+v", db.SubnetIDs)
	}

	// get
	resp = doRequest(t, "GET", base+"/"+db.ID, token, "")
	mustStatus(t, resp, 200)

	// list has 1
	resp = doRequest(t, "GET", base, token, "")
	json.NewDecoder(resp.Body).Decode(&listResp)
	if len(listResp.DBs) != 1 {
		t.Fatalf("expected 1, got %d", len(listResp.DBs))
	}

	// token isolation
	resp = doRequest(t, "GET", base+"/"+db.ID, "other-token", "")
	mustStatus(t, resp, 404)

	// rename
	resp = doRequest(t, "PATCH", base+"/"+db.ID, token, `{"name":"renamed-db"}`)
	mustStatus(t, resp, 200)
	var renamed model.Database
	json.NewDecoder(resp.Body).Decode(&renamed)
	if renamed.Name != "renamed-db" {
		t.Fatalf("expected renamed-db, got %q", renamed.Name)
	}

	// delete
	resp = doRequest(t, "DELETE", base+"/"+db.ID, token, "")
	mustStatus(t, resp, 204)

	// gone
	resp = doRequest(t, "GET", base+"/"+db.ID, token, "")
	mustStatus(t, resp, 404)
}

func TestDatabase_Create_BadRequest(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(store.NewMemoryStore(), nil))
	defer srv.Close()
	token := "tok"
	base := srv.URL + "/v1/databases"

	cases := []string{
		`{}`,
		`{bad json`,
		`{"name":"db","engine":"oracle","version":"1","plan":"small","subnet_ids":["x"]}`,
		`{"name":"db","engine":"postgres","plan":"small","subnet_ids":["x"]}`,
		`{"name":"db","engine":"postgres","version":"15","plan":"gigantic","subnet_ids":["x"]}`,
		`{"name":"db","engine":"postgres","version":"15","plan":"small"}`,
		`{"name":"db","engine":"postgres","version":"15","plan":"small","subnet_ids":[]}`,
		`{"name":"db","engine":"postgres","version":"15","plan":"small","subnet_ids":["nonexistent"]}`,
	}
	for _, body := range cases {
		resp := doRequest(t, "POST", base, token, body)
		if resp.StatusCode < 400 {
			t.Errorf("body %q: expected error, got %d", body, resp.StatusCode)
		}
	}
}

func TestDatabase_AllEngines(t *testing.T) {
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

	for _, engine := range []string{"postgres", "mysql", "mariadb"} {
		body := fmt.Sprintf(`{"name":"db-%s","engine":"%s","version":"8","plan":"small","subnet_ids":["%s"]}`, engine, engine, sub.ID)
		resp := doRequest(t, "POST", srv.URL+"/v1/databases", token, body)
		mustStatus(t, resp, 201)
	}
}

func TestDatabase_Delete_NotFound(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(store.NewMemoryStore(), nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "DELETE", srv.URL+"/v1/databases/nonexistent", "tok", ""), 404)
}

func TestDatabase_Rename_NotFound(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(store.NewMemoryStore(), nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "PATCH", srv.URL+"/v1/databases/nonexistent", "tok", `{"name":"x"}`), 404)
}

func TestDatabase_Rename_BadRequest(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(store.NewMemoryStore(), nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "PATCH", srv.URL+"/v1/databases/x", "tok", `{}`), 400)
}

// --- error path tests ---

func TestDatabase_List_StoreError(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(newErrStore(), nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "GET", srv.URL+"/v1/databases", "tok", ""), 500)
}

func TestDatabase_Create_GetSubnetStoreError(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(newErrStore(), nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "POST", srv.URL+"/v1/databases", "tok",
		`{"name":"db","engine":"postgres","version":"15","plan":"small","subnet_ids":["sub-1"]}`), 500)
}

func TestDatabase_Create_CreateStoreError(t *testing.T) {
	fs := newErrStore()
	fs.getSubnet = okGetSubnet
	srv := httptest.NewServer(api.NewServer(fs, nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "POST", srv.URL+"/v1/databases", "tok",
		`{"name":"db","engine":"postgres","version":"15","plan":"small","subnet_ids":["sub-1"]}`), 500)
}

func TestDatabase_Get_StoreError(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(newErrStore(), nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "GET", srv.URL+"/v1/databases/db-1", "tok", ""), 500)
}

func TestDatabase_Delete_GetStoreError(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(newErrStore(), nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "DELETE", srv.URL+"/v1/databases/db-1", "tok", ""), 500)
}

func TestDatabase_Delete_DeleteStoreError(t *testing.T) {
	fs := newErrStore()
	fs.getDatabase = okGetDatabase
	srv := httptest.NewServer(api.NewServer(fs, nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "DELETE", srv.URL+"/v1/databases/db-1", "tok", ""), 500)
}

func TestDatabase_Rename_GetStoreError(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(newErrStore(), nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "PATCH", srv.URL+"/v1/databases/db-1", "tok", `{"name":"x"}`), 500)
}

func TestDatabase_Rename_RenameStoreError(t *testing.T) {
	fs := newErrStore()
	fs.getDatabase = okGetDatabase
	srv := httptest.NewServer(api.NewServer(fs, nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "PATCH", srv.URL+"/v1/databases/db-1", "tok", `{"name":"x"}`), 500)
}
