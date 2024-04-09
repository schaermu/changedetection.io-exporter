// SPDX-FileCopyrightText: 2024 Stefan Schärmeli <schaermu@pm.me>
// SPDX-License-Identifier: MIT
package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func CreateSimpleTestApiServer(t *testing.T, url string, payloadFile string) *httptest.Server {
	return CreateTestApiServer(t, map[string]string{url: payloadFile})
}

func CreateTestApiServer(t *testing.T, urlPayloadMap map[string]string) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		payloadFile, ok := urlPayloadMap[req.URL.String()]
		if !ok {
			// return 404 for unknown requests
			rw.WriteHeader(http.StatusNotFound)
		} else {
			if payloadFile == "" {
				rw.Write([]byte("OK"))
			} else {
				expected, err := os.ReadFile(payloadFile)
				if err != nil {
					rw.WriteHeader(http.StatusNotFound)
				} else {
					rw.Write(expected)
				}
			}

		}
	}))
	return server
}

func TestGetRequestApiKey(t *testing.T) {
	server := CreateSimpleTestApiServer(t, "/api/v1/watch", "")
	defer server.Close()

	api := NewApiClient(server.URL, "foo-bar-key")
	request, err := api.getRequest("GET", "/watch", nil)
	ok(t, err)
	equals(t, "foo-bar-key", request.Header.Get("x-api-key"))
}

func TestGetWatches(t *testing.T) {
	// Start a local HTTP server
	server := CreateSimpleTestApiServer(t, "/api/v1/watch", "./test/json/getWatches.json")
	defer server.Close()

	api := NewApiClient(server.URL, "foo-bar-key")
	watches, err := api.getWatches()
	ok(t, err)
	equals(t, 2, len(watches))
	equals(t, "Random Quote", watches["6a4b7d5c-fee4-4616-9f43-4ac97046b595"].Title)
}

func TestGetLatestPriceSnapshot(t *testing.T) {
	id := "6a4b7d5c-fee4-4616-9f43-4ac97046b595"

	// Start a local HTTP server
	server := CreateSimpleTestApiServer(t, fmt.Sprintf("/api/v1/watch/%s/history/latest", id), fmt.Sprintf("./test/json/getLatestPriceSnapshot_%s.json", id))
	defer server.Close()

	api := NewApiClient(server.URL, "foo-bar-key")
	priceData, err := api.getLatestPriceSnapshot("6a4b7d5c-fee4-4616-9f43-4ac97046b595")
	ok(t, err)
	equals(t, float64(100), priceData.Price)
	equals(t, "USD", priceData.Currency)
}

func TestSetBaseUrl(t *testing.T) {
	api := NewApiClient("http://localhost:8080", "foo-bar-key")
	api.SetBaseUrl("http://localhost:8081")
	equals(t, "http://localhost:8081/api/v1", api.baseUrl)
}

func TestGetLastestPriceSnapshot_NotFound(t *testing.T) {
	existingId := "6a4b7d5c-fee4-4616-9f43-4ac97046b595"
	lookupId := "foo-bar-id"

	// Start a local HTTP server
	server := CreateSimpleTestApiServer(t, fmt.Sprintf("/api/v1/watch/%s/history/latest", existingId), "")
	defer server.Close()

	api := NewApiClient(server.URL, "foo-bar-key")
	priceData, err := api.getLatestPriceSnapshot(lookupId)
	equals(t, fmt.Errorf("watch %s not found", lookupId), err)
	equals(t, (*PriceData)(nil), priceData)
}
