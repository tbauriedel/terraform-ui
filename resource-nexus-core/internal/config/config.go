package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"golang.org/x/time/rate"
)

// Config represents the configuration of resource-nexus-core.
type Config struct {
	Logging  Logger   `json:"logging"`
	Database Database `json:"database"`
	Listener Listener `json:"listener"`
}

// Listener represents the listener configuration.
type Listener struct {
	ListenAddr          string        `json:"listenAddr"`
	ReadTimeout         time.Duration `json:"readTimeout"`
	IdleTimeout         time.Duration `json:"idleTimeout"`
	TLSEnabled          bool          `json:"tlsEnabled"`
	TLSSkipVerify       bool          `json:"tlsSkipVerify"`
	TLSCertPath         string        `json:"tlsCertFile"`
	TLSKeyPath          string        `json:"tlsKeyFile"`
	RateLimitGeneration rate.Limit    `json:"rateLimitGeneration"`
	RateLimitBucketSize int           `json:"rateLimitBucketSize"`
}

// Logger represents the logging configuration.
type Logger struct {
	Type  string `json:"type"`
	File  string `json:"file"`
	Level string `json:"level"`
}

// Database represents the database configuration.
type Database struct {
	Address       string `json:"address"`
	Port          int    `json:"port"`
	User          string `json:"user"`
	Password      string `json:"password"`
	Name          string `json:"name"`
	TLSSkipVerify bool   `json:"tlsSkipVerify"`
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
		Database: Database{},
		Listener: Listener{
			ListenAddr:          ":4890",
			ReadTimeout:         10 * time.Second,
			IdleTimeout:         120 * time.Second,
			TLSSkipVerify:       false,
			RateLimitGeneration: 5,
			RateLimitBucketSize: 25,
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
