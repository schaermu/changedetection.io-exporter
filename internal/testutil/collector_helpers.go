package testutil

import (
	"os"
	"path"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/schaermu/changedetection.io-exporter/pkg/data"
)

// createCollectorTestDb creates a test database with two watch items and returns the UUID of the second item and the database.
func NewCollectorTestDb() (string, map[string]*data.WatchItem) {
	watchDb := NewWatchDb(0)
	uuid1, watch1 := NewTestItem("Item 1", 100, "USD", 20, 15, 10)
	uuid2, watch2 := NewTestItem("Item 2", 200, "USD", 20, 15, 10)
	watchDb[uuid1] = watch1
	watchDb[uuid2] = watch2
	return uuid2, watchDb
}

func ExpectMetrics(t *testing.T, c prometheus.Collector, fixture string, expectedMetrics ...string) {
	exp, err := os.Open(GetFixturePath(path.Join("metrics", fixture)))
	if err != nil {
		t.Fatalf("Error opening fixture file %q: %v", fixture, err)
	}
	if err := testutil.CollectAndCompare(c, exp, expectedMetrics...); err != nil {
		t.Fatalf("Unexpected metrics returned: %v", err)
	}
}

func ExpectMetricCount(t *testing.T, c prometheus.Collector, expected int, expectedMetrics ...string) {
	count := testutil.CollectAndCount(c, expectedMetrics...)
	total := int(expected * len(expectedMetrics))
	if count != total {
		t.Fatalf("Expected %d metrics, got %d", total, count)
	}
}
