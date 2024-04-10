// SPDX-FileCopyrightText: 2024 Stefan Schärmeli <schaermu@pm.me>
// SPDX-License-Identifier: MIT
package cdio

import (
	"fmt"
	"testing"

	"github.com/schaermu/changedetection.io-exporter/internal/testutil"
)

func TestGetRequestApiKey(t *testing.T) {
	server := testutil.CreateSimpleTestApiServer(t, "/api/v1/watch", "")
	defer server.Close()

	api := NewApiClient(server.URL, "foo-bar-key")
	request, err := api.getRequest("GET", "/watch", nil)
	testutil.Ok(t, err)
	testutil.Equals(t, "foo-bar-key", request.Header.Get("x-api-key"))
}

func TestGetWatches(t *testing.T) {
	// Start a local HTTP server
	server := testutil.CreateSimpleTestApiServer(t, "/api/v1/watch", testutil.GetFixturePath("json/getWatches.json"))
	defer server.Close()

	api := NewApiClient(server.URL, "foo-bar-key")
	watches, err := api.GetWatches()
	testutil.Ok(t, err)
	testutil.Equals(t, 2, len(watches))
	testutil.Equals(t, "Random Quote", watches["6a4b7d5c-fee4-4616-9f43-4ac97046b595"].Title)
}

func TestGetLatestPriceSnapshot(t *testing.T) {
	id := "6a4b7d5c-fee4-4616-9f43-4ac97046b595"

	// Start a local HTTP server
	server := testutil.CreateSimpleTestApiServer(t, fmt.Sprintf("/api/v1/watch/%s/history/latest", id), testutil.GetFixturePath("json/getLatestPriceSnapshot_4ac97046b595.json"))
	defer server.Close()

	api := NewApiClient(server.URL, "foo-bar-key")
	priceData, err := api.GetLatestPriceSnapshot("6a4b7d5c-fee4-4616-9f43-4ac97046b595")
	testutil.Ok(t, err)
	testutil.Equals(t, float64(100), priceData.Price)
	testutil.Equals(t, "USD", priceData.Currency)
}

func TestSetBaseUrl(t *testing.T) {
	api := NewApiClient("http://localhost:8080", "foo-bar-key")
	api.SetBaseUrl("http://localhost:8081")
	testutil.Equals(t, "http://localhost:8081/api/v1", api.baseUrl)
}

func TestGetLastestPriceSnapshot_NotFound(t *testing.T) {
	existingId := "6a4b7d5c-fee4-4616-9f43-4ac97046b595"
	lookupId := "foo-bar-id"

	// Start a local HTTP server
	server := testutil.CreateSimpleTestApiServer(t, fmt.Sprintf("/api/v1/watch/%s/history/latest", existingId), "")
	defer server.Close()

	api := NewApiClient(server.URL, "foo-bar-key")
	priceData, err := api.GetLatestPriceSnapshot(lookupId)
	testutil.Equals(t, fmt.Errorf("watch %s not found", lookupId), err)
	testutil.Equals(t, (*PriceData)(nil), priceData)
}
