package api_test

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/we-work-in-the-cloud/nullcloud/backend/internal/api"
	"github.com/we-work-in-the-cloud/nullcloud/backend/internal/model"
	"github.com/we-work-in-the-cloud/nullcloud/backend/internal/store"
)

func TestRegions_MissingAuth(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(store.NewMemoryStore(), nil))
	defer srv.Close()
	mustStatus(t, doRequest(t, "GET", srv.URL+"/v1/regions", "", ""), 401)
}

func TestRegions_List(t *testing.T) {
	srv := httptest.NewServer(api.NewServer(store.NewMemoryStore(), nil))
	defer srv.Close()

	resp := doRequest(t, "GET", srv.URL+"/v1/regions", "tok", "")
	mustStatus(t, resp, 200)

	var result struct {
		Regions []model.Region `json:"regions"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(result.Regions) != 6 {
		t.Fatalf("expected 6 regions, got %d", len(result.Regions))
	}

	// verify all expected regions are present
	names := make(map[string]bool)
	for _, r := range result.Regions {
		names[r.Name] = true
		if len(r.Zones) != 3 {
			t.Errorf("region %s: expected 3 zones, got %d", r.Name, len(r.Zones))
		}
		// verify zones are named correctly (e.g. us-east-1, us-east-2, us-east-3)
		for i, z := range r.Zones {
			expected := r.Name + "-" + string(rune('1'+i))
			if z.Name != expected {
				t.Errorf("region %s zone %d: expected %s, got %s", r.Name, i, expected, z.Name)
			}
		}
	}
	for _, want := range []string{"us-east", "us-west", "us-central", "eu-central", "eu-east", "eu-west"} {
		if !names[want] {
			t.Errorf("region %s not found in response", want)
		}
	}
}
