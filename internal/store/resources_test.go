package store_test

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/we-work-in-the-cloud/nullcloud/backend/internal/model"
	"github.com/we-work-in-the-cloud/nullcloud/backend/internal/store"
)

// ---- helpers ----

func newLB(id, name string) model.LoadBalancer {
	return model.LoadBalancer{ID: id, Name: name, Status: "active", CRN: "crn:lb:" + id,
		Protocol: "http", Port: 80, CreatedAt: time.Now().UTC()}
}

func newBucket(id, name string) model.Bucket {
	return model.Bucket{ID: id, Name: name, Status: "available", CRN: "crn:bkt:" + id,
		Region: "us-east-1", CreatedAt: time.Now().UTC()}
}

func newDB(id, name string) model.Database {
	return model.Database{ID: id, Name: name, Status: "available", CRN: "crn:db:" + id,
		Engine: "postgres", Version: "15", Plan: "small", SubnetIDs: []string{"sub-1"}, CreatedAt: time.Now().UTC()}
}

func newCluster(id, name string) model.KubernetesCluster {
	return model.KubernetesCluster{ID: id, Name: name, Status: "running", CRN: "crn:k8s:" + id,
		Version: "1.30", NodeCount: 2, SubnetIDs: []string{"sub-1"}, CreatedAt: time.Now().UTC()}
}

// ---- MemoryStore ----

func TestMemoryStore_LoadBalancerLifecycle(t *testing.T) {
	s := store.NewMemoryStore()
	ctx := context.Background()
	token := "tok1"

	list, err := s.ListLoadBalancers(ctx, token)
	if err != nil || len(list) != 0 {
		t.Fatalf("expected empty list: %v %v", list, err)
	}

	lb := newLB("lb-1", "my-lb")
	if err := s.CreateLoadBalancer(ctx, token, lb); err != nil {
		t.Fatal(err)
	}

	got, ok, err := s.GetLoadBalancer(ctx, token, "lb-1")
	if err != nil || !ok || got.Name != "my-lb" {
		t.Fatalf("unexpected get: %v %v %v", got, ok, err)
	}

	_, ok2, _ := s.GetLoadBalancer(ctx, "other-token", "lb-1")
	if ok2 {
		t.Fatal("token isolation failed")
	}

	list, _ = s.ListLoadBalancers(ctx, token)
	if len(list) != 1 {
		t.Fatalf("expected 1, got %d", len(list))
	}

	if err := s.RenameLoadBalancer(ctx, token, "lb-1", "renamed-lb"); err != nil {
		t.Fatal(err)
	}
	got, _, _ = s.GetLoadBalancer(ctx, token, "lb-1")
	if got.Name != "renamed-lb" {
		t.Fatalf("expected renamed-lb, got %q", got.Name)
	}

	if err := s.RenameLoadBalancer(ctx, token, "nonexistent", "x"); err == nil {
		t.Fatal("expected error for missing LB")
	}

	if err := s.DeleteLoadBalancer(ctx, token, "lb-1"); err != nil {
		t.Fatal(err)
	}
	_, ok3, _ := s.GetLoadBalancer(ctx, token, "lb-1")
	if ok3 {
		t.Fatal("expected deleted")
	}
}

func TestMemoryStore_BucketLifecycle(t *testing.T) {
	s := store.NewMemoryStore()
	ctx := context.Background()
	token := "tok1"

	list, err := s.ListBuckets(ctx, token)
	if err != nil || len(list) != 0 {
		t.Fatalf("expected empty: %v %v", list, err)
	}

	bkt := newBucket("bkt-1", "my-bucket")
	if err := s.CreateBucket(ctx, token, bkt); err != nil {
		t.Fatal(err)
	}

	got, ok, err := s.GetBucket(ctx, token, "bkt-1")
	if err != nil || !ok || got.Region != "us-east-1" {
		t.Fatalf("unexpected get: %v %v %v", got, ok, err)
	}

	_, ok2, _ := s.GetBucket(ctx, "other", "bkt-1")
	if ok2 {
		t.Fatal("token isolation failed")
	}

	list, _ = s.ListBuckets(ctx, token)
	if len(list) != 1 {
		t.Fatalf("expected 1, got %d", len(list))
	}

	if err := s.RenameBucket(ctx, token, "bkt-1", "renamed-bucket"); err != nil {
		t.Fatal(err)
	}
	got, _, _ = s.GetBucket(ctx, token, "bkt-1")
	if got.Name != "renamed-bucket" {
		t.Fatalf("expected renamed-bucket, got %q", got.Name)
	}

	if err := s.RenameBucket(ctx, token, "nonexistent", "x"); err == nil {
		t.Fatal("expected error for missing bucket")
	}

	if err := s.DeleteBucket(ctx, token, "bkt-1"); err != nil {
		t.Fatal(err)
	}
	_, ok3, _ := s.GetBucket(ctx, token, "bkt-1")
	if ok3 {
		t.Fatal("expected deleted")
	}
}

