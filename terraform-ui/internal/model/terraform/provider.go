package model

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Provider represents a Terraform provider
//
// Options are custom options for the provider as JSON
type Provider struct {
	RequiredTerraformVersion string // "v1.1.0"
	ProviderName             string // "proxmox"
	Source                   string // "Telmate/proxmox"
	Version                  string // "3.0.2-rc06"
	Options                  []byte // JSON content. Needs to be validated!
}

// HasOptions checks if the provider has options
func (p *Provider) HasOptions() bool {
	if p.Options == nil {
		return false
	}
	return true
}

// ValidateOptions validates whether the options are valid JSON
func (p *Provider) ValidateOptions() error {
	if !p.HasOptions() {
		return nil
	}

	var js json.RawMessage
	err := json.Unmarshal(p.Options, &js)
	if err != nil {
		return errors.New(fmt.Sprint("options are non valid json", "err", err))
	}

	return nil
}

// GetProviderConfig returns the Terraform provider configuration as JSON object
func (p *Provider) GetProviderConfig() (string, error) {
	if err := p.ValidateOptions(); err != nil {
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
		return "", errors.New(fmt.Sprint("failed to marshal provider config", "err", err))
	}

	return string(data), nil
}
