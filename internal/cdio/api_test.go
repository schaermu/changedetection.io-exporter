// SPDX-FileCopyrightText: 2024 Stefan Schärmeli <schaermu@pm.me>
// SPDX-License-Identifier: MIT
package cdio

import (
	"fmt"
	"testing"

	"github.com/schaermu/changedetection.io-exporter/internal/data"
	"github.com/schaermu/changedetection.io-exporter/internal/testutil"
)

func TestGetRequestApiKey(t *testing.T) {
	watchDb := testutil.NewWatchDb(2)
	server := testutil.CreateTestApiServer(t, watchDb)
	defer server.Close()

	api := NewApiClient(server.URL(), "foo-bar-key")
	request, err := api.getRequest("GET", "/watch", nil)
	testutil.Ok(t, err)
	testutil.Equals(t, "foo-bar-key", request.Header.Get("x-api-key"))
}

func TestGetWatches(t *testing.T) {
	watchDb := testutil.NewWatchDb(1)
	uuid, watch := testutil.NewTestItem("Test Me", 100, "USD")
	watchDb[uuid] = watch
	server := testutil.CreateTestApiServer(t, watchDb)
	defer server.Close()

	api := NewApiClient(server.URL(), "foo-bar-key")
	watches, err := api.GetWatches()
	testutil.Ok(t, err)
	testutil.Equals(t, 2, len(watches))
	testutil.Equals(t, "Test Me", watches[uuid].Title)
}

func TestGetLatestPriceSnapshot(t *testing.T) {
	watchDb := testutil.NewWatchDb(1)
	uuid, watchItem := testutil.NewTestItem("Test Me", 100, "USD")
	watchDb[uuid] = watchItem
	server := testutil.CreateTestApiServer(t, watchDb)
	defer server.Close()

	api := NewApiClient(server.URL(), "foo-bar-key")
	priceData, err := api.GetLatestPriceSnapshot(uuid)
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
	watchDb := testutil.NewWatchDb(2)
	server := testutil.CreateTestApiServer(t, watchDb)
	defer server.Close()

	nonExistingId := "i-surely-do-not-exist"

	api := NewApiClient(server.URL(), "foo-bar-key")
	priceData, err := api.GetLatestPriceSnapshot(nonExistingId)
	testutil.Equals(t, fmt.Errorf("watch %s not found", nonExistingId), err)
	testutil.Equals(t, (*data.PriceData)(nil), priceData)
}
