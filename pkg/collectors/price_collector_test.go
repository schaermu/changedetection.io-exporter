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
	"github.com/schaermu/changedetection.io-exporter/pkg/cdio"
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

// createCollectorTestDb creates a test database with two watch items and returns the UUID of the second item and the database.
func createCollectorTestDb() (string, map[string]*data.WatchItem) {
	watchDb := testutil.NewWatchDb(0)
	uuid1, watch1 := testutil.NewTestItem("Item 1", 100, "USD")
	uuid2, watch2 := testutil.NewTestItem("Item 2", 200, "USD")
	watchDb[uuid1] = watch1
	watchDb[uuid2] = watch2
	return uuid2, watchDb
}

func createTestClient(url string) *cdio.ApiClient {
	return cdio.NewApiClient(url, "foo-bar-key")
}

func TestPriceCollector(t *testing.T) {
	_, watchDb := createCollectorTestDb()
	server := testutil.CreateTestApiServer(t, watchDb)
	defer server.Close()

	client := createTestClient(server.URL())
	c, err := NewPriceCollector(client)
	if err != nil {
		t.Fatal(err)
	}

	expectMetricCount(t, c, 2)
	expectMetrics(t, c, "price_metrics.prom")
}

func TestPriceCollector_RemoveWatchDuringRuntime(t *testing.T) {
	keyToRemove, watchDb := createCollectorTestDb()
	server := testutil.CreateTestApiServer(t, watchDb)
	defer server.Close()

	client := createTestClient(server.URL())
	c, err := NewPriceCollector(client)
	if err != nil {
		t.Fatal(err)
	}
	expectMetricCount(t, c, 2)

	delete(watchDb, keyToRemove)

	expectMetricCount(t, c, 1)
	expectMetrics(t, c, "price_metrics_autounregister.prom")
}

func TestPriceCollector_NewWatchDuringRuntime(t *testing.T) {
	_, watchDb := createCollectorTestDb()
	server := testutil.CreateTestApiServer(t, watchDb)
	defer server.Close()

	client := createTestClient(server.URL())
	c, err := NewPriceCollector(client)
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

func TestPriceCollector_HandlesArrayResponse(t *testing.T) {
	_, watchDb := createCollectorTestDb()
	server := testutil.CreateTestApiServer(t, watchDb, testutil.WithPricesAsArray())
	defer server.Close()

	client := createTestClient(server.URL())
	c, err := NewPriceCollector(client)
	if err != nil {
		t.Fatal(err)
	}

	expectMetricCount(t, c, 2)
	expectMetrics(t, c, "price_metrics.prom")
}
