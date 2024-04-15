// SPDX-FileCopyrightText: 2024 Stefan Schärmeli <schaermu@pm.me>
// SPDX-License-Identifier: MIT
package cdio

import (
	"fmt"
	"testing"

	"github.com/schaermu/changedetection.io-exporter/internal/testutil"
	"github.com/schaermu/changedetection.io-exporter/pkg/data"
)

func TestGetRequestApiKey(t *testing.T) {
	watchDb := testutil.NewWatchDb(2)
	server := testutil.CreateTestApiServer(t, watchDb)
	defer server.Close()

	apiKey := "manual-api-key"

	api := NewApiClient(server.URL(), apiKey)
	request, err := api.getRequest("GET", "/watch", nil)
	testutil.Ok(t, err)
	testutil.Equals(t, apiKey, request.Header.Get("x-api-key"))
}

func TestGetWatches(t *testing.T) {
	watchDb := testutil.NewWatchDb(1)
	uuid, watch := testutil.NewTestItem("Test Me", 100, "USD", 20, 15, 10)
	watchDb[uuid] = watch
	server := testutil.CreateTestApiServer(t, watchDb)
	defer server.Close()

	api := NewTestApiClient(server.URL())
	watches, err := api.GetWatches()
	testutil.Ok(t, err)
	testutil.Equals(t, 2, len(watches))
	testutil.Equals(t, "Test Me", watches[uuid].Title)
}

func TestGetWatchData(t *testing.T) {
	watchDb := testutil.NewWatchDb(1)
	uuid, watch := testutil.NewTestItem("Test Me", 100, "USD", 20, 15, 10)
	watch.CheckCount = 2
	watchDb[uuid] = watch
	server := testutil.CreateTestApiServer(t, watchDb)
	defer server.Close()

	api := NewTestApiClient(server.URL())
	watchItem, err := api.GetWatchData(uuid)

	testutil.Ok(t, err)
	testutil.Equals(t, "Test Me", watchItem.Title)
	testutil.Equals(t, 2, watchItem.CheckCount)
}

func TestGetWatchData_NotFound(t *testing.T) {
	watchDb := testutil.NewWatchDb(2)
	server := testutil.CreateTestApiServer(t, watchDb)
	defer server.Close()

	nonExistingId := "i-surely-do-not-exist"

	api := NewTestApiClient(server.URL())
	watchItem, err := api.GetWatchData(nonExistingId)

	testutil.Equals(t, fmt.Errorf("watch %s not found", nonExistingId), err)
	testutil.Equals(t, (*data.WatchItem)(nil), watchItem)
}

func TestGetLatestPriceSnapshot(t *testing.T) {
	watchDb := testutil.NewWatchDb(1)
	uuid, watchItem := testutil.NewTestItem("Test Me", 100, "USD", 20, 15, 10)
	watchDb[uuid] = watchItem
	server := testutil.CreateTestApiServer(t, watchDb)
	defer server.Close()

	api := NewTestApiClient(server.URL())
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

	api := NewTestApiClient(server.URL())
	priceData, err := api.GetLatestPriceSnapshot(nonExistingId)
	testutil.Equals(t, fmt.Errorf("watch %s not found", nonExistingId), err)
	testutil.Equals(t, (*data.PriceData)(nil), priceData)
}