func TestMemoryStore_DatabaseLifecycle(t *testing.T) {
	s := store.NewMemoryStore()
	ctx := context.Background()
	token := "tok1"

	list, err := s.ListDatabases(ctx, token)
	if err != nil || len(list) != 0 {
		t.Fatalf("expected empty: %v %v", list, err)
	}

	db := newDB("db-1", "my-db")
	if err := s.CreateDatabase(ctx, token, db); err != nil {
		t.Fatal(err)
	}

	got, ok, err := s.GetDatabase(ctx, token, "db-1")
	if err != nil || !ok || got.Engine != "postgres" {
		t.Fatalf("unexpected get: %v %v %v", got, ok, err)
	}
	if len(got.SubnetIDs) != 1 || got.SubnetIDs[0] != "sub-1" {
		t.Fatalf("unexpected subnet_ids: %v", got.SubnetIDs)
	}

	_, ok2, _ := s.GetDatabase(ctx, "other", "db-1")
	if ok2 {
		t.Fatal("token isolation failed")
	}

	list, _ = s.ListDatabases(ctx, token)
	if len(list) != 1 {
		t.Fatalf("expected 1, got %d", len(list))
	}

	if err := s.RenameDatabase(ctx, token, "db-1", "renamed-db"); err != nil {
		t.Fatal(err)
	}
	got, _, _ = s.GetDatabase(ctx, token, "db-1")
	if got.Name != "renamed-db" {
		t.Fatalf("expected renamed-db, got %q", got.Name)
	}

	if err := s.RenameDatabase(ctx, token, "nonexistent", "x"); err == nil {
		t.Fatal("expected error for missing database")
	}

	if err := s.DeleteDatabase(ctx, token, "db-1"); err != nil {
		t.Fatal(err)
	}
	_, ok3, _ := s.GetDatabase(ctx, token, "db-1")
	if ok3 {
		t.Fatal("expected deleted")
	}
}

func TestMemoryStore_KubernetesClusterLifecycle(t *testing.T) {
	s := store.NewMemoryStore()
	ctx := context.Background()
	token := "tok1"

	list, err := s.ListKubernetesClusters(ctx, token)
	if err != nil || len(list) != 0 {
		t.Fatalf("expected empty: %v %v", list, err)
	}

	cluster := newCluster("k8s-1", "my-cluster")
	if err := s.CreateKubernetesCluster(ctx, token, cluster); err != nil {
		t.Fatal(err)
	}

	got, ok, err := s.GetKubernetesCluster(ctx, token, "k8s-1")
	if err != nil || !ok || got.Version != "1.30" {
		t.Fatalf("unexpected get: %v %v %v", got, ok, err)
	}
	if len(got.SubnetIDs) != 1 || got.SubnetIDs[0] != "sub-1" {
		t.Fatalf("unexpected subnet_ids: %v", got.SubnetIDs)
	}

	_, ok2, _ := s.GetKubernetesCluster(ctx, "other", "k8s-1")
	if ok2 {
		t.Fatal("token isolation failed")
	}

	list, _ = s.ListKubernetesClusters(ctx, token)
	if len(list) != 1 {
		t.Fatalf("expected 1, got %d", len(list))
	}

	if err := s.RenameKubernetesCluster(ctx, token, "k8s-1", "renamed-cluster"); err != nil {
		t.Fatal(err)
	}
	got, _, _ = s.GetKubernetesCluster(ctx, token, "k8s-1")
	if got.Name != "renamed-cluster" {
		t.Fatalf("expected renamed-cluster, got %q", got.Name)
	}

	if err := s.RenameKubernetesCluster(ctx, token, "nonexistent", "x"); err == nil {
		t.Fatal("expected error for missing cluster")
	}

	if err := s.DeleteKubernetesCluster(ctx, token, "k8s-1"); err != nil {
		t.Fatal(err)
	}
	_, ok3, _ := s.GetKubernetesCluster(ctx, token, "k8s-1")
	if ok3 {
		t.Fatal("expected deleted")
	}
}

