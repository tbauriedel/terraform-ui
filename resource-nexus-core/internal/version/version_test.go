package version

import "testing"

func TestGetVersion(t *testing.T) {
	if GetVersion() != VERSION {
		t.Fatal("Version does not match expected value")
	}
}
