package terraform

import (
	"fmt"
	"os"

	"github.com/tbauriedel/resource-nexus-core/internal/utils/fileutils"
)

// Instance represents a terraform instance.
type Instance struct {
	ExecutablePath string
	BaseDir        string
	tmpWorkDir     string
}

// NewInstance creates a new terraform instance.
// Validates the provided executable and creates a temporary working directory.
func NewInstance(executable string, basedir string) (*Instance, error) {
	i := getDefaults()

	i.ExecutablePath = executable
	i.BaseDir = basedir

	// check executable and create tmp workdir
	err := i.prepare()
	if err != nil {
		return nil, fmt.Errorf("failed to prepare terraform instance: %w", err)
	}

	return i, nil
}

// Cleanup removes the temporary working directory.
func (i *Instance) Cleanup() error {
	err := os.Remove(i.tmpWorkDir)
	if err != nil {
		return fmt.Errorf("failed to cleanup temporary terraform working directory: %w", err)
	}

	return nil
}

// getDefaults returns the default values for the terraform instance.
func getDefaults() *Instance {
	return &Instance{
		ExecutablePath: "/usr/local/bin/terraform",
		BaseDir:        "/tmp",
	}
}

// prepare checks if the terraform executable is available and executable.
// Creates the temporary working directory and saves it into the Instance.
func (i *Instance) prepare() error {
	// Check if the executable path is set
	if i.ExecutablePath == "" {
		return fmt.Errorf("terraform executable path is empty")
	}

	// Check if the executable is a valid full path
	if !fileutils.IsPath(i.ExecutablePath) {
		return fmt.Errorf(
			"terraform executable path '%s' is not valid. Please provide full path to executable",
			i.ExecutablePath,
		)
	}

	// Check if the executable is available
	ok := fileutils.FileExists(i.ExecutablePath)
	if !ok {
		return fmt.Errorf("terraform executable not found at '%s'", i.ExecutablePath)
	}

	// Check if the executable is executable
	ok, err := fileutils.IsExecutableForUser(i.ExecutablePath)
	if err != nil {
		return fmt.Errorf("failed to check terraform executable permissions: %w", err)
	}

	if !ok {
		return fmt.Errorf("terraform executable '%s' is not executable", i.ExecutablePath)
	}

	// check base dir is not empty
	if i.BaseDir == "" {
		return fmt.Errorf("terraform base directory is empty")
	}

	// Create working directory
	tmpDir, err := os.MkdirTemp(i.BaseDir, "tmp-nexus-resource-core-")
	if err != nil {
		return fmt.Errorf("failed to create terraform working directory: %w", err)
	}

	i.tmpWorkDir = tmpDir

	return nil
}
