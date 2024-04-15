// SPDX-FileCopyrightText: 2024 Stefan Schärmeli <schaermu@pm.me>
// SPDX-License-Identifier: MIT
package collectors

import (
	"testing"

	"github.com/schaermu/changedetection.io-exporter/internal/testutil"
	"github.com/schaermu/changedetection.io-exporter/pkg/cdio"
	"github.com/schaermu/changedetection.io-exporter/pkg/data"
)

var (
	expectedSystemMetrics = []string{
		"changedetectionio_system_overdue_watch_count",
		"changedetectionio_system_queue_size",
		"changedetectionio_system_uptime",
		"changedetectionio_system_watch_count",
	}
)

func TestSystemCollector(t *testing.T) {
	lastId, watchDb := testutil.NewCollectorTestDb()
	server := testutil.CreateTestApiServer(t, watchDb, testutil.WithSystemInfo(&data.SystemInfo{
		QueueSize:      0,
		WatchCount:     len(watchDb),
		OverdueWatches: []string{lastId},
		Uptime:         1111.1,
		Version:        "0.1.1",
	}))
	defer server.Close()

	client := cdio.NewTestApiClient(server.URL())
	c, err := NewSystemCollector(client)
	if err != nil {
		t.Fatal(err)
	}

	testutil.ExpectMetricCount(t, c, 1, expectedSystemMetrics...)
	testutil.ExpectMetrics(t, c, "system_metrics.prom", expectedSystemMetrics...)
}