func TestMemoryStore_VPCRename(t *testing.T) {
	s := store.NewMemoryStore()
	ctx := context.Background()
	token := "tok1"

	vpc := model.VPC{ID: "vpc-1", Name: "original", Status: "available", CreatedAt: time.Now()}
	s.CreateVPC(ctx, token, vpc)

	if err := s.RenameVPC(ctx, token, "vpc-1", "renamed"); err != nil {
		t.Fatal(err)
	}
	got, _, _ := s.GetVPC(ctx, token, "vpc-1")
	if got.Name != "renamed" {
		t.Fatalf("expected renamed, got %q", got.Name)
	}

	if err := s.RenameVPC(ctx, token, "nonexistent", "x"); err == nil {
		t.Fatal("expected error")
	}
}

func TestMemoryStore_SubnetRename(t *testing.T) {
	s := store.NewMemoryStore()
	ctx := context.Background()
	token := "tok1"

	sub := model.Subnet{ID: "sub-1", Name: "original", Status: "available", CreatedAt: time.Now()}
	s.CreateSubnet(ctx, token, sub)

	if err := s.RenameSubnet(ctx, token, "sub-1", "renamed"); err != nil {
		t.Fatal(err)
	}
	got, _, _ := s.GetSubnet(ctx, token, "sub-1")
	if got.Name != "renamed" {
		t.Fatalf("expected renamed, got %q", got.Name)
	}

	if err := s.RenameSubnet(ctx, token, "nonexistent", "x"); err == nil {
		t.Fatal("expected error")
	}
}

func TestMemoryStore_VSIRename(t *testing.T) {
	s := store.NewMemoryStore()
	ctx := context.Background()
	token := "tok1"

	vsi := model.VSI{ID: "vsi-1", Name: "original", Status: "running", CreatedAt: time.Now()}
	s.CreateVSI(ctx, token, vsi)

	if err := s.RenameVSI(ctx, token, "vsi-1", "renamed"); err != nil {
		t.Fatal(err)
	}
	got, _, _ := s.GetVSI(ctx, token, "vsi-1")
	if got.Name != "renamed" {
		t.Fatalf("expected renamed, got %q", got.Name)
	}

	if err := s.RenameVSI(ctx, token, "nonexistent", "x"); err == nil {
		t.Fatal("expected error")
	}
}

// ---- JSONFileStore ----

func TestJSONFileStore_LoadBalancerLifecycle(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	ctx := context.Background()
	token := "tok1"

	s, _ := store.NewJSONFileStore(path)

	lb := newLB("lb-1", "my-lb")
	if err := s.CreateLoadBalancer(ctx, token, lb); err != nil {
		t.Fatal(err)
	}

	got, ok, err := s.GetLoadBalancer(ctx, token, "lb-1")
	if err != nil || !ok || got.Protocol != "http" {
		t.Fatalf("unexpected: %v %v %v", got, ok, err)
	}

	s.CreateLoadBalancer(ctx, token, newLB("lb-2", "lb2"))
	list, _ := s.ListLoadBalancers(ctx, token)
	if len(list) != 2 {
		t.Fatalf("expected 2, got %d", len(list))
	}

	if err := s.RenameLoadBalancer(ctx, token, "lb-1", "renamed-lb"); err != nil {
		t.Fatal(err)
	}

	if err := s.RenameLoadBalancer(ctx, token, "nonexistent", "x"); err == nil {
		t.Fatal("expected error")
	}

	if err := s.DeleteLoadBalancer(ctx, token, "lb-1"); err != nil {
		t.Fatal(err)
	}
	list, _ = s.ListLoadBalancers(ctx, token)
	if len(list) != 1 {
		t.Fatalf("expected 1 after delete, got %d", len(list))
	}
}

