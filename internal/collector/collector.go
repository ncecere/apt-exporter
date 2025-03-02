package collector

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/ncecere/apt-exporter/internal/config"
	"github.com/ncecere/apt-exporter/internal/metrics"
)

// Collector handles the collection of APT metrics.
type Collector struct {
	cfg     *config.Config
	metrics *metrics.Metrics
	logger  *log.Logger
}

// New creates a new Collector instance.
func New(cfg *config.Config, metrics *metrics.Metrics) *Collector {
	// Create a logger with appropriate prefix
	logger := log.New(os.Stdout, "apt-collector: ", log.LstdFlags)

	return &Collector{
		cfg:     cfg,
		metrics: metrics,
		logger:  logger,
	}
}

// Start begins periodic collection of metrics.
func (c *Collector) Start(ctx context.Context) {
	// Collect metrics immediately on startup
	c.collect(ctx)

	// Set up ticker for periodic collection
	ticker := time.NewTicker(time.Duration(c.cfg.CheckIntervalSeconds) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			c.collect(ctx)
		case <-ctx.Done():
			c.logger.Println("Stopping metrics collection")
			return
		}
	}
}

// collect gathers all metrics and updates the Prometheus gauges.
func (c *Collector) collect(ctx context.Context) {
	c.logger.Println("Collecting APT metrics")
	startTime := time.Now()

	// Create a context with timeout for running external commands
	cmdCtx, cancel := context.WithTimeout(ctx, time.Duration(c.cfg.CommandTimeoutSeconds)*time.Second)
	defer cancel()

	// Track collection success
	success := true

	// Collect metrics and track success
	if err := c.checkUpdates(cmdCtx); err != nil {
		c.logger.Printf("Error checking updates: %v", err)
		success = false
	}

	if err := c.checkLastUpdateTime(); err != nil {
		c.logger.Printf("Error checking last update time: %v", err)
		success = false
	}

	if err := c.checkRebootRequired(); err != nil {
		c.logger.Printf("Error checking reboot required: %v", err)
		success = false
	}

	// Update collection metrics
	c.metrics.CollectionSuccess.Set(boolToFloat64(success))
	c.metrics.CollectionDurationSeconds.Set(time.Since(startTime).Seconds())
	c.metrics.LastCollectionTimestamp.Set(float64(time.Now().Unix()))
}

// checkUpdates collects information about available updates.
func (c *Collector) checkUpdates(ctx context.Context) error {
	if _, err := os.Stat(c.cfg.AptCheckPath); err != nil {
		c.metrics.UpdatesAvailable.Set(0)
		c.metrics.SecurityUpdatesAvailable.Set(0)
		return fmt.Errorf("apt-check not found at %s: %w", c.cfg.AptCheckPath, err)
	}

	out, err := exec.CommandContext(ctx, c.cfg.AptCheckPath).Output()
	if err != nil {
		c.metrics.UpdatesAvailable.Set(0)
		c.metrics.SecurityUpdatesAvailable.Set(0)
		return fmt.Errorf("error running apt-check: %w", err)
	}

	parts := strings.Split(strings.TrimSpace(string(out)), ";")
	if len(parts) < 2 {
		c.metrics.UpdatesAvailable.Set(0)
		c.metrics.SecurityUpdatesAvailable.Set(0)
		return fmt.Errorf("unexpected output format from apt-check: %s", string(out))
	}

	// Parse regular updates
	regularUpdates, err := strconv.Atoi(parts[0])
	if err != nil {
		c.metrics.UpdatesAvailable.Set(0)
		return fmt.Errorf("failed to parse update count: %w", err)
	}
	c.metrics.UpdatesAvailable.Set(float64(regularUpdates))

	// Parse security updates
	securityUpdates, err := strconv.Atoi(parts[1])
	if err != nil {
		c.metrics.SecurityUpdatesAvailable.Set(0)
		return fmt.Errorf("failed to parse security update count: %w", err)
	}
	c.metrics.SecurityUpdatesAvailable.Set(float64(securityUpdates))

	return nil
}

// checkLastUpdateTime checks when the last update was performed.
func (c *Collector) checkLastUpdateTime() error {
	info, err := os.Stat(c.cfg.UpdateStampPath)
	if err != nil {
		c.metrics.SecondsSinceLastUpdate.Set(0)
		return fmt.Errorf("failed to stat update stamp file: %w", err)
	}

	seconds := time.Since(info.ModTime()).Seconds()
	c.metrics.SecondsSinceLastUpdate.Set(seconds)
	return nil
}

// checkRebootRequired checks if a reboot is required.
func (c *Collector) checkRebootRequired() error {
	_, err := os.Stat(c.cfg.RebootRequiredFile)
	if err == nil {
		c.metrics.RebootRequired.Set(1)
		return nil
	} else if os.IsNotExist(err) {
		c.metrics.RebootRequired.Set(0)
		return nil
	}

	c.metrics.RebootRequired.Set(0)
	return fmt.Errorf("error checking reboot required file: %w", err)
}

// boolToFloat64 converts a boolean to a float64 (1.0 for true, 0.0 for false)
func boolToFloat64(b bool) float64 {
	if b {
		return 1.0
	}
	return 0.0
}
