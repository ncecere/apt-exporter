package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Gauge is an interface that allows us to use both prometheus.Gauge and test gauges
type Gauge interface {
	Set(float64)
}

// TestGauge is a mock implementation of Gauge for testing
type TestGauge struct {
	value float64
}

// Set sets the gauge to an arbitrary value
func (g *TestGauge) Set(val float64) {
	g.value = val
}

// Get returns the current value of the gauge (for testing)
func (g *TestGauge) Get() float64 {
	return g.value
}

// Metrics holds all the Prometheus metrics for the APT exporter.
type Metrics struct {
	// Core metrics
	UpdatesAvailable         Gauge
	SecurityUpdatesAvailable Gauge
	SecondsSinceLastUpdate   Gauge
	RebootRequired           Gauge

	// Collector metrics
	CollectionSuccess         Gauge
	CollectionDurationSeconds Gauge
	LastCollectionTimestamp   Gauge
}

// NewMetrics creates and registers all metrics with the provided prefix.
func NewMetrics(prefix string) *Metrics {
	m := &Metrics{
		// Core metrics
		UpdatesAvailable: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: prefix + "_updates_available",
			Help: "Number of available package updates",
		}),
		SecurityUpdatesAvailable: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: prefix + "_security_updates_available",
			Help: "Number of available security updates",
		}),
		SecondsSinceLastUpdate: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: prefix + "_seconds_since_last_update",
			Help: "Seconds since last successful apt update",
		}),
		RebootRequired: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: prefix + "_reboot_required",
			Help: "1 if a reboot is required, 0 otherwise",
		}),

		// Collector metrics
		CollectionSuccess: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: prefix + "_collector_success",
			Help: "1 if the last collection was successful, 0 otherwise",
		}),
		CollectionDurationSeconds: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: prefix + "_collector_duration_seconds",
			Help: "Duration of the last collection in seconds",
		}),
		LastCollectionTimestamp: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: prefix + "_collector_last_timestamp",
			Help: "Timestamp of the last collection",
		}),
	}

	// Register all metrics with Prometheus
	prometheus.MustRegister(
		// Core metrics
		m.UpdatesAvailable.(prometheus.Collector),
		m.SecurityUpdatesAvailable.(prometheus.Collector),
		m.SecondsSinceLastUpdate.(prometheus.Collector),
		m.RebootRequired.(prometheus.Collector),

		// Collector metrics
		m.CollectionSuccess.(prometheus.Collector),
		m.CollectionDurationSeconds.(prometheus.Collector),
		m.LastCollectionTimestamp.(prometheus.Collector),
	)

	return m
}

// NewTestMetrics creates a new Metrics instance with TestGauge implementations for testing
func NewTestMetrics() *Metrics {
	return &Metrics{
		// Core metrics
		UpdatesAvailable:         &TestGauge{},
		SecurityUpdatesAvailable: &TestGauge{},
		SecondsSinceLastUpdate:   &TestGauge{},
		RebootRequired:           &TestGauge{},

		// Collector metrics
		CollectionSuccess:         &TestGauge{},
		CollectionDurationSeconds: &TestGauge{},
		LastCollectionTimestamp:   &TestGauge{},
	}
}