func TestJSONFileStore_LoadBalancerPersistence(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	ctx := context.Background()
	token := "tok1"

	s1, _ := store.NewJSONFileStore(path)
	s1.CreateLoadBalancer(ctx, token, newLB("lb-1", "persisted"))

	s2, _ := store.NewJSONFileStore(path)
	got, ok, _ := s2.GetLoadBalancer(ctx, token, "lb-1")
	if !ok || got.Name != "persisted" {
		t.Fatalf("LB not persisted: ok=%v name=%q", ok, got.Name)
	}
}

func TestJSONFileStore_BucketLifecycle(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	ctx := context.Background()
	token := "tok1"

	s, _ := store.NewJSONFileStore(path)

	bkt := newBucket("bkt-1", "my-bucket")
	if err := s.CreateBucket(ctx, token, bkt); err != nil {
		t.Fatal(err)
	}

	got, ok, err := s.GetBucket(ctx, token, "bkt-1")
	if err != nil || !ok || got.Region != "us-east-1" {
		t.Fatalf("unexpected: %v %v %v", got, ok, err)
	}

	s.CreateBucket(ctx, token, newBucket("bkt-2", "bkt2"))
	list, _ := s.ListBuckets(ctx, token)
	if len(list) != 2 {
		t.Fatalf("expected 2, got %d", len(list))
	}

	if err := s.RenameBucket(ctx, token, "bkt-1", "renamed-bucket"); err != nil {
		t.Fatal(err)
	}

	if err := s.RenameBucket(ctx, token, "nonexistent", "x"); err == nil {
		t.Fatal("expected error")
	}

	if err := s.DeleteBucket(ctx, token, "bkt-1"); err != nil {
		t.Fatal(err)
	}
	list, _ = s.ListBuckets(ctx, token)
	if len(list) != 1 {
		t.Fatalf("expected 1 after delete, got %d", len(list))
	}
}

func TestJSONFileStore_BucketPersistence(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	ctx := context.Background()
	token := "tok1"

	s1, _ := store.NewJSONFileStore(path)
	s1.CreateBucket(ctx, token, newBucket("bkt-1", "persisted"))

	s2, _ := store.NewJSONFileStore(path)
	got, ok, _ := s2.GetBucket(ctx, token, "bkt-1")
	if !ok || got.Name != "persisted" {
		t.Fatalf("bucket not persisted: ok=%v name=%q", ok, got.Name)
	}
}

func TestJSONFileStore_DatabaseLifecycle(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	ctx := context.Background()
	token := "tok1"

	s, _ := store.NewJSONFileStore(path)

	db := newDB("db-1", "my-db")
	if err := s.CreateDatabase(ctx, token, db); err != nil {
		t.Fatal(err)
	}

	got, ok, err := s.GetDatabase(ctx, token, "db-1")
	if err != nil || !ok || got.Engine != "postgres" {
		t.Fatalf("unexpected: %v %v %v", got, ok, err)
	}
	if len(got.SubnetIDs) != 1 {
		t.Fatalf("expected 1 subnet_id, got %v", got.SubnetIDs)
	}

	s.CreateDatabase(ctx, token, newDB("db-2", "db2"))
	list, _ := s.ListDatabases(ctx, token)
	if len(list) != 2 {
		t.Fatalf("expected 2, got %d", len(list))
	}

	if err := s.RenameDatabase(ctx, token, "db-1", "renamed-db"); err != nil {
		t.Fatal(err)
	}

	if err := s.RenameDatabase(ctx, token, "nonexistent", "x"); err == nil {
		t.Fatal("expected error")
	}

	if err := s.DeleteDatabase(ctx, token, "db-1"); err != nil {
		t.Fatal(err)
	}
	list, _ = s.ListDatabases(ctx, token)
	if len(list) != 1 {
		t.Fatalf("expected 1 after delete, got %d", len(list))
	}
}

func TestJSONFileStore_DatabasePersistence(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	ctx := context.Background()
	token := "tok1"

	s1, _ := store.NewJSONFileStore(path)
	s1.CreateDatabase(ctx, token, newDB("db-1", "persisted"))

	s2, _ := store.NewJSONFileStore(path)
	got, ok, _ := s2.GetDatabase(ctx, token, "db-1")
	if !ok || got.Name != "persisted" {
		t.Fatalf("database not persisted: ok=%v name=%q", ok, got.Name)
	}
}

