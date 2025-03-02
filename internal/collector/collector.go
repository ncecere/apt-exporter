package collector

import (
	"context"
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

	// Create a context with timeout for running external commands
	cmdCtx, cancel := context.WithTimeout(ctx, time.Duration(c.cfg.CommandTimeoutSeconds)*time.Second)
	defer cancel()

	c.checkUpdates(cmdCtx)
	c.checkLastUpdateTime()
	c.checkRebootRequired()
}

// checkUpdates collects information about available updates.
func (c *Collector) checkUpdates(ctx context.Context) {
	if _, err := os.Stat(c.cfg.AptCheckPath); err != nil {
		c.logger.Printf("apt-check not found at %s: %v", c.cfg.AptCheckPath, err)
		return
	}

	out, err := exec.CommandContext(ctx, c.cfg.AptCheckPath).Output()
	if err != nil {
		c.logger.Printf("Error running apt-check: %v", err)
		return
	}

	parts := strings.Split(strings.TrimSpace(string(out)), ";")
	if len(parts) < 2 {
		c.logger.Printf("Unexpected output format from apt-check: %s", string(out))
		return
	}

	// Parse regular updates
	if n, err := strconv.Atoi(parts[0]); err == nil {
		c.metrics.UpdatesAvailable.Set(float64(n))
	} else {
		c.logger.Printf("Failed to parse update count: %v", err)
	}

	// Parse security updates
	if n, err := strconv.Atoi(parts[1]); err == nil {
		c.metrics.SecurityUpdatesAvailable.Set(float64(n))
	} else {
		c.logger.Printf("Failed to parse security update count: %v", err)
	}
}

// checkLastUpdateTime checks when the last update was performed.
func (c *Collector) checkLastUpdateTime() {
	info, err := os.Stat(c.cfg.UpdateStampPath)
	if err != nil {
		c.logger.Printf("Failed to stat update stamp file: %v", err)
		c.metrics.SecondsSinceLastUpdate.Set(0)
		return
	}

	seconds := time.Since(info.ModTime()).Seconds()
	c.metrics.SecondsSinceLastUpdate.Set(seconds)
}

// checkRebootRequired checks if a reboot is required.
func (c *Collector) checkRebootRequired() {
	_, err := os.Stat(c.cfg.RebootRequiredFile)
	if err == nil {
		c.metrics.RebootRequired.Set(1)
	} else if os.IsNotExist(err) {
		c.metrics.RebootRequired.Set(0)
	} else {
		c.logger.Printf("Error checking reboot required file: %v", err)
	}
}
