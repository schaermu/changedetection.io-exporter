// SPDX-FileCopyrightText: 2024 Stefan Schärmeli <schaermu@pm.me>
// SPDX-License-Identifier: MIT
package main

import (
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

const (
	namespace = "changedetectionio"
)

type priceMetric struct {
	desc      *prometheus.Desc
	apiClient *ApiClient
	UUID      string
}

func newPriceMetric(labels prometheus.Labels, apiClient *ApiClient, uuid string) priceMetric {
	return priceMetric{
		desc: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "watch", "price"),
			"Current price of an offer type watch",
			nil, labels,
		),
		apiClient: apiClient,
		UUID:      uuid,
	}
}

func (m priceMetric) Describe(ch chan<- *prometheus.Desc) {
	ch <- m.desc
}

func (m priceMetric) Collect(ch chan<- prometheus.Metric) {
	// get latest snapshot
	if pData, err := m.apiClient.getLatestPriceSnapshot(m.UUID); err == nil {
		ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, pData.Price)
	} else {
		// error while fetching latest value for metric, unregister
		log.Errorf("error while fetching price snapshot %v", err)
		prometheus.Unregister(m)
	}
}

type priceCollector struct {
	ApiClient    *ApiClient
	priceMetrics map[string]priceMetric
}

func NewPriceCollector(endpoint string, key string) (*priceCollector, error) {
	// load all registered watches from changedetection.io API
	client := NewApiClient(endpoint, key)
	watches, err := client.getWatches()
	if err != nil {
		return nil, err
	}

	log.Infof("Loaded %d watches from changedetection.io API", len(watches))

	// configure price metrics for each watch
	priceMetrics := make(map[string]priceMetric)
	for id, watch := range watches {
		priceMetrics[id] = newPriceMetric(prometheus.Labels{"title": watch.Title}, client, id)
	}

	return &priceCollector{
		ApiClient:    client,
		priceMetrics: priceMetrics,
	}, nil
}

func (c *priceCollector) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range c.priceMetrics {
		metric.Describe(ch)
	}
}

func (c *priceCollector) Collect(ch chan<- prometheus.Metric) {
	// check for new watches before collecting metrics
	watches, err := c.ApiClient.getWatches()
	if err != nil {
		log.Errorf("error while fetching watches: %v", err)
	} else {
		for id, watch := range watches {
			if _, ok := c.priceMetrics[id]; !ok {
				c.priceMetrics[id] = newPriceMetric(prometheus.Labels{"title": watch.Title}, c.ApiClient, id)
				prometheus.MustRegister(c.priceMetrics[id])
				log.Infof("Picked up new watch %s, registered as metric %s", watch.Title, id)
			}
		}
	}

	for _, metric := range c.priceMetrics {
		metric.Collect(ch)
	}
}
