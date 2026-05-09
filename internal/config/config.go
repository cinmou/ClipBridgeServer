// SPDX-License-Identifier: GPL-3.0-only

package config

import (
	"fmt"
	"net"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Auth    AuthConfig    `yaml:"auth"`
	Server  ServerConfig  `yaml:"server"`
	Storage StorageConfig `yaml:"storage"`
	Limits  LimitsConfig  `yaml:"limits"`
	Cleaner CleanerConfig `yaml:"cleaner"`
}

// AuthConfig contains the shared bearer token used by API clients.
type AuthConfig struct {
	Token string `yaml:"token"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type StorageConfig struct {
	DataDir        string `yaml:"data_dir"`
	DatabasePath   string `yaml:"database_path"`
	TTLHours       int    `yaml:"ttl_hours"`
	MaxItems       int    `yaml:"max_items"`
	MaxTotalSizeMB int    `yaml:"max_total_size_mb"`
}

// LimitsConfig contains the request limits that are actively enforced by the
// current implementation. We only enable fields once code actually uses them so
// the config surface stays honest across phases.
type LimitsConfig struct {
	MinTextBytes    int `yaml:"min_text_bytes"`
	MaxTextBytes    int `yaml:"max_text_bytes"`
	MaxRequestBytes int `yaml:"max_request_bytes"`
}

// CleanerConfig controls whether the background retention worker runs and how
// often it wakes up to reload the persisted cleanup policy.
type CleanerConfig struct {
	Enabled         bool `yaml:"enabled"`
	IntervalMinutes int  `yaml:"interval_minutes"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config file %q: %w", path, err)
	}

	cfg := Default()
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parse config file %q: %w", path, err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("validate config file %q: %w", path, err)
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	if c.Server.Host == "" {
		return fmt.Errorf("server.host must not be empty")
	}

	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("server.port must be between 1 and 65535")
	}

	if ip := net.ParseIP(c.Server.Host); ip == nil && c.Server.Host != "localhost" {
		return fmt.Errorf("server.host must be a valid IP address or localhost")
	}

	if c.Auth.Token == "" {
		return fmt.Errorf("auth.token must not be empty")
	}

	if c.Storage.DataDir == "" {
		return fmt.Errorf("storage.data_dir must not be empty")
	}

	if c.Storage.DatabasePath == "" {
		return fmt.Errorf("storage.database_path must not be empty")
	}

	if c.Storage.TTLHours <= 0 {
		return fmt.Errorf("storage.ttl_hours must be greater than 0")
	}

	if c.Storage.MaxItems <= 0 {
		return fmt.Errorf("storage.max_items must be greater than 0")
	}

	if c.Storage.MaxTotalSizeMB <= 0 {
		return fmt.Errorf("storage.max_total_size_mb must be greater than 0")
	}

	if c.Limits.MinTextBytes < 0 {
		return fmt.Errorf("limits.min_text_bytes must be greater than or equal to 0")
	}

	if c.Limits.MaxTextBytes <= 0 {
		return fmt.Errorf("limits.max_text_bytes must be greater than 0")
	}

	if c.Limits.MaxRequestBytes <= 0 {
		return fmt.Errorf("limits.max_request_bytes must be greater than 0")
	}

	if c.Limits.MinTextBytes > c.Limits.MaxTextBytes {
		return fmt.Errorf("limits.min_text_bytes must not be greater than limits.max_text_bytes")
	}

	if c.Cleaner.IntervalMinutes <= 0 {
		return fmt.Errorf("cleaner.interval_minutes must be greater than 0")
	}

	return nil
}

func (c ServerConfig) Address() string {
	return net.JoinHostPort(c.Host, strconv.Itoa(c.Port))
}
