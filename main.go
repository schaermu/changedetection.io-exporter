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

var (
	port     = os.Getenv("PORT")
	logLevel = os.Getenv("LOG_LEVEL")
	apiUrl   = os.Getenv("CDIO_API_BASE_URL")
	apiKey   = os.Getenv("CDIO_API_KEY")
)

func init() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05.000000",
	})

	switch logLevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "fatal":
		log.SetLevel(log.FatalLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}
}

func main() {
	if port == "" {
		port = "9123"
	}
	if apiUrl == "" || apiKey == "" {
		log.Fatal("CDIO_API_BASE_URL and CDIO_API_KEY environment variables must be set")
		os.Exit(1)
	}

	client := cdio.NewApiClient(apiUrl, apiKey)
	registry := prometheus.NewPedanticRegistry()

	// register default collectors
	registry.MustRegister(
		promcollectors.NewProcessCollector(promcollectors.ProcessCollectorOpts{}),
		promcollectors.NewGoCollector(),
	)

	// register changedetection.io collectors
	registry.MustRegister(
		collectors.NewSystemCollector(client),
		collectors.NewWatchCollector(client),
		collectors.NewPriceCollector(client),
	)

	// register prometheus handler
	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{
		ErrorLog: log.StandardLogger(),
	}))
	log.Info(fmt.Sprintf("Beginning to serve on port %s", port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
