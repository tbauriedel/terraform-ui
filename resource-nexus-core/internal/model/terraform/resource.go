package tfModel

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/tbauriedel/resource-nexus-core/internal/utils/fileutils"
)

type Resource struct {
	ResourceType string
	Name         string
	Options      []byte
}

// HasOptions checks if the provider has options.
func (r *Resource) HasOptions() bool {
	return r.Options != nil
}

// ValidateOptionsSyntax validates whether the options are valid JSON.
func (r *Resource) ValidateOptionsSyntax() error {
	if !r.HasOptions() {
		return nil
	}

	var js json.RawMessage

	err := json.Unmarshal(r.Options, &js)
	if err != nil {
		return fmt.Errorf("options are non valid json: %w", err)
	}

	return nil
}

// GetResourceConfig returns the Terraform resource configuration as JSON string.
func (r *Resource) GetResourceConfig() (string, error) {
	err := r.ValidateOptionsSyntax()
	if err != nil {
		return "", err
	}

	result := map[string]any{
		"resource": map[string]any{
			r.ResourceType: map[string]any{
				r.Name: map[string]any{},
			},
		},
	}

	// Add resource options if defined
	if r.HasOptions() {
		result["resource"].(map[string]any)[r.ResourceType].(map[string]any)[r.Name] = json.RawMessage(r.Options)
	}

	// Marshal result to JSON
	data, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("failed to marshal resource config: %w", err)
	}

	return string(data), nil
}

// WriteToFile writes the Terraform resource configuration to a file.
// workdir is the terraform working directory for the resource that will be managed.
func (r *Resource) WriteToFile(workdir string) error {
	filename := filepath.Join(workdir, "resource.tf.json")

	f, err := fileutils.OpenFile(filename)
	if err != nil {
		return fmt.Errorf("cant write resource file: %w", err)
	}

	defer f.Close()

	conf, err := r.GetResourceConfig()
	if err != nil {
		return fmt.Errorf("cant write resource file: %w", err)
	}

	_, err = f.WriteString(conf)
	if err != nil {
		return fmt.Errorf("cant write resource file: %w", err)
	}

	return nil
}
