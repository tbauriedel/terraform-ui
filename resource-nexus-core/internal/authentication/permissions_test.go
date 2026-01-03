package authentication

import (
	"testing"
)

func TestPermissions(t *testing.T) {
	// get all permissions and check return type
	var p map[string]string
	p = permissions()

	if len(p) == 0 {
		t.Fatal("permissions map is empty")
	}

	if permissions()["/system/health"] != "system:health:get" {
		t.Fatal("wrong permission returned")
	}
}

func TestGetPermissionForPath(t *testing.T) {
	p, ok := GetPermissionForPath("/system/health")
	if !ok {
		t.Fatal("permission not found")
	}

	if p != "system:health:get" {
		t.Fatal("wrong permission returned")
	}

	p, ok = GetPermissionForPath("/system/health2")
	if ok {
		t.Fatal("permission found for non existing path")
	}
}

func TestBuildPermissionString(t *testing.T) {
	actual := BuildPermissionString("security", "user", "create")
	expected := "security:user:create"

	if actual != expected {
		t.Fatalf("Expected: %s, Actual: %s", expected, actual)
	}
}
