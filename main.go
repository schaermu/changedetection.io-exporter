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

	log "github.com/sirupsen/logrus"
)

func main() {
	var (
		port   = os.Getenv("PORT")
		apiUrl = os.Getenv("CDIO_API_BASE_URL")
		apiKey = os.Getenv("CDIO_API_KEY")
	)

	if port == "" {
		port = "8123"
	}
	if apiUrl == "" || apiKey == "" {
		log.Fatal("CDIO_API_BASE_URL and CDIO_API_KEY environment variables must be set")
		os.Exit(1)
	}

	collector, err := NewPriceCollector(apiUrl, apiKey)
	if err != nil {
		log.Fatal(err)
	}
	prometheus.MustRegister(collector)
	http.Handle("/", promhttp.Handler())
	log.Info(fmt.Sprintf("Beginning to serve on port %s", port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}
