// SPDX-FileCopyrightText: 2024 Stefan Schärmeli <schaermu@pm.me>
// SPDX-License-Identifier: MIT
package collectors

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/schaermu/changedetection.io-exporter/pkg/cdio"
	"github.com/schaermu/changedetection.io-exporter/pkg/data"
	log "github.com/sirupsen/logrus"
)

const (
	namespace = "changedetectionio"
)

type priceCollector struct {
	sync.RWMutex

	ApiClient *cdio.ApiClient
	watches   map[string]data.WatchItem
	price     *prometheus.Desc
}

func NewPriceCollector(endpoint string, key string) (*priceCollector, error) {
	// load all registered watches from changedetection.io API
	client := cdio.NewApiClient(endpoint, key)
	watches, err := client.GetWatches()
	if err != nil {
		return nil, err
	}

	log.Infof("Loaded %d watches from changedetection.io API", len(watches))

	return &priceCollector{
		ApiClient: client,
		watches:   watches,
		price: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "watch", "price"),
			"Current price of an offer type watch",
			[]string{"title"}, nil,
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
	} else {
		c.watches = watches
	}

	for uuid, watch := range c.watches {
		// get latest price snapshot
		if pData, err := c.ApiClient.GetLatestPriceSnapshot(uuid); err == nil {
			ch <- prometheus.MustNewConstMetric(c.price, prometheus.GaugeValue, pData.Price, []string{watch.Title}...)
		}
	}
}
