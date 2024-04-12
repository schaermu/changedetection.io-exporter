// SPDX-FileCopyrightText: 2024 Stefan Schärmeli <schaermu@pm.me>
// SPDX-License-Identifier: MIT
package cdio

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/schaermu/changedetection.io-exporter/internal/data"
)

type ApiClient struct {
	Client  *http.Client
	baseUrl string
	key     string
}

func NewApiClient(baseUrl string, key string) *ApiClient {
	return &ApiClient{
		Client:  &http.Client{},
		baseUrl: fmt.Sprintf("%s/api/v1", baseUrl),
		key:     key,
	}
}

func (client *ApiClient) SetBaseUrl(baseUrl string) {
	client.baseUrl = fmt.Sprintf("%s/api/v1", baseUrl)
}

func (client *ApiClient) getRequest(method string, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, fmt.Sprintf("%s/%s", client.baseUrl, url), body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("x-api-key", client.key)
	return req, nil
}

func (client *ApiClient) GetWatches() (map[string]data.WatchItem, error) {
	req, err := client.getRequest("GET", "watch", nil)
	if err != nil {
		return nil, err
	}

	res, err := client.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	watches := make(map[string]data.WatchItem)
	err = json.NewDecoder(res.Body).Decode(&watches)
	if err != nil {
		return nil, err
	}
	return watches, nil
}

func (client *ApiClient) GetLatestPriceSnapshot(id string) (*data.PriceData, error) {
	req, err := client.getRequest("GET", fmt.Sprintf("watch/%s/history/latest", id), nil)
	if err != nil {
		return nil, err
	}
	res, err := client.Client.Do(req)

	if res.StatusCode == 404 {
		// watch not found, was probably removed
		return nil, fmt.Errorf("watch %s not found", id)
	}

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var priceData = data.PriceData{}
	err = json.NewDecoder(res.Body).Decode(&priceData)
	if err != nil {
		return nil, fmt.Errorf("error while decoding price data: %v", err)
	}
	return &priceData, nil
}
