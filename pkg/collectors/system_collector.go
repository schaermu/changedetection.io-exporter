// SPDX-FileCopyrightText: 2024 Stefan Schärmeli <schaermu@pm.me>
// SPDX-License-Identifier: MIT
package collectors

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/schaermu/changedetection.io-exporter/pkg/cdio"
	log "github.com/sirupsen/logrus"
)

type systemCollector struct {
	baseCollector

	queueSize    *prometheus.Desc
	overdueCount *prometheus.Desc
	uptime       *prometheus.Desc
	watchCount   *prometheus.Desc
}

func NewSystemCollector(client *cdio.ApiClient) (*systemCollector, error) {
	return &systemCollector{
		baseCollector: *newBaseCollector(client),
		queueSize: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "system", "queue_size"),
			"Current changedetection.io instance queue size",
			[]string{"version"}, nil,
		),
		watchCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "system", "watch_count"),
			"Current changedetection.io instance watch count",
			[]string{"version"}, nil,
		),
		overdueCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "system", "overdue_watch_count"),
			"Current changedetection.io instance overdue watch count",
			[]string{"version"}, nil,
		),
		uptime: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "system", "uptime"),
			"Current changedetection.io instance system uptime",
			[]string{"version"}, nil,
		),
	}, nil
}

func (c *systemCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.queueSize
	ch <- c.watchCount
	ch <- c.overdueCount
	ch <- c.uptime
}

func (c *systemCollector) Collect(ch chan<- prometheus.Metric) {
	c.RLock()
	defer c.RUnlock()

	// check for new watches before collecting metrics
	system, err := c.ApiClient.GetSystemInfo()
	if err != nil {
		log.Errorf("error while fetching system info: %v", err)
	} else {
		ch <- prometheus.MustNewConstMetric(c.queueSize, prometheus.GaugeValue, float64(system.QueueSize), system.Version)
		ch <- prometheus.MustNewConstMetric(c.watchCount, prometheus.GaugeValue, float64(system.WatchCount), system.Version)
		ch <- prometheus.MustNewConstMetric(c.overdueCount, prometheus.GaugeValue, float64(len(system.OverdueWatches)), system.Version)
		ch <- prometheus.MustNewConstMetric(c.uptime, prometheus.GaugeValue, float64(system.Uptime), system.Version)
	}
}
