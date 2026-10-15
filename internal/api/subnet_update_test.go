package api_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/we-work-in-the-cloud/nullcloud/backend/internal/api"
	"github.com/we-work-in-the-cloud/nullcloud/backend/internal/model"
	"github.com/we-work-in-the-cloud/nullcloud/backend/internal/store"
)

func TestSubnet_CIDRChange_CombinedUpdate(t *testing.T) {
	for _, storeType := range []string{"memory", "jsonfile"} {
		for _, tc := range []struct {
			name   string
			cidr   string
			status int
		}{
			{"outside_range", "192.168.1.0/24", 400},
			{"invalid_cidr", "invalid", 400},
			{"compatible_range", "10.0.0.0/23", 200},
			{"unchanged_range", "10.0.0.0/24", 200},
		} {
			t.Run(storeType+"/"+tc.name, func(t *testing.T) {
				ctx := context.Background()
				const token = "tok"
				path := filepath.Join(t.TempDir(), "state.json")
				var s store.Store = store.NewMemoryStore()
				if storeType == "jsonfile" {
					var err error
					s, err = store.NewJSONFileStore(path)
					if err != nil {
						t.Fatal(err)
					}
				}
				sub := model.Subnet{ID: "sub-1", Name: "original", CIDRBlock: "10.0.0.0/24"}
				if err := s.CreateSubnet(ctx, token, sub); err != nil {
					t.Fatal(err)
				}
				for _, seed := range []struct {
					token string
					vsi   model.VSI
				}{
					{token, model.VSI{ID: "vsi-1", SubnetID: sub.ID, PrimaryIP: "10.0.0.10"}},
					{token, model.VSI{ID: "vsi-2", SubnetID: "other-subnet", PrimaryIP: "192.168.1.10"}},
					{"other-token", model.VSI{ID: "vsi-3", SubnetID: sub.ID, PrimaryIP: "192.168.1.20"}},
				} {
					if err := s.CreateVSI(ctx, seed.token, seed.vsi); err != nil {
						t.Fatal(err)
					}
				}
				srv := httptest.NewServer(api.NewServer(s, nil))
				defer srv.Close()
				body := fmt.Sprintf(`{"name":"renamed","cidr_block":%q}`, tc.cidr)
				resp := doRequest(t, "PATCH", srv.URL+"/v1/subnets/"+sub.ID, token, body)
				mustStatus(t, resp, tc.status)
				resp.Body.Close()

				wantName, wantCIDR := sub.Name, sub.CIDRBlock
				if tc.status == 200 {
					wantName, wantCIDR = "renamed", tc.cidr
				}
				resp = doRequest(t, "GET", srv.URL+"/v1/subnets/"+sub.ID, token, "")
				mustStatus(t, resp, 200)
				defer resp.Body.Close()
				var got model.Subnet
				if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
					t.Fatal(err)
				}
				if got.Name != wantName || got.CIDRBlock != wantCIDR {
					t.Fatalf("subnet after PATCH: name=%q CIDR=%q; want name=%q CIDR=%q", got.Name, got.CIDRBlock, wantName, wantCIDR)
				}
				if storeType == "jsonfile" {
					reloaded, err := store.NewJSONFileStore(path)
					if err != nil {
						t.Fatal(err)
					}
					persisted, ok, err := reloaded.GetSubnet(ctx, token, sub.ID)
					if err != nil || !ok || persisted.Name != wantName || persisted.CIDRBlock != wantCIDR {
						t.Fatalf("persisted subnet: %+v, ok=%v, err=%v", persisted, ok, err)
					}
				}
			})
		}
	}
}
