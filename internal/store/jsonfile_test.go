package store_test

import (
	"context"
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
