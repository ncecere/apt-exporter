package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ncecere/apt-exporter/internal/collector"
	"github.com/ncecere/apt-exporter/internal/config"
	"github.com/ncecere/apt-exporter/internal/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Version information set by build flags
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	// Parse command-line flags
	configPath := flag.String("config", "config.yml", "Path to YAML configuration file")
	showVersion := flag.Bool("version", false, "Show version information")
	flag.Parse()

	// Show version information if requested
	if *showVersion {
		fmt.Printf("apt-exporter version %s (commit: %s, built at: %s)\n", version, commit, date)
		os.Exit(0)
	}

	// Set up logging
	logger := log.New(os.Stdout, "apt-exporter: ", log.LstdFlags)
	logger.Printf("Starting APT exporter version %s", version)

	// Load configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		logger.Fatalf("Failed to load configuration: %v", err)
	}
	logger.Printf("Configuration loaded from %s", *configPath)

	// Configure logging level (basic implementation)
	// A more sophisticated logging library could be used here

	// Initialize metrics
	m := metrics.NewMetrics(cfg.MetricPrefix)
	logger.Printf("Metrics initialized with prefix: %s", cfg.MetricPrefix)

	// Create collector
	c := collector.New(cfg, m)

	// Set up HTTP server for metrics endpoint
	http.Handle(cfg.MetricsEndpoint, promhttp.Handler())
	logger.Printf("Metrics endpoint registered at %s", cfg.MetricsEndpoint)

	// Create a context that will be canceled on SIGINT or SIGTERM
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigChan
		logger.Printf("Received signal: %v", sig)
		cancel()
	}()

	// Start metrics collection in a goroutine
	go c.Start(ctx)

	// Start HTTP server
	server := &http.Server{
		Addr: cfg.ListenAddress,
	}

	// Start server in a goroutine
	go func() {
		logger.Printf("Starting metrics server on %s", cfg.ListenAddress)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("HTTP server error: %v", err)
		}
	}()

	// Wait for context cancellation (from signal handler)
	<-ctx.Done()
	logger.Println("Shutting down...")

	// Create a context with timeout for graceful shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Printf("HTTP server shutdown error: %v", err)
	}

	logger.Println("APT exporter stopped")
}
