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
)

func expectMetrics(t *testing.T, c prometheus.Collector, fixture string, metricNames ...string) {
	exp, err := os.Open(testutil.GetFixturePath(path.Join("metrics", fixture)))
	if err != nil {
		t.Fatalf("Error opening fixture file %q: %v", fixture, err)
	}
	if err := promtest.CollectAndCompare(c, exp, metricNames...); err != nil {
		t.Fatalf("Unexpected metrics returned: %v", err)
	}
}

func expectMetricCount(t *testing.T, c prometheus.Collector, expected int, metricnames ...string) {
	count := promtest.CollectAndCount(c, metricnames...)
	if count != expected {
		t.Fatalf("Expected %d metrics, got %d", expected, count)
	}
}

func TestPriceCollector(t *testing.T) {
	server := testutil.CreateTestApiServer(t, map[string]string{
		"/api/v1/watch": testutil.GetFixturePath("json/getWatches.json"),
		"/api/v1/watch/6a4b7d5c-fee4-4616-9f43-4ac97046b595/history/latest": testutil.GetFixturePath("json/getLatestPriceSnapshot_4ac97046b595.json"),
		"/api/v1/watch/e6f5fd5c-dbfe-468b-b8f3-f9d6ff5ad69b/history/latest": testutil.GetFixturePath("json/getLatestPriceSnapshot_f9d6ff5ad69b.json"),
	})
	defer server.Close()
	c, err := NewPriceCollector(server.URL, "foo-bar-key")
	if err != nil {
		t.Fatal(err)
	}
	expectMetrics(t, c, "price_metrics.prom", "changedetectionio_watch_price")
	expectMetricCount(t, c, 2, "changedetectionio_watch_price")
}

func TestAutoUnregisterCollector(t *testing.T) {
	// we deliberately do not register a history for one of the watches to provoke an unregister event
	server := testutil.CreateTestApiServer(t, map[string]string{
		"/api/v1/watch": testutil.GetFixturePath("json/getWatches.json"),
		"/api/v1/watch/6a4b7d5c-fee4-4616-9f43-4ac97046b595/history/latest": testutil.GetFixturePath("json/getLatestPriceSnapshot_4ac97046b595.json"),
	})
	defer server.Close()
	c, err := NewPriceCollector(server.URL, "foo-bar-key")
	if err != nil {
		t.Fatal(err)
	}
	expectMetrics(t, c, testutil.GetFixturePath("metrics/price_metrics_autounregister.prom"), "changedetectionio_watch_price")
	expectMetricCount(t, c, 1, "changedetectionio_watch_price")
}

/*
func TestAutoregisterPriceCollector(t *testing.T) {
	server := CreateTestApiServer(t, map[string]string{
		"/api/v1/watch": "./test/json/getWatches_single.json",
		"/api/v1/watch/6a4b7d5c-fee4-4616-9f43-4ac97046b595/history/latest": "./test/json/getLatestPriceSnapshot_6a4b7d5c-fee4-4616-9f43-4ac97046b595.json",
	})
	defer server.Close()
	c, err := NewPriceCollector(server.URL, "foo-bar-key")
	if err != nil {
		t.Fatal(err)
	}
	expectMetricCount(t, c, 1, "changedetectionio_watch_price")

	// now add a new watch and expect the collector to pick it up
	newServer := CreateTestApiServer(t, map[string]string{
		"/api/v1/watch": "./test/json/getWatches.json",
		"/api/v1/watch/6a4b7d5c-fee4-4616-9f43-4ac97046b595/history/latest": "./test/json/getLatestPriceSnapshot_6a4b7d5c-fee4-4616-9f43-4ac97046b595.json",
		"/api/v1/watch/e6f5fd5c-dbfe-468b-b8f3-f9d6ff5ad69b/history/latest": "./test/json/getLatestPriceSnapshot_e6f5fd5c-dbfe-468b-b8f3-f9d6ff5ad69b.json",
	})
	defer newServer.Close()
	c.ApiClient.SetBaseUrl(newServer.URL)

	expectMetrics(t, c, "price_metrics.prom", "changedetectionio_watch_price")
	expectMetricCount(t, c, 2, "changedetectionio_watch_price")
}
*/
