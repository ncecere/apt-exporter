package config

import (
	"fmt"
	"os"

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
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var conf Config
	if err := yaml.Unmarshal(data, &conf); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Validate configuration
	if err := conf.validate(); err != nil {
		return nil, err
	}

	return &conf, nil
}

// validate checks if the configuration is valid.
func (c *Config) validate() error {
	if c.CheckIntervalSeconds <= 0 {
		return fmt.Errorf("check_interval_seconds must be positive")
	}
	if c.ListenAddress == "" {
		return fmt.Errorf("listen_address cannot be empty")
	}
	if c.CommandTimeoutSeconds <= 0 {
		return fmt.Errorf("command_timeout_seconds must be positive")
	}
	if c.MetricsEndpoint == "" {
		return fmt.Errorf("metrics_endpoint cannot be empty")
	}
	if c.MetricPrefix == "" {
		return fmt.Errorf("metric_prefix cannot be empty")
	}
	return nil
}
