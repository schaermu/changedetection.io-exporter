// SPDX-FileCopyrightText: 2024 Stefan Schärmeli <schaermu@pm.me>
// SPDX-License-Identifier: MIT
package main

import (
	"fmt"
	"net/http"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/schaermu/changedetection.io-exporter/pkg/cdio"
	"github.com/schaermu/changedetection.io-exporter/pkg/collectors"

	promcollectors "github.com/prometheus/client_golang/prometheus/collectors"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.000000",
	})

	log.SetLevel(log.DebugLevel)
}

func main() {
	var (
		port   = os.Getenv("PORT")
		apiUrl = os.Getenv("CDIO_API_BASE_URL")
		apiKey = os.Getenv("CDIO_API_KEY")
	)

	if port == "" {
		port = "9123"
	}
	if apiUrl == "" || apiKey == "" {
		log.Fatal("CDIO_API_BASE_URL and CDIO_API_KEY environment variables must be set")
		os.Exit(1)
	}

	client := cdio.NewApiClient(apiUrl, apiKey)
	registry := prometheus.NewPedanticRegistry()
	registry.MustRegister(
		promcollectors.NewProcessCollector(promcollectors.ProcessCollectorOpts{}),
		promcollectors.NewGoCollector(),
	)

	priceCollector, err := collectors.NewPriceCollector(client)
	if err != nil {
		log.Fatal(err)
	} else {
		registry.MustRegister(priceCollector)
	}

	watchCollector, err := collectors.NewWatchCollector(client)
	if err != nil {
		log.Fatal(err)
	} else {
		registry.MustRegister(watchCollector)
	}

	systemCollector, err := collectors.NewSystemCollector(client)
	if err != nil {
		log.Fatal(err)
	} else {
		registry.MustRegister(systemCollector)
	}

	http.Handle("/", promhttp.HandlerFor(registry, promhttp.HandlerOpts{
		ErrorLog: log.StandardLogger(),
	}))
	log.Info(fmt.Sprintf("Beginning to serve on port %s", port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
