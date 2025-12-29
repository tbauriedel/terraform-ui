package fileutils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestOpenFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("../../../test/testdata", "tmp")
	if err != nil {
		t.Fatal(err)
	}

	defer os.RemoveAll(tmpDir)

	f, err := OpenFile(filepath.Join(tmpDir, "test.txt"))
	if err != nil {
		t.Fatal(err)
	}

	f.Close()
}

func TestFileExists(t *testing.T) {
	// config_test.json exists in repo for sure. So we use it here
	if !FileExists("../../../test/testdata/config/config_test.json") {
		t.Fatal("file does not exist")
	}
}

func TestIsPath(t *testing.T) {
	invalid := "test.txt"
	if IsPath(invalid) {
		t.Fatal("path should not be valid")
	}

	valid := "/test.txt"
	if !IsPath(valid) {
		t.Fatal("path should be valid")
	}
}

func TestIsExecutableForUser(t *testing.T) {
	invalid := "../../../test/testdata/files/empty-non-executable"

	ok, err := IsExecutableForUser(invalid)
	if err != nil || ok {
		t.Fatal(err)
	}

	valid := "../../../test/testdata/files/empty-executable"

	ok, err = IsExecutableForUser(valid)
	if err != nil || !ok {
		t.Fatal(err)
	}
}
