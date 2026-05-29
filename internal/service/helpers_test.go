package service

import "testing"

func TestIsStateLevelTenant(t *testing.T) {
	if !isStateLevelTenant("pb") {
		t.Fatal("expected state-level tenant")
	}
	if isStateLevelTenant("pb.amritsar") {
		t.Fatal("expected ulb-level tenant")
	}
}
