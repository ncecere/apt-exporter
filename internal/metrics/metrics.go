package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Metrics holds all the Prometheus metrics for the APT exporter.
type Metrics struct {
	UpdatesAvailable         prometheus.Gauge
	SecurityUpdatesAvailable prometheus.Gauge
	SecondsSinceLastUpdate   prometheus.Gauge
	RebootRequired           prometheus.Gauge
}

// NewMetrics creates and registers all metrics with the provided prefix.
func NewMetrics(prefix string) *Metrics {
	m := &Metrics{
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
	}

	// Register all metrics with Prometheus
	prometheus.MustRegister(
		m.UpdatesAvailable,
		m.SecurityUpdatesAvailable,
		m.SecondsSinceLastUpdate,
		m.RebootRequired,
	)

	return m
}
