package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Config struct {
	Logging  Logger   `json:"logging"`
	Database Database `json:"database"`
	Listener Listener `json:"listener"`
}

type Listener struct {
	ListenAddr    string        `json:"listen_addr"`
	ReadTimeout   time.Duration `json:"read_timeout"`
	IdleTimeout   time.Duration `json:"idle_timeout"`
	TlsEnabled    bool          `json:"tls_enabled"`
	TlsSkipVerify bool          `json:"tls_skip_verify"`
	TlsCertPath   string        `json:"tls_cert_file"`
	TlsKeyPath    string        `json:"tls_key_file"`
}

type Logger struct {
	Type  string `json:"type"`
	File  string `json:"file"`
	Level string `json:"level"`
}

// conn_str=postgres://user:pass@localhost:5432/dbname?sslmode=disable
type Database struct {
	Address       string `json:"address"`
	Port          int    `json:"port"`
	User          string `json:"user"`
	Password      string `json:"password"`
	Name          string `json:"name"`
	TlsSkipVerify bool   `json:"tls_skip_verify"`
}

var (
	RedactionPlaceholder = "********"
)

// LoadDefaults returns a Config struct with default values
func LoadDefaults() Config {
	return Config{
		Logging: Logger{
			Type:  "stdout",
			File:  "terraform-ui-core.log",
			Level: "info",
		},
		Database: Database{},
		Listener: Listener{
			ListenAddr:    ":4890",
			ReadTimeout:   10 * time.Second,
			IdleTimeout:   120 * time.Second,
			TlsSkipVerify: false,
		},
	}
}

// LoadFromJSONFile reads the config file and return a Config struct
func LoadFromJSONFile(file string) (Config, error) {
	data, err := os.ReadFile(file)
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

// GetConfigRedacted returns a new Config object, with sensitive testdata masked with RedactionPlaceholder
//
// Sensitive testdata includes:
//   - Database.User
//   - Database.Password
func (c Config) GetConfigRedacted() Config {
	sanitized := c

	if c.Database.User != "" {
		sanitized.Database.User = RedactionPlaceholder
	}

	if c.Database.Password != "" {
		sanitized.Database.Password = RedactionPlaceholder
	}

	return sanitized
}
