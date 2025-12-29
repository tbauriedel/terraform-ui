package tf

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"
)

func (p *TerraformProvider) TestHasOptions(t *testing.T) {
	// Test non-empty options
	p.Options = []byte("{'hello':'world'}")
	if !p.HasOptions() {
		t.Fatal("options should be set")
	}

	// Test empty options
	p.Options = nil
	if p.HasOptions() {
		t.Fatal("options should not be set")
	}
}

func (p *TerraformProvider) TestValidateOptionsSyntax(t *testing.T) {
	// Test valid options
	p.Options = []byte(`{"hello":"world"}`)

	if err := p.ValidateOptionsSyntax(); err != nil {
		t.Fatal("options should be valid")
	}

	// Test invalid options
	p.Options = []byte("hello''world")
	if err := p.ValidateOptionsSyntax(); err == nil {
		t.Fatal("options should be invalid")
	}
}

func (p *TerraformProvider) TestGetProviderConfig(t *testing.T) {
	p = &TerraformProvider{
		RequiredTerraformVersion: ">= 0.1.0",
		ProviderName:             "aws",
		Source:                   "hashicorp/aws",
		Version:                  "3.28.0",
		Options:                  []byte(`{"region":"eu-central-1"}`),
	}

	config, err := p.GetProviderConfig()
	if err != nil {
		t.Fatalf("failed to create provider config: %s", err.Error())
	}

	validConfig := `{"provider":{"aws":{"region":"eu-central-1"}},"terraform":{"required_providers":{"aws":{"source":"hashicorp/aws","version":"3.28.0"}},"required_version":">= 0.1.0"}}`

	var actual, expected interface{}

	json.Unmarshal([]byte(config), &actual)
	json.Unmarshal([]byte(validConfig), &expected)

	if !reflect.DeepEqual(actual, expected) {
		t.Fatal("configs do not match")
	}
}

func (p *TerraformProvider) TestWriteToFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("../../../test/testdata", "tmp")
	if err != nil {
		t.Fatal(err)
	}

	defer os.RemoveAll(tmpDir)

	p = &TerraformProvider{
		RequiredTerraformVersion: ">= 0.1.0",
		ProviderName:             "aws",
		Source:                   "hashicorp/aws",
		Version:                  "3.28.0",
		Options:                  []byte(`{"region":"eu-central-1"}`),
	}

	err = p.WriteToFile(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(tmpDir + "/provider.tf.json")
	if err != nil {
		t.Fatal(err)
	}

	validConfig := `{"provider":{"aws":{"region":"eu-central-1"}},"terraform":{"required_providers":{"aws":{"source":"hashicorp/aws","version":"3.28.0"}},"required_version":">= 0.1.0"}}`

	var actual, expected interface{}

	json.Unmarshal(content, &actual)
	json.Unmarshal([]byte(validConfig), &expected)

	if !reflect.DeepEqual(actual, expected) {
		t.Fatal("configs do not match")
	}
}
