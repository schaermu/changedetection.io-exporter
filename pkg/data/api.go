// SPDX-FileCopyrightText: 2024 Stefan Schärmeli <schaermu@pm.me>
// SPDX-License-Identifier: MIT
package data

type WatchItem struct {
	LastChanged            int64      `json:"last_changed"`
	LastChecked            int64      `json:"last_checked"`
	LastError              bool       `json:"last_error"`
	Title                  string     `json:"title"`
	CheckCount             int        `json:"check_count,omitempty"`
	FetchTime              float64    `json:"fetch_time,omitempty"`
	NotificationAlertCount int        `json:"notification_alert_count,omitempty"`
	LastCheckStatus        int        `json:"last_check_status,omitempty"`
	PriceData              *PriceData `json:"price,omitempty"`
}

type PriceData struct {
	Price        float64 `json:"price"`
	Currency     string  `json:"priceCurrency"`
	Availability string  `json:"availability"`
}
