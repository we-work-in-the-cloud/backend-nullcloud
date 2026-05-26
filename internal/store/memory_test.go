package store_test

import (
	"context"
	"testing"
	"time"

	"github.com/we-work-in-the-cloud/nullcloud/backend/internal/model"
	"github.com/we-work-in-the-cloud/nullcloud/backend/internal/store"
)

func TestMemoryStore_VPCLifecycle(t *testing.T) {
	s := store.NewMemoryStore()
	ctx := context.Background()
	token := "tok1"

	vpcs, err := s.ListVPCs(ctx, token)
	if err != nil || len(vpcs) != 0 {
		t.Fatalf("expected empty list, got %v %v", vpcs, err)
	}

	vpc := model.VPC{ID: "vpc-1", Name: "test", Status: "available", CreatedAt: time.Now()}
	if err := s.CreateVPC(ctx, token, vpc); err != nil {
		t.Fatal(err)
	}

	got, ok, err := s.GetVPC(ctx, token, "vpc-1")
	if err != nil || !ok || got.Name != "test" {
		t.Fatalf("unexpected: %v %v %v", got, ok, err)
	}

	_, ok2, _ := s.GetVPC(ctx, "other-token", "vpc-1")
	if ok2 {
		t.Fatal("other token should not see resource")
	}

	list, _ := s.ListVPCs(ctx, token)
	if len(list) != 1 {
		t.Fatalf("expected 1 VPC, got %d", len(list))
	}

	if err := s.DeleteVPC(ctx, token, "vpc-1"); err != nil {
		t.Fatal(err)
	}

	_, ok3, _ := s.GetVPC(ctx, token, "vpc-1")
	if ok3 {
		t.Fatal("expected deleted")
	}
}

func TestMemoryStore_SubnetLifecycle(t *testing.T) {
	s := store.NewMemoryStore()
	ctx := context.Background()
	token := "tok1"

	list, _ := s.ListSubnets(ctx, token)
	if len(list) != 0 {
		t.Fatal("expected empty")
	}

	sub := model.Subnet{ID: "sub-1", Name: "s1", VPCID: "vpc-1", Status: "available", CreatedAt: time.Now()}
	s.CreateSubnet(ctx, token, sub)

	got, ok, _ := s.GetSubnet(ctx, token, "sub-1")
	if !ok || got.VPCID != "vpc-1" {
		t.Fatal("unexpected get result")
	}

	_, ok2, _ := s.GetSubnet(ctx, "other", "sub-1")
	if ok2 {
		t.Fatal("token isolation failed")
	}

	s.DeleteSubnet(ctx, token, "sub-1")
	_, ok3, _ := s.GetSubnet(ctx, token, "sub-1")
	if ok3 {
		t.Fatal("expected deleted")
	}
}

func TestMemoryStore_VSILifecycle(t *testing.T) {
	s := store.NewMemoryStore()
	ctx := context.Background()
	token := "tok1"

	list, _ := s.ListVSIs(ctx, token)
	if len(list) != 0 {
		t.Fatal("expected empty")
	}

	vsi := model.VSI{ID: "vsi-1", Name: "v1", SubnetID: "sub-1", VPCID: "vpc-1", Status: "running", CreatedAt: time.Now()}
	s.CreateVSI(ctx, token, vsi)

	got, ok, _ := s.GetVSI(ctx, token, "vsi-1")
	if !ok || got.SubnetID != "sub-1" {
		t.Fatal("unexpected get result")
	}

	_, ok2, _ := s.GetVSI(ctx, "other", "vsi-1")
	if ok2 {
		t.Fatal("token isolation failed")
	}

	s.DeleteVSI(ctx, token, "vsi-1")
	_, ok3, _ := s.GetVSI(ctx, token, "vsi-1")
	if ok3 {
		t.Fatal("expected deleted")
	}
}
