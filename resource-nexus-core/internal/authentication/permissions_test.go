package authentication

import "testing"

func TestBuildPermissionString(t *testing.T) {
	actual := BuildPermissionString("get", "user")
	expected := "get:user"

	if actual != expected {
		t.Fatalf("Expected: %s, Actual: %s", expected, actual)
	}
}
