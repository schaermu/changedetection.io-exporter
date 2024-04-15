// SPDX-FileCopyrightText: 2024 Stefan Schärmeli <schaermu@pm.me>
// SPDX-License-Identifier: MIT
package data

import (
	"net/url"
)

type WatchItem struct {
	LastChanged            int64      `json:"last_changed"`
	LastChecked            int64      `json:"last_checked"`
	LastError              bool       `json:"last_error"`
	Title                  string     `json:"title"`
	Url                    string     `json:"url"`
	CheckCount             int        `json:"check_count,omitempty"`
	FetchTime              float64    `json:"fetch_time,omitempty"`
	NotificationAlertCount int        `json:"notification_alert_count,omitempty"`
	LastCheckStatus        int        `json:"last_check_status,omitempty"`
	PriceData              *PriceData `json:"price,omitempty"`
}

func (w *WatchItem) GetMetrics() ([]string, error) {
	url, err := url.Parse(w.Url)
	if err != nil {
		return nil, err
	}
	return []string{w.Title, url.Host}, nil
}

type PriceData struct {
	Price        float64 `json:"price"`
	Currency     string  `json:"priceCurrency"`
	Availability string  `json:"availability"`
}
