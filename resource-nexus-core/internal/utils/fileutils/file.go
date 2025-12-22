package fileutils

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
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

// IsPath returns true if the string is a path.
func IsPath(s string) bool {
	return filepath.Base(s) != s
}

// IsExecutableForUser checks if the file is executable for the current user.
// Returns true if the file is executable, false otherwise.
//
// Tested for user, group and others.
func IsExecutableForUser(path string) (bool, error) {
	// get file info
	info, err := os.Stat(path)
	if err != nil {
		return false, fmt.Errorf("failed to get file info: %w", err)
	}

	// get file ownership and permissions
	fStat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return false, nil
	}

	// get file permissions
	fMode := info.Mode().Perm()

	// get uid and gid of current user
	uid := uint32(os.Getuid()) //nolint:gosec
	gid := uint32(os.Getgid()) //nolint:gosec

	// is executable for user
	if fStat.Uid == uid {
		return fMode&0100 != 0, nil
	}

	// is executable for group
	if fStat.Gid == gid {
		return fMode&0010 != 0, nil
	}

	// is executable for others
	return fMode&0001 != 0, nil
}
