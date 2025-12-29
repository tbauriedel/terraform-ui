package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Config represents the configuration of resource-nexus-core.
type Config struct {
	Logging  Logger   `json:"logging"`
	Database Database `json:"database"`
	Listener Listener `json:"listener"`
	Security Security `json:"security"`
}

var (
	redactionPlaceholder = "********" //nolint:gochecknoglobals
)

// LoadDefaults returns a Config struct with default values.
func LoadDefaults() Config {
	return Config{
		Logging: Logger{
			Type:  "stdout",
			File:  "resource-nexus-core.log",
			Level: "info",
		},
		Database: Database{
			Address:  "localhost",
			Port:     5432,
			User:     "",
			Password: "",
			Name:     "resource-nexus",
			TLSMode:  "verify-full",
		},
		Listener: Listener{
			ListenAddr:                 ":4890",
			ReadTimeout:                10 * time.Second,
			IdleTimeout:                120 * time.Second,
			TLSSkipVerify:              false,
			GlobalRateLimitGeneration:  5,
			GlobalRateLimitBucketSize:  25,
			IpBasedRateLimitGeneration: 2,
			IpBasedRateLimitBucketSize: 10,
		},
		Security: Security{
			PasswordHashing: HashingParams{
				Iterations:   3,
				MemoryCost:   64 * 1024, // 64MB
				ThreadsCount: 1,
				KeyLength:    32,
				SaltLength:   16,
			},
		},
	}
}

// LoadFromJSONFile reads the config file and return a Config struct.
//
// The function will try to read the provided file. Be sure, the file is not sensitive.
func LoadFromJSONFile(file string) (Config, error) {
	data, err := os.ReadFile(file) // #nosec G304 file path is provided by a trusted CLI flag
	if err != nil {
		return Config{}, fmt.Errorf("failed to read config file. %w", err)
	}

	c := LoadDefaults()

	err = json.Unmarshal(data, &c)
	if err != nil {
		return Config{}, fmt.Errorf("failed to parse config file. %w", err)
	}

	return c, nil
}

// GetConfigRedacted returns a new Config object, with sensitive testdata masked with redactionPlaceholder
//
// Sensitive testdata includes:
//   - Database.User
//   - Database.Password
func (c Config) GetConfigRedacted() Config {
	sanitized := c

	if c.Database.User != "" {
		sanitized.Database.User = redactionPlaceholder
	}

	if c.Database.Password != "" {
		sanitized.Database.Password = redactionPlaceholder
	}

	return sanitized
}
