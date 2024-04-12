// SPDX-FileCopyrightText: 2024 Stefan Schärmeli <schaermu@pm.me>
// SPDX-License-Identifier: MIT
package collectors

import (
	"os"
	"path"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	promtest "github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/schaermu/changedetection.io-exporter/internal/testutil"
	"github.com/schaermu/changedetection.io-exporter/pkg/data"
)

var (
	expectedMetrics = []string{"changedetectionio_watch_price"}
)

func expectMetrics(t *testing.T, c prometheus.Collector, fixture string) {
	exp, err := os.Open(testutil.GetFixturePath(path.Join("metrics", fixture)))
	if err != nil {
		t.Fatalf("Error opening fixture file %q: %v", fixture, err)
	}
	if err := promtest.CollectAndCompare(c, exp, expectedMetrics...); err != nil {
		t.Fatalf("Unexpected metrics returned: %v", err)
	}
}

func expectMetricCount(t *testing.T, c prometheus.Collector, expected int) {
	count := promtest.CollectAndCount(c, expectedMetrics...)
	if count != expected {
		t.Fatalf("Expected %d metrics, got %d", expected, count)
	}
}

func createCollectorTestDb() map[string]*data.WatchItem {
	watchDb := testutil.NewWatchDb(0)
	uuid1, watch1 := testutil.NewTestItem("Item 1", 100, "USD")
	uuid2, watch2 := testutil.NewTestItem("Item 2", 200, "USD")
	watchDb[uuid1] = watch1
	watchDb[uuid2] = watch2
	return watchDb
}

func TestPriceCollector(t *testing.T) {
	watchDb := createCollectorTestDb()
	server := testutil.CreateTestApiServer(t, watchDb)
	defer server.Close()

	c, err := NewPriceCollector(server.URL(), "foo-bar-key")
	if err != nil {
		t.Fatal(err)
	}

	expectMetricCount(t, c, 2)
	expectMetrics(t, c, "price_metrics.prom")
}

/*
Test commented out for now due to flakiness, something seems to be off with the testing registry.
See https://github.com/schaermu/changedetection.io-exporter/issues/7.

func TestAutoUnregisterCollector(t *testing.T) {
	watchDb := createCollectorTestDb()
	server := testutil.CreateTestApiServer(t, watchDb)
	defer server.Close()

	c, err := NewPriceCollector(server.URL(), "foo-bar-key")
	if err != nil {
		t.Fatal(err)
	}

	keyToRemove := maps.Keys(watchDb)[len(watchDb)-1]
	delete(watchDb, keyToRemove)

	expectMetricCount(t, c, 1)
	expectMetrics(t, c, "price_metrics_autounregister.prom")
}
*/

/*
Test commented out because of some weirdness with the testing registry, see https://stackoverflow.com/questions/78297112/how-to-test-dynamic-metric-registration-in-custom-prometheus-exporter.
See issue https://github.com/schaermu/changedetection.io-exporter/issues/7.

func TestAutoregisterPriceCollector(t *testing.T) {
	watchDb := createCollectorTestDb()
	server := testutil.CreateTestApiServer(t, watchDb)
	defer server.Close()

	c, err := NewPriceCollector(server.URL(), "foo-bar-key")
	if err != nil {
		t.Fatal(err)
	}
	expectMetricCount(t, c, 2)

	// now add a new watch and expect the collector to pick it up
	uuid, newItem := testutil.NewTestItem("Item 3", 300, "USD")
	watchDb[uuid] = newItem

	expectMetricCount(t, c, 3)
	expectMetrics(t, c, "price_metrics_autoregister.prom")
}
*/
