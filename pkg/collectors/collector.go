// SPDX-FileCopyrightText: 2024 Stefan Schärmeli <schaermu@pm.me>
// SPDX-License-Identifier: MIT
package collectors

import (
	"sync"

	"github.com/schaermu/changedetection.io-exporter/pkg/cdio"
	"github.com/schaermu/changedetection.io-exporter/pkg/data"
)

const (
	namespace = "changedetectionio"
)

type baseCollector struct {
	sync.RWMutex

	ApiClient *cdio.ApiClient
}

type baseWatchCollector struct {
	baseCollector

	watches map[string]data.WatchItem
}

func newBaseWatchCollector(client *cdio.ApiClient, watches map[string]data.WatchItem) *baseWatchCollector {
	return &baseWatchCollector{
		baseCollector: baseCollector{
			ApiClient: client,
		},
		watches: watches,
	}
}
