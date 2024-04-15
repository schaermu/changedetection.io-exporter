// SPDX-FileCopyrightText: 2024 Stefan Schärmeli <schaermu@pm.me>
// SPDX-License-Identifier: MIT
package collectors

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/schaermu/changedetection.io-exporter/pkg/cdio"
	log "github.com/sirupsen/logrus"
)

type watchCollector struct {
	baseCollector

	checkCount             *prometheus.Desc
	fetchTime              *prometheus.Desc
	notificationAlertCount *prometheus.Desc
	lastCheckStatus        *prometheus.Desc
}

func NewWatchCollector(client *cdio.ApiClient) (*watchCollector, error) {
	return &watchCollector{
		baseCollector: *newBaseCollector(client),
		checkCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "watch", "check_count"),
			"Number of checks for a watch",
			labels, nil,
		),
		fetchTime: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "watch", "fetch_time"),
			"Time it took to fetch the watch",
			labels, nil,
		),
		notificationAlertCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "watch", "notification_alert_count"),
			"Number of notification alerts for a watch",
			labels, nil,
		),
		lastCheckStatus: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "watch", "last_check_status"),
			"Status of the last check for a watch",
			labels, nil,
		),
	}, nil
}

func (c *watchCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.checkCount
	ch <- c.fetchTime
	ch <- c.notificationAlertCount
	ch <- c.lastCheckStatus
}

func (c *watchCollector) Collect(ch chan<- prometheus.Metric) {
	c.RLock()
	defer c.RUnlock()

	// check for new watches before collecting metrics
	watches, err := c.ApiClient.GetWatches()
	if err != nil {
		log.Errorf("error while fetching watches: %v", err)
	}

	for uuid := range watches {
		// get latest watch data
		if watchData, err := c.ApiClient.GetWatchData(uuid); err == nil {
			if metricLabels, err := watchData.GetMetrics(); err != nil {
				log.Error(err)
				continue
			} else {
				ch <- prometheus.MustNewConstMetric(c.checkCount, prometheus.CounterValue, float64(watchData.CheckCount), metricLabels...)
				ch <- prometheus.MustNewConstMetric(c.fetchTime, prometheus.GaugeValue, watchData.FetchTime, metricLabels...)
				ch <- prometheus.MustNewConstMetric(c.notificationAlertCount, prometheus.CounterValue, float64(watchData.NotificationAlertCount), metricLabels...)
				ch <- prometheus.MustNewConstMetric(c.lastCheckStatus, prometheus.GaugeValue, float64(watchData.LastCheckStatus), metricLabels...)
			}
		} else {
			log.Error(err)
		}
	}
}
