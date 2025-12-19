package resource

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/tbauriedel/terraform-ui/internal/utils/fileutils"
)

type Resource struct {
	ResourceType string
	Name         string
	Options      []byte
}

// HasOptions checks if the provider has options
func (r *Resource) HasOptions() bool {
	if r.Options == nil {
		return false
	}
	return true
}

// ValidateOptionsSyntax validates whether the options are valid JSON
func (r *Resource) ValidateOptionsSyntax() error {
	if !r.HasOptions() {
		return nil
	}

	var js json.RawMessage
	err := json.Unmarshal(r.Options, &js)
	if err != nil {
		return errors.New(fmt.Sprint("options are non valid json", "error", err))
	}

	return nil
}

// GetResourceConfig returns the Terraform resource configuration as JSON string
func (r *Resource) GetResourceConfig() (string, error) {
	if err := r.ValidateOptionsSyntax(); err != nil {
		return "", err
	}

	result := map[string]interface{}{
		"resource": map[string]interface{}{
			r.ResourceType: map[string]interface{}{
				r.Name: map[string]interface{}{},
			},
		},
	}

	// Add resource options if defined
	if r.HasOptions() {
		result["resource"].(map[string]interface{})[r.ResourceType].(map[string]interface{})[r.Name] = json.RawMessage(r.Options)
	}

	// Marshal result to JSON
	data, err := json.Marshal(result)
	if err != nil {
		return "", errors.New(fmt.Sprint("failed to marshal resource config", "error", err))
	}

	return string(data), nil
}

// WriteToFile writes the Terraform resource configuration to a file.
// workdir is the terraform working directory for the resource that will be managed
func (r *Resource) WriteToFile(workdir string) error {
	filename := filepath.Join(workdir, "resource.tf.json")

	f, err := fileutils.OpenFile(filename)
	if err != nil {
		return errors.New(fmt.Sprint("cant write resource file", "error", err))
	}

	defer f.Close()

	conf, err := r.GetResourceConfig()
	if err != nil {
		return errors.New(fmt.Sprint("failed to create resource config", "error", err))
	}

	_, err = f.Write([]byte(conf))
	if err != nil {
		return errors.New(fmt.Sprint("cant write resource file", "error", err))
	}

	return nil
}
