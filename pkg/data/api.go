// SPDX-FileCopyrightText: 2024 Stefan Schärmeli <schaermu@pm.me>
// SPDX-License-Identifier: MIT
package data

import (
	"fmt"
	"net/url"
)

type StringBoolean bool

func (sb *StringBoolean) UnmarshalJSON(data []byte) error {
	if string(data) == "false" {
		*sb = false
	} else {
		*sb = true
	}
	return nil
}

type WatchItem struct {
	LastChanged            int64         `json:"last_changed"`
	LastChecked            int64         `json:"last_checked"`
	LastError              StringBoolean `json:"last_error"`
	Title                  string        `json:"title"`
	Url                    string        `json:"url"`
	CheckCount             int           `json:"check_count,omitempty"`
	FetchTime              float64       `json:"fetch_time,omitempty"`
	NotificationAlertCount int           `json:"notification_alert_count,omitempty"`
	LastCheckStatus        int           `json:"last_check_status,omitempty"`
	PriceData              *PriceData    `json:"price,omitempty"`
}

type PriceData struct {
	Price        float64 `json:"price"`
	Currency     string  `json:"priceCurrency"`
	Availability string  `json:"availability"`
}

type SystemInfo struct {
	Version        string   `json:"version"`
	Uptime         float64  `json:"uptime"`
	WatchCount     int      `json:"watch_count"`
	OverdueWatches []string `json:"overdue_watches"`
	QueueSize      int      `json:"queue_size"`
}

func (w *WatchItem) GetMetrics() ([]string, error) {
	url, err := url.ParseRequestURI(w.Url)
	if err != nil {
		return nil, err
	} else if url.Host == "" {
		return nil, fmt.Errorf("host is empty")
	}

	if w.Title == "" {
		return nil, fmt.Errorf("title is empty")
	}
	return []string{w.Title, url.Host}, nil
}
