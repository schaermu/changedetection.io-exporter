// SPDX-FileCopyrightText: 2024 Stefan Schärmeli <schaermu@pm.me>
// SPDX-License-Identifier: MIT
package collectors

import (
	"testing"

	"github.com/schaermu/changedetection.io-exporter/internal/testutil"
	"github.com/schaermu/changedetection.io-exporter/pkg/cdio"
)

var (
	expectedPriceMetrics = []string{"changedetectionio_watch_price"}
)

func TestPriceCollector(t *testing.T) {
	_, watchDb := testutil.NewCollectorTestDb()
	server := testutil.CreateTestApiServer(t, watchDb)
	defer server.Close()

	client := cdio.NewTestApiClient(server.URL())
	c := NewPriceCollector(client)

	testutil.ExpectMetricCount(t, c, 2, expectedPriceMetrics...)
	testutil.ExpectMetrics(t, c, "price_metrics.prom", expectedPriceMetrics...)
}

func TestPriceCollector_RemoveWatchDuringRuntime(t *testing.T) {
	keyToRemove, watchDb := testutil.NewCollectorTestDb()
	server := testutil.CreateTestApiServer(t, watchDb)
	defer server.Close()

	client := cdio.NewTestApiClient(server.URL())
	c := NewPriceCollector(client)

	testutil.ExpectMetricCount(t, c, 2, expectedPriceMetrics...)

	delete(watchDb, keyToRemove)

	testutil.ExpectMetricCount(t, c, 1, expectedPriceMetrics...)
	testutil.ExpectMetrics(t, c, "price_metrics_autounregister.prom", expectedPriceMetrics...)
}

func TestPriceCollector_NewWatchDuringRuntime(t *testing.T) {
	_, watchDb := testutil.NewCollectorTestDb()
	server := testutil.CreateTestApiServer(t, watchDb)
	defer server.Close()

	client := cdio.NewTestApiClient(server.URL())
	c := NewPriceCollector(client)

	testutil.ExpectMetricCount(t, c, 2, expectedPriceMetrics...)

	// now add a new watch and expect the collector to pick it up
	uuid, newItem := testutil.NewTestItem("Item 3", 300, "USD", 0, 0, 0)
	watchDb[uuid] = newItem

	testutil.ExpectMetricCount(t, c, 3, expectedPriceMetrics...)
	testutil.ExpectMetrics(t, c, "price_metrics_autoregister.prom", expectedPriceMetrics...)
}

func TestPriceCollector_HandlesArrayResponse(t *testing.T) {
	_, watchDb := testutil.NewCollectorTestDb()
	server := testutil.CreateTestApiServer(t, watchDb, testutil.WithPricesAsArray())
	defer server.Close()

	client := cdio.NewTestApiClient(server.URL())
	c := NewPriceCollector(client)

	testutil.ExpectMetricCount(t, c, 2, expectedPriceMetrics...)
	testutil.ExpectMetrics(t, c, "price_metrics.prom", expectedPriceMetrics...)
}
