package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config holds configuration parameters for the APT exporter.
type Config struct {
	CheckIntervalSeconds  int    `yaml:"check_interval_seconds"`
	ListenAddress         string `yaml:"listen_address"`          // e.g. ":9100"
	AptCheckPath          string `yaml:"apt_check_path"`          // e.g. "/usr/lib/update-notifier/apt-check"
	UpdateStampPath       string `yaml:"update_stamp_path"`       // e.g. "/var/lib/apt/periodic/update-success-stamp"
	RebootRequiredFile    string `yaml:"reboot_required_file"`    // e.g. "/var/run/reboot-required"
	LogLevel              string `yaml:"log_level"`               // e.g. "info", "debug"
	CommandTimeoutSeconds int    `yaml:"command_timeout_seconds"` // e.g. 10
	MetricsEndpoint       string `yaml:"metrics_endpoint"`        // e.g. "/metrics"
	MetricPrefix          string `yaml:"metric_prefix"`           // e.g. "ubuntu"
}

// Load reads a YAML configuration file and returns a Config struct.
func Load(path string) (*Config, error) {
	// Resolve absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve config path: %w", err)
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var conf Config
	if err := yaml.Unmarshal(data, &conf); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Validate configuration
	if err := conf.validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &conf, nil
}

// validate checks if the configuration is valid.
func (c *Config) validate() error {
	// Validate numeric values
	if c.CheckIntervalSeconds <= 0 {
		return fmt.Errorf("check_interval_seconds must be positive")
	}
	if c.CommandTimeoutSeconds <= 0 {
		return fmt.Errorf("command_timeout_seconds must be positive")
	}

	// Validate string values
	if c.ListenAddress == "" {
		return fmt.Errorf("listen_address cannot be empty")
	}
	if c.MetricsEndpoint == "" {
		return fmt.Errorf("metrics_endpoint cannot be empty")
	}
	if c.MetricPrefix == "" {
		return fmt.Errorf("metric_prefix cannot be empty")
	}

	// Validate log level
	switch c.LogLevel {
	case "debug", "info", "warn", "error":
		// Valid log levels
	default:
		return fmt.Errorf("invalid log_level: %s (must be one of: debug, info, warn, error)", c.LogLevel)
	}

	// Ensure metrics endpoint starts with a slash
	if c.MetricsEndpoint[0] != '/' {
		c.MetricsEndpoint = "/" + c.MetricsEndpoint
	}

	return nil
}

// ValidateFilePaths checks if the file paths in the configuration exist.
// This is separate from validate() because we may want to skip this check in tests.
func (c *Config) ValidateFilePaths() error {
	// Check if apt-check exists
	if _, err := os.Stat(c.AptCheckPath); err != nil {
		return fmt.Errorf("apt_check_path %s is not accessible: %w", c.AptCheckPath, err)
	}

	// Check if update stamp directory exists (the file itself may not exist yet)
	updateStampDir := filepath.Dir(c.UpdateStampPath)
	if _, err := os.Stat(updateStampDir); err != nil {
		return fmt.Errorf("update_stamp_path directory %s is not accessible: %w", updateStampDir, err)
	}

	// Check if reboot required directory exists (the file itself may not exist yet)
	rebootRequiredDir := filepath.Dir(c.RebootRequiredFile)
	if _, err := os.Stat(rebootRequiredDir); err != nil {
		return fmt.Errorf("reboot_required_file directory %s is not accessible: %w", rebootRequiredDir, err)
	}

	return nil
}
