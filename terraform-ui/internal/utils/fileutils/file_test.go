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
