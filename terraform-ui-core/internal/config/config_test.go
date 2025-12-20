package config

import "testing"

func TestLoadDefaults(t *testing.T) {
	_ = LoadDefaults()
	// only returns config struct, nothing to test
}

func TestLoadFromFile(t *testing.T) {
	c, err := LoadFromJSONFile("../../test/testdata/config/config_test.json")
	if err != nil {
		t.Fatal(err)
	}

	if c.Logging.Level != "info" {
		t.Fatalf("Actual: %s, Expected: %s", c.Logging.Level, "info")
	}
}

func TestGetConfigRedacted(t *testing.T) {
	c := LoadDefaults()
	c.Database.User = "foo"
	c.Database.Password = "bar"

	sanitized := c.GetConfigRedacted()

	if sanitized.Database.Password != RedactionPlaceholder || sanitized.Database.User != RedactionPlaceholder {
		t.Fatal("password and user should be redacted")
	}
}
