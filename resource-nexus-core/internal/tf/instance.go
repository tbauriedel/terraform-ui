package tf

import (
	"fmt"
	"os"
	"time"

	"github.com/tbauriedel/resource-nexus-core/internal/common/fileutils"
)

// TerraformInstance represents a terraform instance.
type TerraformInstance struct {
	ExecutablePath    string
	BaseDir           string
	tmpWorkDir        string
	CommandTimeout    time.Duration
	ConfigCreated     bool
	WorkspacePrepared bool
}

// NewInstance creates a new terraform instance.
// Validates the provided executable and creates a temporary working directory.
func NewInstance(executable string, basedir string) (*TerraformInstance, error) {
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
func (tf *TerraformInstance) Cleanup() error {
	err := os.Remove(tf.tmpWorkDir)
	if err != nil {
		return fmt.Errorf("failed to cleanup temporary terraform working directory: %w", err)
	}

	return nil
}

// getDefaults returns the default values for the terraform instance.
func getDefaults() *TerraformInstance {
	return &TerraformInstance{
		ExecutablePath:    "/usr/local/bin/terraform",
		BaseDir:           "/tmp",
		ConfigCreated:     false,
		WorkspacePrepared: false,
		CommandTimeout:    10 * time.Minute,
	}
}

// prepare checks if the terraform executable is available and executable.
// Creates the temporary working directory and saves it into the TerraformInstance.
func (tf *TerraformInstance) prepare() error {
	// Check if the executable path is set
	if tf.ExecutablePath == "" {
		return fmt.Errorf("terraform executable path is empty")
	}

	// Check if the executable is a valid full path
	if !fileutils.IsPath(tf.ExecutablePath) {
		return fmt.Errorf(
			"terraform executable path '%s' is not valid. Please provide full path to executable",
			tf.ExecutablePath,
		)
	}

	// Check if the executable is available
	ok := fileutils.FileExists(tf.ExecutablePath)
	if !ok {
		return fmt.Errorf("terraform executable not found at '%s'", tf.ExecutablePath)
	}

	// Check if the executable is executable
	ok, err := fileutils.IsExecutableForUser(tf.ExecutablePath)
	if err != nil {
		return fmt.Errorf("failed to check terraform executable permissions: %w", err)
	}

	if !ok {
		return fmt.Errorf("terraform executable '%s' is not executable", tf.ExecutablePath)
	}

	// check base dir is not empty
	if tf.BaseDir == "" {
		return fmt.Errorf("terraform base directory is empty")
	}

	// Create working directory
	tmpDir, err := os.MkdirTemp(tf.BaseDir, "tmp-nexus-resource-core-")
	if err != nil {
		return fmt.Errorf("failed to create terraform working directory: %w", err)
	}

	tf.tmpWorkDir = tmpDir

	tf.WorkspacePrepared = true

	return nil
}
