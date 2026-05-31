package api_test

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/we-work-in-the-cloud/nullcloud/backend/internal/api"
	"github.com/we-work-in-the-cloud/nullcloud/backend/internal/model"
	"github.com/we-work-in-the-cloud/nullcloud/backend/internal/store"
)

func okGetBucket(_ context.Context, _, _ string) (model.Bucket, bool, error) {
	return model.Bucket{ID: "bkt-1", Name: "bucket", Status: "available", Region: "us-east-1"}, true, nil
}

func TestBucket_MissingAuth(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(store.NewMemoryStore(), nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "GET", srv.URL+"/v1/buckets", "", ""), 401)
}

func TestBucket_Lifecycle(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(store.NewMemoryStore(), nil))
	defer srv.Close()

	token := "test-token"
	base := srv.URL + "/v1/buckets"

	// empty list
	resp := doRequest(t, "GET", base, token, "")
	mustStatus(t, resp, 200)
	var listResp struct {
		Buckets []model.Bucket `json:"buckets"`
	}
	json.NewDecoder(resp.Body).Decode(&listResp)
	if len(listResp.Buckets) != 0 {
		t.Fatal("expected empty list")
	}

	// create with explicit region
	resp = doRequest(t, "POST", base, token, `{"name":"my-bucket","region":"eu-west-1"}`)
	mustStatus(t, resp, 201)
	var bkt model.Bucket
	json.NewDecoder(resp.Body).Decode(&bkt)
	if bkt.ID == "" || bkt.Name != "my-bucket" || bkt.Region != "eu-west-1" || bkt.Status != "available" || bkt.CRN == "" {
		t.Fatalf("unexpected bucket: %+v", bkt)
	}

	// get
	resp = doRequest(t, "GET", base+"/"+bkt.ID, token, "")
	mustStatus(t, resp, 200)

	// list has 1
	resp = doRequest(t, "GET", base, token, "")
	json.NewDecoder(resp.Body).Decode(&listResp)
	if len(listResp.Buckets) != 1 {
		t.Fatalf("expected 1, got %d", len(listResp.Buckets))
	}

	// token isolation
	resp = doRequest(t, "GET", base+"/"+bkt.ID, "other-token", "")
	mustStatus(t, resp, 404)

	// rename
	resp = doRequest(t, "PATCH", base+"/"+bkt.ID, token, `{"name":"renamed-bucket"}`)
	mustStatus(t, resp, 200)
	var renamed model.Bucket
	json.NewDecoder(resp.Body).Decode(&renamed)
	if renamed.Name != "renamed-bucket" {
		t.Fatalf("expected renamed-bucket, got %q", renamed.Name)
	}

	// delete
	resp = doRequest(t, "DELETE", base+"/"+bkt.ID, token, "")
	mustStatus(t, resp, 204)

	// gone
	resp = doRequest(t, "GET", base+"/"+bkt.ID, token, "")
	mustStatus(t, resp, 404)
}

func TestBucket_Create_DefaultRegion(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(store.NewMemoryStore(), nil))
	defer srv.Close()

	resp := doRequest(t, "POST", srv.URL+"/v1/buckets", "tok", `{"name":"my-bucket"}`)
	mustStatus(t, resp, 201)
	var bkt model.Bucket
	json.NewDecoder(resp.Body).Decode(&bkt)
	if bkt.Region != "us-east-1" {
		t.Fatalf("expected default region us-east-1, got %q", bkt.Region)
	}
}

func TestBucket_Create_BadRequest(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(store.NewMemoryStore(), nil))
	defer srv.Close()
	token := "tok"
	base := srv.URL + "/v1/buckets"

	cases := []string{`{}`, `{bad json`}
	for _, body := range cases {
		resp := doRequest(t, "POST", base, token, body)
		mustStatus(t, resp, 400)
	}
}

func TestBucket_Delete_NotFound(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(store.NewMemoryStore(), nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "DELETE", srv.URL+"/v1/buckets/nonexistent", "tok", ""), 404)
}

func TestBucket_Rename_NotFound(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(store.NewMemoryStore(), nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "PATCH", srv.URL+"/v1/buckets/nonexistent", "tok", `{"name":"x"}`), 404)
}

func TestBucket_Rename_BadRequest(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(store.NewMemoryStore(), nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "PATCH", srv.URL+"/v1/buckets/x", "tok", `{}`), 400)
}

// --- error path tests ---

func TestBucket_List_StoreError(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(newErrStore(), nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "GET", srv.URL+"/v1/buckets", "tok", ""), 500)
}

func TestBucket_Create_StoreError(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(newErrStore(), nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "POST", srv.URL+"/v1/buckets", "tok", `{"name":"b"}`), 500)
}

func TestBucket_Get_StoreError(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(newErrStore(), nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "GET", srv.URL+"/v1/buckets/bkt-1", "tok", ""), 500)
}

func TestBucket_Delete_GetStoreError(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(newErrStore(), nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "DELETE", srv.URL+"/v1/buckets/bkt-1", "tok", ""), 500)
}

func TestBucket_Delete_DeleteStoreError(t *testing.T) {
	fs := newErrStore()
	fs.getBucket = okGetBucket
	srv := httptest.NewServer(api.NewServer(fs, nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "DELETE", srv.URL+"/v1/buckets/bkt-1", "tok", ""), 500)
}

func TestBucket_Rename_GetStoreError(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(newErrStore(), nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "PATCH", srv.URL+"/v1/buckets/bkt-1", "tok", `{"name":"x"}`), 500)
}

func TestBucket_Rename_RenameStoreError(t *testing.T) {
	fs := newErrStore()
	fs.getBucket = okGetBucket
	srv := httptest.NewServer(api.NewServer(fs, nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "PATCH", srv.URL+"/v1/buckets/bkt-1", "tok", `{"name":"x"}`), 500)
}
