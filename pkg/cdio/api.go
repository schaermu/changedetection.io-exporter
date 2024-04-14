// SPDX-FileCopyrightText: 2024 Stefan Schärmeli <schaermu@pm.me>
// SPDX-License-Identifier: MIT
package cdio

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/schaermu/changedetection.io-exporter/pkg/data"
	log "github.com/sirupsen/logrus"
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
	targetUrl := fmt.Sprintf("%s/%s", client.baseUrl, url)
	log.Debugf("curl \"%s\" -H\"x-api-key:%s\"", targetUrl, client.key)
	req, err := http.NewRequest(method, targetUrl, body)
	if err != nil {
		log.Debug(err)
		return nil, err
	}
	req.Header.Add("x-api-key", client.key)
	return req, nil
}

func (client *ApiClient) GetWatches() (map[string]*data.WatchItem, error) {
	req, err := client.getRequest("GET", "watch", nil)
	if err != nil {
		return nil, err
	}

	res, err := client.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	watches := make(map[string]*data.WatchItem)
	err = json.NewDecoder(res.Body).Decode(&watches)
	if err != nil {
		return nil, err
	}
	return watches, nil
}

func (client *ApiClient) GetWatchData(id string) (*data.WatchItem, error) {
	req, err := client.getRequest("GET", fmt.Sprintf("watch/%s", id), nil)
	if err != nil {
		return nil, err
	}

	res, err := client.Client.Do(req)
	if err != nil {
		return nil, err
	}

	switch res.StatusCode {
	case 404:
		// watch not found, was probably removed
		return nil, fmt.Errorf("watch %s not found", id)
	}
	defer res.Body.Close()

	var watchItem = data.WatchItem{}
	err = json.NewDecoder(res.Body).Decode(&watchItem)
	if err != nil {
		return nil, err
	}
	return &watchItem, nil
}

func (client *ApiClient) GetLatestPriceSnapshot(id string) (*data.PriceData, error) {
	req, err := client.getRequest("GET", fmt.Sprintf("watch/%s/history/latest", id), nil)
	if err != nil {
		return nil, err
	}

	res, err := client.Client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == 404 {
		// watch not found, was probably removed
		return nil, fmt.Errorf("watch %s not found", id)
	}

	defer res.Body.Close()

	bodyText, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var priceData = data.PriceData{}
	err = json.Unmarshal(bodyText, &priceData)
	if err != nil {
		// check if the error is due to the response being an array
		if err.Error() == "json: cannot unmarshal array into Go value of type data.PriceData" {
			log.Debug("price data is an array, trying to decode as array")
			var priceDataArray []data.PriceData
			err = json.Unmarshal(bodyText, &priceDataArray)
			if err != nil {
				log.Error(err)
				return nil, err
			}
			return &priceDataArray[0], nil
		}
		log.Error(err)
		return nil, err
	}
	return &priceData, nil
}
