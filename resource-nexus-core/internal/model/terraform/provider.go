package tfModel

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/tbauriedel/resource-nexus-core/internal/utils/fileutils"
)

// Provider represents a Terraform provider
//
// Options are custom options for the provider as JSON.
type Provider struct {
	RequiredTerraformVersion string `json:"requiredVersion"` // "v1.1.0"
	ProviderName             string `json:"providerName"`    // "proxmox"
	Source                   string `json:"source"`          // "Telmate/proxmox"
	Version                  string `json:"version"`         // "3.0.2-rc06"
	Options                  []byte `json:"options"`         // Provider settings. JSON content
}

// HasOptions checks if the provider has options.
func (p *Provider) HasOptions() bool {
	return p.Options != nil
}

// ValidateOptionsSyntax validates whether the options are valid JSON.
func (p *Provider) ValidateOptionsSyntax() error {
	if !p.HasOptions() {
		return nil
	}

	var js json.RawMessage

	err := json.Unmarshal(p.Options, &js)
	if err != nil {
		return fmt.Errorf("options are non valid json: %w", err)
	}

	return nil
}

// GetProviderConfig returns the Terraform provider configuration as JSON string.
func (p *Provider) GetProviderConfig() (string, error) {
	err := p.ValidateOptionsSyntax()
	if err != nil {
		return "", err
	}

	result := map[string]any{
		"terraform": map[string]any{
			"required_version": p.RequiredTerraformVersion, // Required Terraform version
			"required_providers": map[string]any{
				p.ProviderName: map[string]any{ // Provider block
					"source":  p.Source,
					"version": p.Version,
				},
			},
		},
	}

	// Add provider options if defined
	if p.HasOptions() {
		result["provider"] = map[string]any{ // Provider settings block
			p.ProviderName: json.RawMessage(p.Options), // custom options for provider
		}
	}

	// Marshal result to JSON
	data, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("failed to marshal provider config: %w", err)
	}

	return string(data), nil
}

// WriteToFile writes the Terraform provider configuration to a file.
// workdir is the terraform working directory for the resource that will be managed.
func (p *Provider) WriteToFile(workdir string) error {
	filename := filepath.Join(workdir, "provider.tf.json")

	f, err := fileutils.OpenFile(filename)
	if err != nil {
		return fmt.Errorf("cant write provider file: %w", err)
	}
	defer f.Close()

	conf, err := p.GetProviderConfig()
	if err != nil {
		return fmt.Errorf("cant write provider file: %w", err)
	}

	_, err = f.WriteString(conf)
	if err != nil {
		return fmt.Errorf("cant write provider file: %w", err)
	}

	return nil
}
