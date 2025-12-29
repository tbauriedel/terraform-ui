package tf

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"
)

func (r *TerraformResource) TestHasOptions(t *testing.T) {
	r.Options = []byte("{'foo':'bar}")
	if !r.HasOptions() {
		t.Fatal("options should be set")
	}

	r.Options = nil
	if r.HasOptions() {
		t.Fatal("options should not be set")
	}
}

func (r *TerraformResource) TestValidateOptionsSyntax(t *testing.T) {
	// Test valid options
	r.Options = []byte("{'foo':'bar}")

	if err := r.ValidateOptionsSyntax(); err == nil {
		t.Fatal("options should be invalid")
	}

	// Test invalid options
	r.Options = []byte("{'foo':a'bar}")
	if err := r.ValidateOptionsSyntax(); err == nil {
		t.Fatal("options should be invalid")
	}
}

func (r *TerraformResource) TestGetResourceConfig(t *testing.T) {
	r = &TerraformResource{
		ResourceType: "aws_instance",
		Name:         "test",
	}

	config, err := r.GetResourceConfig()
	if err != nil {
		t.Fatal(fmt.Sprint("failed to create resource config", "error", err))
	}

	validConfig := `{"resource":{"aws_instance":{"test":{}}}}`

	var actual, expected interface{}

	json.Unmarshal([]byte(config), &actual)
	json.Unmarshal([]byte(validConfig), &expected)

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("configs do not match.\nactual: %s\nexpected: %s\n", config, validConfig)
	}
}

func (r *TerraformResource) TestWriteToFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("../../../test/testdata", "tmp")
	if err != nil {
		t.Fatal(err)
	}

	defer os.RemoveAll(tmpDir)

	r = &TerraformResource{
		ResourceType: "aws_instance",
		Name:         "test",
	}

	err = r.WriteToFile(tmpDir)
	if err != nil {
		t.Fatal(err)
	}

	content, err := os.ReadFile(tmpDir + "/resource.tf.json")
	if err != nil {
		t.Fatal(err)
	}

	validConfig := `{"resource":{"aws_instance":{"test":{}}}}`

	var actual, expected interface{}

	json.Unmarshal(content, &actual)
	json.Unmarshal([]byte(validConfig), &expected)

	if !reflect.DeepEqual(actual, expected) {
		t.Fatal("configs do not match")
	}
}
