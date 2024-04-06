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

// Define a struct for you collector that contains pointers
// to prometheus descriptors for each metric you wish to expose.
// Note you can also include fields of other types if they provide utility
// but we just won't be exposing them as metrics.
type priceMetric struct {
	desc      *prometheus.Desc
	apiClient *cdioApiClient
	UUID      string
}

func newPriceMetric(labels prometheus.Labels, apiClient *cdioApiClient, uuid string) priceMetric {
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
		ch <- prometheus.MustNewConstMetric(m.desc, prometheus.GaugeValue, float64(pData.Price))
	} else {
		// error while fetching latest value for metric, unregister
		log.Errorf("error while fetching price snapshot %v", err)
		prometheus.Unregister(m)
	}
}

type priceCollector struct {
	priceMetrics map[string]priceMetric
	apiClient    *cdioApiClient
}

func newPriceCollector(endpoint string, key string) (*priceCollector, error) {
	// load all registered watches from changedetection.io API
	client := newCdioApiClient(endpoint, key)
	watches := client.getWatches()

	log.Infof("Loaded %d watches from changedetection.io API", len(watches))

	// configure price metrics for each watch
	priceMetrics := make(map[string]priceMetric)
	for id, watch := range watches {
		priceMetrics[id] = newPriceMetric(prometheus.Labels{"title": watch.Title}, client, id)
	}

	return &priceCollector{
		apiClient:    client,
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
	watches := c.apiClient.getWatches()
	for id, watch := range watches {
		if _, ok := c.priceMetrics[id]; !ok {
			c.priceMetrics[id] = newPriceMetric(prometheus.Labels{"title": watch.Title}, c.apiClient, id)
			prometheus.MustRegister(c.priceMetrics[id])
		}
	}

	for _, metric := range c.priceMetrics {
		metric.Collect(ch)
	}
}
