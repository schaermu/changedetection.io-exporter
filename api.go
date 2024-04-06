// SPDX-FileCopyrightText: 2024 Stefan Schärmeli <schaermu@pm.me>
// SPDX-License-Identifier: MIT
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type cdioApiClient struct {
	baseUrl string
	key     string
	client  *http.Client
}

type WatchItem struct {
	LastChanged int64  `json:"last_changed"`
	LastChecked int64  `json:"last_checked"`
	LastError   bool   `json:"last_error"`
	Title       string `json:"title"`
}

type PriceData struct {
	Price        int32  `json:"price"`
	Currency     string `json:"priceCurrency"`
	Availability string `json:"availability"`
}

func newCdioApiClient(baseUrl string, key string) *cdioApiClient {
	return &cdioApiClient{
		baseUrl: fmt.Sprintf("%s/api/v1", baseUrl),
		key:     key,
		client:  &http.Client{},
	}
}

func (client *cdioApiClient) getWatches() map[string]WatchItem {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/watch", client.baseUrl), nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("x-api-key", client.key)
	res, err := client.client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	watches := make(map[string]WatchItem)
	err = json.NewDecoder(res.Body).Decode(&watches)
	if err != nil {
		panic(err)
	}
	return watches
}

func (client *cdioApiClient) getLatestPriceSnapshot(id string) (*PriceData, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/watch/%s/history/latest", client.baseUrl, id), nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("x-api-key", client.key)
	res, err := client.client.Do(req)

	if res.StatusCode == 404 {
		// watch not found, was probably removed
		return nil, fmt.Errorf("watch %s not found", id)
	}

	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	var priceData = PriceData{}
	err = json.NewDecoder(res.Body).Decode(&priceData)
	if err != nil {
		panic(err)
	}
	return &priceData, nil
}