func TestJSONFileStore_KubernetesClusterLifecycle(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	ctx := context.Background()
	token := "tok1"

	s, _ := store.NewJSONFileStore(path)

	cluster := newCluster("k8s-1", "my-cluster")
	if err := s.CreateKubernetesCluster(ctx, token, cluster); err != nil {
		t.Fatal(err)
	}

	got, ok, err := s.GetKubernetesCluster(ctx, token, "k8s-1")
	if err != nil || !ok || got.Version != "1.30" {
		t.Fatalf("unexpected: %v %v %v", got, ok, err)
	}
	if len(got.SubnetIDs) != 1 {
		t.Fatalf("expected 1 subnet_id, got %v", got.SubnetIDs)
	}

	s.CreateKubernetesCluster(ctx, token, newCluster("k8s-2", "cluster2"))
	list, _ := s.ListKubernetesClusters(ctx, token)
	if len(list) != 2 {
		t.Fatalf("expected 2, got %d", len(list))
	}

	if err := s.RenameKubernetesCluster(ctx, token, "k8s-1", "renamed-cluster"); err != nil {
		t.Fatal(err)
	}

	if err := s.RenameKubernetesCluster(ctx, token, "nonexistent", "x"); err == nil {
		t.Fatal("expected error")
	}

	if err := s.DeleteKubernetesCluster(ctx, token, "k8s-1"); err != nil {
		t.Fatal(err)
	}
	list, _ = s.ListKubernetesClusters(ctx, token)
	if len(list) != 1 {
		t.Fatalf("expected 1 after delete, got %d", len(list))
	}
}

func TestJSONFileStore_KubernetesClusterPersistence(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	ctx := context.Background()
	token := "tok1"

	s1, _ := store.NewJSONFileStore(path)
	s1.CreateKubernetesCluster(ctx, token, newCluster("k8s-1", "persisted"))

	s2, _ := store.NewJSONFileStore(path)
	got, ok, _ := s2.GetKubernetesCluster(ctx, token, "k8s-1")
	if !ok || got.Name != "persisted" {
		t.Fatalf("cluster not persisted: ok=%v name=%q", ok, got.Name)
	}
}

func TestJSONFileStore_VPCRename(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	ctx := context.Background()
	token := "tok1"

	s, _ := store.NewJSONFileStore(path)
	s.CreateVPC(ctx, token, model.VPC{ID: "vpc-1", Name: "original", Status: "available", CreatedAt: time.Now().UTC()})

	if err := s.RenameVPC(ctx, token, "vpc-1", "renamed"); err != nil {
		t.Fatal(err)
	}
	got, _, _ := s.GetVPC(ctx, token, "vpc-1")
	if got.Name != "renamed" {
		t.Fatalf("expected renamed, got %q", got.Name)
	}

	if err := s.RenameVPC(ctx, token, "nonexistent", "x"); err == nil {
		t.Fatal("expected error")
	}
}

func TestJSONFileStore_SubnetRename(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	ctx := context.Background()
	token := "tok1"

	s, _ := store.NewJSONFileStore(path)
	s.CreateSubnet(ctx, token, model.Subnet{ID: "sub-1", Name: "original", Status: "available", CreatedAt: time.Now().UTC()})

	if err := s.RenameSubnet(ctx, token, "sub-1", "renamed"); err != nil {
		t.Fatal(err)
	}
	got, _, _ := s.GetSubnet(ctx, token, "sub-1")
	if got.Name != "renamed" {
		t.Fatalf("expected renamed, got %q", got.Name)
	}

	if err := s.RenameSubnet(ctx, token, "nonexistent", "x"); err == nil {
		t.Fatal("expected error")
	}
}

func TestJSONFileStore_VSIRename(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state.json")
	ctx := context.Background()
	token := "tok1"

	s, _ := store.NewJSONFileStore(path)
	s.CreateVSI(ctx, token, model.VSI{ID: "vsi-1", Name: "original", Status: "running", CreatedAt: time.Now().UTC()})

	if err := s.RenameVSI(ctx, token, "vsi-1", "renamed"); err != nil {
		t.Fatal(err)
	}
	got, _, _ := s.GetVSI(ctx, token, "vsi-1")
	if got.Name != "renamed" {
		t.Fatalf("expected renamed, got %q", got.Name)
	}

	if err := s.RenameVSI(ctx, token, "nonexistent", "x"); err == nil {
		t.Fatal("expected error")
	}
}
