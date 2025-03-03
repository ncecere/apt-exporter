package collector

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ncecere/apt-exporter/internal/config"
	"github.com/ncecere/apt-exporter/internal/metrics"
)

func TestCollector(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir := t.TempDir()

	// Create mock files
	aptCheckPath := filepath.Join(tmpDir, "apt-check")
	updateStampPath := filepath.Join(tmpDir, "update-success-stamp")
	rebootRequiredFile := filepath.Join(tmpDir, "reboot-required")

	// Create a mock apt-check script that outputs "5;2" to stderr (like the real apt-check)
	mockAptCheckContent := `#!/bin/sh
echo "5;2" >&2
`
	if err := os.WriteFile(aptCheckPath, []byte(mockAptCheckContent), 0755); err != nil {
		t.Fatalf("Failed to create mock apt-check: %v", err)
	}

	// Create a mock update stamp file
	updateTime := time.Now().Add(-24 * time.Hour) // 1 day ago
	if err := os.WriteFile(updateStampPath, []byte(""), 0644); err != nil {
		t.Fatalf("Failed to create mock update stamp: %v", err)
	}
	if err := os.Chtimes(updateStampPath, updateTime, updateTime); err != nil {
		t.Fatalf("Failed to set mock update stamp time: %v", err)
	}

	// Create a mock reboot required file
	if err := os.WriteFile(rebootRequiredFile, []byte(""), 0644); err != nil {
		t.Fatalf("Failed to create mock reboot required file: %v", err)
	}

	// Create a test configuration
	cfg := &config.Config{
		CheckIntervalSeconds:  300,
		ListenAddress:         ":9100",
		AptCheckPath:          aptCheckPath,
		UpdateStampPath:       updateStampPath,
		RebootRequiredFile:    rebootRequiredFile,
		LogLevel:              "info",
		CommandTimeoutSeconds: 10,
		MetricsEndpoint:       "/metrics",
		MetricPrefix:          "test",
	}

	// Create metrics (using the test metrics implementation)
	m := metrics.NewTestMetrics()

	// Create collector
	c := New(cfg, m)

	// Create a context with a short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// Collect metrics
	c.collect(ctx)

	// Check metrics values
	updatesAvailable := m.UpdatesAvailable.(*metrics.TestGauge).Get()
	if updatesAvailable != 5 {
		t.Errorf("Expected UpdatesAvailable to be 5, got %f", updatesAvailable)
	}

	securityUpdates := m.SecurityUpdatesAvailable.(*metrics.TestGauge).Get()
	if securityUpdates != 2 {
		t.Errorf("Expected SecurityUpdatesAvailable to be 2, got %f", securityUpdates)
	}

	// Check seconds since last update (should be around 24 hours = 86400 seconds)
	// Allow for some tolerance in the comparison
	secondsSinceUpdate := m.SecondsSinceLastUpdate.(*metrics.TestGauge).Get()
	if secondsSinceUpdate < 86000 || secondsSinceUpdate > 87000 {
		t.Errorf("Expected SecondsSinceLastUpdate to be around 86400, got %f", secondsSinceUpdate)
	}

	rebootRequired := m.RebootRequired.(*metrics.TestGauge).Get()
	if rebootRequired != 1 {
		t.Errorf("Expected RebootRequired to be 1, got %f", rebootRequired)
	}

	collectionSuccess := m.CollectionSuccess.(*metrics.TestGauge).Get()
	if collectionSuccess != 1 {
		t.Errorf("Expected CollectionSuccess to be 1, got %f", collectionSuccess)
	}

	// Test with missing reboot required file
	if err := os.Remove(rebootRequiredFile); err != nil {
		t.Fatalf("Failed to remove mock reboot required file: %v", err)
	}

	// Collect metrics again
	c.collect(ctx)

	// Check that reboot required is now 0
	rebootRequired = m.RebootRequired.(*metrics.TestGauge).Get()
	if rebootRequired != 0 {
		t.Errorf("Expected RebootRequired to be 0 after file removal, got %f", rebootRequired)
	}

	// Test with invalid apt-check script (outputs to stderr like the real apt-check)
	invalidAptCheckContent := `#!/bin/sh
echo "invalid" >&2
`
	if err := os.WriteFile(aptCheckPath, []byte(invalidAptCheckContent), 0755); err != nil {
		t.Fatalf("Failed to update mock apt-check: %v", err)
	}

	// Collect metrics again
	c.collect(ctx)

	// Check that collection success is still 1 (we now handle invalid output gracefully)
	collectionSuccess = m.CollectionSuccess.(*metrics.TestGauge).Get()
	if collectionSuccess != 1 {
		t.Errorf("Expected CollectionSuccess to be 1 after invalid apt-check (handled gracefully), got %f", collectionSuccess)
	}

	// Check that updates are set to 0 when format is invalid
	updatesAvailable = m.UpdatesAvailable.(*metrics.TestGauge).Get()
	if updatesAvailable != 0 {
		t.Errorf("Expected UpdatesAvailable to be 0 after invalid apt-check, got %f", updatesAvailable)
	}
}
