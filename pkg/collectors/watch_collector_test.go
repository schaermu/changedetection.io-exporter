// SPDX-FileCopyrightText: 2024 Stefan Schärmeli <schaermu@pm.me>
// SPDX-License-Identifier: MIT
package collectors

import (
	"testing"

	"github.com/schaermu/changedetection.io-exporter/internal/testutil"
	"github.com/schaermu/changedetection.io-exporter/pkg/cdio"
)

var (
	expectedWatchMetrics = []string{
		"changedetectionio_watch_check_count",
		"changedetectionio_watch_fetch_time",
		"changedetectionio_watch_notification_alert_count",
		"changedetectionio_watch_last_check_status",
	}
)

func TestWatchCollector(t *testing.T) {
	_, watchDb := testutil.NewCollectorTestDb()
	server := testutil.CreateTestApiServer(t, watchDb)
	defer server.Close()

	client := cdio.NewTestApiClient(server.URL())
	c := NewWatchCollector(client)

	//testutil.ExpectMetricCount(t, c, 2, expectedWatchMetrics...)
	testutil.ExpectMetrics(t, c, "watch_metrics.prom", expectedWatchMetrics...)
}
