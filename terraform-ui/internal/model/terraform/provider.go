package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/tbauriedel/terraform-ui/internal/utils/fileutils"
)

// Provider represents a Terraform provider
//
// Options are custom options for the provider as JSON
type Provider struct {
	RequiredTerraformVersion string `json:"required_version"` // "v1.1.0"
	ProviderName             string `json:"provider-name"`    // "proxmox"
	Source                   string `json:"source"`           // "Telmate/proxmox"
	Version                  string `json:"version"`          // "3.0.2-rc06"
	Options                  []byte `json:"options"`          // Provider settings. JSON content
}

// HasOptions checks if the provider has options
func (p *Provider) HasOptions() bool {
	if p.Options == nil {
		return false
	}
	return true
}

// ValidateOptionsSyntax validates whether the options are valid JSON
func (p *Provider) ValidateOptionsSyntax() error {
	if !p.HasOptions() {
		return nil
	}

	var js json.RawMessage
	err := json.Unmarshal(p.Options, &js)
	if err != nil {
		return errors.New(fmt.Sprint("options are non valid json", "error", err))
	}

	return nil
}

// GetProviderConfig returns the Terraform provider configuration as JSON string
func (p *Provider) GetProviderConfig() (string, error) {
	if err := p.ValidateOptionsSyntax(); err != nil {
		return "", err
	}

	result := map[string]interface{}{
		"terraform": map[string]interface{}{
			"required_version": p.RequiredTerraformVersion, // Required Terraform version
			"required_providers": map[string]interface{}{
				p.ProviderName: map[string]interface{}{ // Provider block
					"source":  p.Source,
					"version": p.Version,
				},
			},
		},
	}

	// Add provider options if defined
	if p.HasOptions() {
		result["provider"] = map[string]interface{}{ // Provider settings block
			p.ProviderName: json.RawMessage(p.Options), // custom options for provider
		}
	}

	// Marshal result to JSON
	data, err := json.Marshal(result)
	if err != nil {
		return "", errors.New(fmt.Sprint("failed to marshal provider config", "error", err))
	}

	return string(data), nil
}

// WriteToFile writes the Terraform provider configuration to a file.
// workdir is the terraform working directory for the resource that will be managed
func (p *Provider) WriteToFile(workdir string) error {
	filename := filepath.Join(workdir, "provider.tf.json")

	f, err := fileutils.OpenFile(filename)
	if err != nil {
		return errors.New(fmt.Sprint("cant write provider file", "error", err))
	}
	defer f.Close()

	conf, err := p.GetProviderConfig()
	if err != nil {
		return errors.New(fmt.Sprint("cant write provider file", "error", err))
	}

	_, err = f.Write([]byte(conf))
	if err != nil {
		return errors.New(fmt.Sprint("cant write provider file", "error", err))
	}

	return nil
}
