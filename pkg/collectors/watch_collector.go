// SPDX-FileCopyrightText: 2024 Stefan Schärmeli <schaermu@pm.me>
// SPDX-License-Identifier: MIT
package collectors

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/schaermu/changedetection.io-exporter/pkg/cdio"
	log "github.com/sirupsen/logrus"
)

type watchCollector struct {
	baseWatchCollector

	checkCount             *prometheus.Desc
	fetchTime              *prometheus.Desc
	notificationAlertCount *prometheus.Desc
	lastCheckStatus        *prometheus.Desc
}

func NewWatchCollector(endpoint string, key string) (*watchCollector, error) {
	// load all registered watches from changedetection.io API
	client := cdio.NewApiClient(endpoint, key)
	watches, err := client.GetWatches()
	if err != nil {
		return nil, err
	}

	log.Infof("Loaded %d watches from changedetection.io API", len(watches))

	return &watchCollector{
		baseWatchCollector: *newBaseWatchCollector(client, watches),
		checkCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "watch", "check_count"),
			"Number of checks for a watch",
			[]string{"title"}, nil,
		),
		fetchTime: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "watch", "fetch_time"),
			"Time it took to fetch the watch",
			[]string{"title"}, nil,
		),
		notificationAlertCount: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "watch", "notification_alert_count"),
			"Number of notification alerts for a watch",
			[]string{"title"}, nil,
		),
		lastCheckStatus: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "watch", "last_check_status"),
			"Status of the last check for a watch",
			[]string{"title"}, nil,
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
	} else {
		c.watches = watches
	}

	for uuid, watch := range c.watches {
		// get latest watch data
		if watchData, err := c.ApiClient.GetWatchData(uuid); err == nil {
			ch <- prometheus.MustNewConstMetric(c.checkCount, prometheus.GaugeValue, float64(watchData.CheckCount), []string{watch.Title}...)
			ch <- prometheus.MustNewConstMetric(c.fetchTime, prometheus.GaugeValue, watchData.FetchTime, []string{watch.Title}...)
			ch <- prometheus.MustNewConstMetric(c.notificationAlertCount, prometheus.GaugeValue, float64(watchData.NotificationAlertCount), []string{watch.Title}...)
			ch <- prometheus.MustNewConstMetric(c.lastCheckStatus, prometheus.GaugeValue, float64(watchData.LastCheckStatus), []string{watch.Title}...)
		}
	}
}
