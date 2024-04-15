// SPDX-FileCopyrightText: 2024 Stefan Schärmeli <schaermu@pm.me>
// SPDX-License-Identifier: MIT
package collectors

import (
	"sync"

	"github.com/schaermu/changedetection.io-exporter/pkg/cdio"
)

var (
	namespace = "changedetectionio"
	labels    = []string{"title", "source"}
)

type baseCollector struct {
	sync.RWMutex

	ApiClient *cdio.ApiClient
}

func newBaseCollector(client *cdio.ApiClient) *baseCollector {
	return &baseCollector{
		ApiClient: client,
	}
}
