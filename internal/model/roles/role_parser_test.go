package roles_test

import (
	"testing"

	"github.com/Yash-Watchguard/Tasknest/internal/model/roles"
)


func TestRoleParser(t *testing.T) {
	tests := []struct {
		role     roles.Role
		expected string
	}{
		{roles.Role(0), "Admin"},
		{roles.Role(1), "Manager"},
		{roles.Role(2), "Employee"},
		{roles.Role(99), ""}, // invalid role
	}

	for _, tt := range tests {
		got := roles.RoleParser(tt.role)
		if got != tt.expected {
			t.Errorf("RoleParser(%d) = %q, want %q", tt.role, got, tt.expected)
		}
	}
}

func TestRoleScan(t *testing.T) {
	var r roles.Role

	// valid scan
	err := r.Scan(int64(1))
	if err != nil {
		t.Errorf("unexpected error scanning: %v", err)
	}
	if r != roles.Role(1) {
		t.Errorf("expected Role(1), got %v", r)
	}

	// invalid scan (wrong type)
	err = r.Scan("not-an-int")
	if err == nil {
		t.Errorf("expected error scanning string, got nil")
	}
}

func TestRoleValue(t *testing.T) {
	r := roles.Role(2)

	val, err := r.Value()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if v, ok := val.(int64); !ok || v != 2 {
		t.Errorf("expected int64(2), got %#v", val)
	}
}
