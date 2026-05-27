package store_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/we-work-in-the-cloud/nullcloud/backend/internal/model"
	"github.com/we-work-in-the-cloud/nullcloud/backend/internal/store"
)

func TestJSONFileStore_NewFileOK(t *testing.T) {
	path := filepath.Join(t.TempDir(), "new.json")
	_, err := store.NewJSONFileStore(path)
	if err != nil {
		t.Fatal(err)
	}
}

func TestJSONFileStore_VPCPersistence(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	ctx := context.Background()
	token := "tok1"

	s1, err := store.NewJSONFileStore(path)
	if err != nil {
		t.Fatal(err)
	}

	vpc := model.VPC{ID: "vpc-1", Name: "persisted", Status: "available", CreatedAt: time.Now().UTC()}
	if err := s1.CreateVPC(ctx, token, vpc); err != nil {
		t.Fatal(err)
	}

	s2, err := store.NewJSONFileStore(path)
	if err != nil {
		t.Fatal(err)
	}

	got, ok, err := s2.GetVPC(ctx, token, "vpc-1")
	if err != nil || !ok {
		t.Fatalf("expected to find vpc after reload: ok=%v err=%v", ok, err)
	}
	if got.Name != "persisted" {
		t.Fatalf("unexpected name: %s", got.Name)
	}
}

func TestJSONFileStore_SubnetPersistence(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	ctx := context.Background()
	token := "tok1"

	s1, _ := store.NewJSONFileStore(path)
	sub := model.Subnet{ID: "sub-1", Name: "s1", VPCID: "vpc-1", Status: "available", CreatedAt: time.Now().UTC()}
	s1.CreateSubnet(ctx, token, sub)

	s2, _ := store.NewJSONFileStore(path)
	got, ok, _ := s2.GetSubnet(ctx, token, "sub-1")
	if !ok || got.VPCID != "vpc-1" {
		t.Fatal("subnet not persisted")
	}
}

func TestJSONFileStore_TokenIsolation(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	ctx := context.Background()

	s, _ := store.NewJSONFileStore(path)
	vpc := model.VPC{ID: "vpc-1", Name: "tok1-vpc", Status: "available", CreatedAt: time.Now().UTC()}
	s.CreateVPC(ctx, "tok1", vpc)

	_, ok, _ := s.GetVPC(ctx, "tok2", "vpc-1")
	if ok {
		t.Fatal("token isolation failed")
	}
}

func TestJSONFileStore_InvalidJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bad.json")
	if err := os.WriteFile(path, []byte(`{bad json`), 0644); err != nil {
		t.Fatal(err)
	}
	_, err := store.NewJSONFileStore(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestJSONFileStore_VPCListAndDelete(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	ctx := context.Background()
	token := "tok1"

	s, _ := store.NewJSONFileStore(path)
	s.CreateVPC(ctx, token, model.VPC{ID: "vpc-1", Name: "a", Status: "available", CreatedAt: time.Now().UTC()})
	s.CreateVPC(ctx, token, model.VPC{ID: "vpc-2", Name: "b", Status: "available", CreatedAt: time.Now().UTC()})

	list, err := s.ListVPCs(ctx, token)
	if err != nil || len(list) != 2 {
		t.Fatalf("expected 2 VPCs, got %d err=%v", len(list), err)
	}

	if err := s.DeleteVPC(ctx, token, "vpc-1"); err != nil {
		t.Fatal(err)
	}
	list, _ = s.ListVPCs(ctx, token)
	if len(list) != 1 {
		t.Fatalf("expected 1 VPC after delete, got %d", len(list))
	}
}

func TestJSONFileStore_SubnetListAndDelete(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	ctx := context.Background()
	token := "tok1"

	s, _ := store.NewJSONFileStore(path)
	s.CreateSubnet(ctx, token, model.Subnet{ID: "sub-1", Name: "s1", VPCID: "v1", Status: "available", CreatedAt: time.Now().UTC()})
	s.CreateSubnet(ctx, token, model.Subnet{ID: "sub-2", Name: "s2", VPCID: "v1", Status: "available", CreatedAt: time.Now().UTC()})

	list, err := s.ListSubnets(ctx, token)
	if err != nil || len(list) != 2 {
		t.Fatalf("expected 2 subnets, got %d err=%v", len(list), err)
	}

	if err := s.DeleteSubnet(ctx, token, "sub-1"); err != nil {
		t.Fatal(err)
	}
	list, _ = s.ListSubnets(ctx, token)
	if len(list) != 1 {
		t.Fatalf("expected 1 subnet after delete, got %d", len(list))
	}
}

func TestJSONFileStore_VSILifecycle(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	ctx := context.Background()
	token := "tok1"

	s, _ := store.NewJSONFileStore(path)

	vsi := model.VSI{ID: "vsi-1", Name: "v1", SubnetID: "sub-1", VPCID: "vpc-1", Status: "running", CreatedAt: time.Now().UTC()}
	if err := s.CreateVSI(ctx, token, vsi); err != nil {
		t.Fatal(err)
	}

	got, ok, err := s.GetVSI(ctx, token, "vsi-1")
	if err != nil || !ok || got.SubnetID != "sub-1" {
		t.Fatalf("unexpected get: ok=%v err=%v", ok, err)
	}

	s.CreateVSI(ctx, token, model.VSI{ID: "vsi-2", Name: "v2", Status: "running", CreatedAt: time.Now().UTC()})
	list, err := s.ListVSIs(ctx, token)
	if err != nil || len(list) != 2 {
		t.Fatalf("expected 2 VSIs, got %d err=%v", len(list), err)
	}

	if err := s.UpdateVSIStatus(ctx, token, "vsi-1", "stopped"); err != nil {
		t.Fatal(err)
	}
	got, _, _ = s.GetVSI(ctx, token, "vsi-1")
	if got.Status != "stopped" {
		t.Fatalf("expected stopped, got %q", got.Status)
	}

	if err := s.UpdateVSIStatus(ctx, token, "nonexistent", "running"); err == nil {
		t.Fatal("expected error for missing VSI")
	}

	if err := s.DeleteVSI(ctx, token, "vsi-1"); err != nil {
		t.Fatal(err)
	}
	list, _ = s.ListVSIs(ctx, token)
	if len(list) != 1 {
		t.Fatalf("expected 1 VSI after delete, got %d", len(list))
	}
}

func TestJSONFileStore_VSIPersistence(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	ctx := context.Background()
	token := "tok1"

	s1, _ := store.NewJSONFileStore(path)
	s1.CreateVSI(ctx, token, model.VSI{ID: "vsi-1", Name: "persisted", Status: "running", CreatedAt: time.Now().UTC()})

	s2, _ := store.NewJSONFileStore(path)
	got, ok, _ := s2.GetVSI(ctx, token, "vsi-1")
	if !ok || got.Name != "persisted" {
		t.Fatalf("VSI not persisted: ok=%v name=%q", ok, got.Name)
	}
}
