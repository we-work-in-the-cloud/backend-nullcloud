package uid_test

import (
	"strings"
	"testing"

	"github.com/we-work-in-the-cloud/nullcloud/backend/internal/uid"
)

func TestNew(t *testing.T) {
	id := uid.New("vpc")
	if !strings.HasPrefix(id, "vpc-") {
		t.Fatalf("expected vpc- prefix, got %q", id)
	}
	if len(id) != len("vpc-")+16 {
		t.Fatalf("unexpected length: %q", id)
	}

	a, b := uid.New("x"), uid.New("x")
	if a == b {
		t.Fatal("expected unique IDs")
	}
}
