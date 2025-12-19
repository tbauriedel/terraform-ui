package fileutils

import "testing"

func TestOpenFile(t *testing.T) {
	f, err := OpenFile("../../../test/testdata/test.txt")
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
