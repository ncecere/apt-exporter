package metrics

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
)

func TestNewMetrics(t *testing.T) {
	// Create metrics with a test prefix
	m := NewMetrics("test", true)

	// Check that all metrics are initialized
	if m.UpdatesAvailable == nil {
		t.Error("UpdatesAvailable metric is nil")
	}
	if m.SecurityUpdatesAvailable == nil {
		t.Error("SecurityUpdatesAvailable metric is nil")
	}
	if m.SecondsSinceLastUpdate == nil {
		t.Error("SecondsSinceLastUpdate metric is nil")
	}
	if m.RebootRequired == nil {
		t.Error("RebootRequired metric is nil")
	}
	if m.CollectionSuccess == nil {
		t.Error("CollectionSuccess metric is nil")
	}
	if m.CollectionDurationSeconds == nil {
		t.Error("CollectionDurationSeconds metric is nil")
	}
	if m.LastCollectionTimestamp == nil {
		t.Error("LastCollectionTimestamp metric is nil")
	}

	// Verify that metrics are registered with Prometheus
	// This is a bit tricky to test directly, so we'll just check that they implement the Collector interface
	_, ok := m.UpdatesAvailable.(prometheus.Collector)
	if !ok {
		t.Error("UpdatesAvailable does not implement prometheus.Collector")
	}
}

func TestNewTestMetrics(t *testing.T) {
	// Create test metrics
	m := NewTestMetrics()

	// Check that all metrics are initialized
	if m.UpdatesAvailable == nil {
		t.Error("UpdatesAvailable metric is nil")
	}
	if m.SecurityUpdatesAvailable == nil {
		t.Error("SecurityUpdatesAvailable metric is nil")
	}
	if m.SecondsSinceLastUpdate == nil {
		t.Error("SecondsSinceLastUpdate metric is nil")
	}
	if m.RebootRequired == nil {
		t.Error("RebootRequired metric is nil")
	}
	if m.CollectionSuccess == nil {
		t.Error("CollectionSuccess metric is nil")
	}
	if m.CollectionDurationSeconds == nil {
		t.Error("CollectionDurationSeconds metric is nil")
	}
	if m.LastCollectionTimestamp == nil {
		t.Error("LastCollectionTimestamp metric is nil")
	}

	// Test that we can get and set values
	testGauge := m.UpdatesAvailable.(*TestGauge)
	testGauge.Set(42)
	if testGauge.Get() != 42 {
		t.Errorf("Expected TestGauge value to be 42, got %f", testGauge.Get())
	}
}
