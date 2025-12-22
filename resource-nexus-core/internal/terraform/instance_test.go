package terraform

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetDefaults(t *testing.T) {
	i := getDefaults()

	if i.ExecutablePath != "/usr/local/bin/terraform" {
		t.Fatal("Default executable path, does not match the expected value")
	}
}

func TestPrepare(t *testing.T) {
	i := Instance{
		ExecutablePath: "../../test/testdata/files/empty-executable",
		BaseDir:        "/tmp",
	}

	defer func() {
		path := filepath.Join(i.BaseDir, "tmp-nexus-resource-core-*")
		err := os.RemoveAll(path)
		if err != nil {
			t.Fatal(err)
		}
	}()

	err := i.prepare()
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewInstance(t *testing.T) {
	_, err := NewInstance("../../test/testdata/files/empty-executable", "/tmp")
	if err != nil {
		t.Fatal(err)
	}
}

func TestCleanup(t *testing.T) {
	i, err := NewInstance("../../test/testdata/files/empty-executable", "/tmp")
	if err != nil {

	}

	err = i.Cleanup()
	if err != nil {
		t.Fatal(err)
	}
}
