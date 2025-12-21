package fileutils

import (
	"fmt"
	"os"
)

// OpenFile opens the provided file. If the file does not exist, it will be created.
// File permissions are set to 0644.
func OpenFile(file string) (*os.File, error) {
	f, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY, 0644) //nolint:gosec
	if err != nil {
		return nil, fmt.Errorf("failed to open / create file: %w", err)
	}

	return f, nil
}

// FileExists returns true if the file exists.
func FileExists(file string) bool {
	_, err := os.Stat(file)

	return !os.IsNotExist(err)
}
