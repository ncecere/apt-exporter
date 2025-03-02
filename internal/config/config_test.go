package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad(t *testing.T) {
	// Create a temporary test config file
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "test_config.yml")

	configContent := `check_interval_seconds: 300
listen_address: ":9100"
apt_check_path: "/usr/lib/update-notifier/apt-check"
update_stamp_path: "/var/lib/apt/periodic/update-success-stamp"
reboot_required_file: "/var/run/reboot-required"
log_level: "info"
command_timeout_seconds: 10
metrics_endpoint: "/metrics"
metric_prefix: "ubuntu"
`

	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write test config file: %v", err)
	}

	// Test loading the config
	cfg, err := Load(configPath)
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Verify config values
	if cfg.CheckIntervalSeconds != 300 {
		t.Errorf("Expected CheckIntervalSeconds to be 300, got %d", cfg.CheckIntervalSeconds)
	}
	if cfg.ListenAddress != ":9100" {
		t.Errorf("Expected ListenAddress to be ':9100', got %s", cfg.ListenAddress)
	}
	if cfg.AptCheckPath != "/usr/lib/update-notifier/apt-check" {
		t.Errorf("Expected AptCheckPath to be '/usr/lib/update-notifier/apt-check', got %s", cfg.AptCheckPath)
	}
	if cfg.MetricPrefix != "ubuntu" {
		t.Errorf("Expected MetricPrefix to be 'ubuntu', got %s", cfg.MetricPrefix)
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name        string
		config      Config
		expectError bool
	}{
		{
			name: "Valid config",
			config: Config{
				CheckIntervalSeconds:  300,
				ListenAddress:         ":9100",
				CommandTimeoutSeconds: 10,
				MetricsEndpoint:       "/metrics",
				MetricPrefix:          "ubuntu",
				LogLevel:              "info",
			},
			expectError: false,
		},
		{
			name: "Invalid check interval",
			config: Config{
				CheckIntervalSeconds:  0,
				ListenAddress:         ":9100",
				CommandTimeoutSeconds: 10,
				MetricsEndpoint:       "/metrics",
				MetricPrefix:          "ubuntu",
				LogLevel:              "info",
			},
			expectError: true,
		},
		{
			name: "Empty listen address",
			config: Config{
				CheckIntervalSeconds:  300,
				ListenAddress:         "",
				CommandTimeoutSeconds: 10,
				MetricsEndpoint:       "/metrics",
				MetricPrefix:          "ubuntu",
				LogLevel:              "info",
			},
			expectError: true,
		},
		{
			name: "Invalid command timeout",
			config: Config{
				CheckIntervalSeconds:  300,
				ListenAddress:         ":9100",
				CommandTimeoutSeconds: 0,
				MetricsEndpoint:       "/metrics",
				MetricPrefix:          "ubuntu",
				LogLevel:              "info",
			},
			expectError: true,
		},
		{
			name: "Empty metrics endpoint",
			config: Config{
				CheckIntervalSeconds:  300,
				ListenAddress:         ":9100",
				CommandTimeoutSeconds: 10,
				MetricsEndpoint:       "",
				MetricPrefix:          "ubuntu",
				LogLevel:              "info",
			},
			expectError: true,
		},
		{
			name: "Empty metric prefix",
			config: Config{
				CheckIntervalSeconds:  300,
				ListenAddress:         ":9100",
				CommandTimeoutSeconds: 10,
				MetricsEndpoint:       "/metrics",
				MetricPrefix:          "",
				LogLevel:              "info",
			},
			expectError: true,
		},
		{
			name: "Invalid log level",
			config: Config{
				CheckIntervalSeconds:  300,
				ListenAddress:         ":9100",
				CommandTimeoutSeconds: 10,
				MetricsEndpoint:       "/metrics",
				MetricPrefix:          "ubuntu",
				LogLevel:              "invalid",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.validate()
			if (err != nil) != tt.expectError {
				t.Errorf("validate() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}
