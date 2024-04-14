// SPDX-FileCopyrightText: 2024 Stefan Schärmeli <schaermu@pm.me>
// SPDX-License-Identifier: MIT
package collectors

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/schaermu/changedetection.io-exporter/pkg/cdio"
	log "github.com/sirupsen/logrus"
)

type priceCollector struct {
	baseCollector

	price *prometheus.Desc
}

func NewPriceCollector(client *cdio.ApiClient) (*priceCollector, error) {
	return &priceCollector{
		baseCollector: *newBaseCollector(client),
		price: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "watch", "price"),
			"Current price of an offer type watch",
			labels, nil,
		),
	}, nil
}

func (c *priceCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.price
}

func (c *priceCollector) Collect(ch chan<- prometheus.Metric) {
	c.RLock()
	defer c.RUnlock()

	// check for new watches before collecting metrics
	watches, err := c.ApiClient.GetWatches()
	if err != nil {
		log.Errorf("error while fetching watches: %v", err)
	}
	log.Infof("Collecting price metrics for %v watches", len(watches))

	for uuid, watch := range watches {
		// get latest price snapshot
		if pData, err := c.ApiClient.GetLatestPriceSnapshot(uuid); err == nil {
			ch <- prometheus.MustNewConstMetric(c.price, prometheus.GaugeValue, pData.Price, []string{watch.Title}...)
		}
	}
}
