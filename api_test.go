package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func CreateTestApiServer(t *testing.T, url string, payloadFile string) *httptest.Server {
	var expected []byte
	var err error
	if payloadFile == "" {
		expected = []byte("OK")
	} else {
		expected, err = os.ReadFile(payloadFile)
		if err != nil {
			t.Fatal(err)
		}
	}

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		// Test request parameters
		equals(t, req.URL.String(), url)
		// Send response to be tested
		rw.Write(expected)
	}))
	return server
}

func TestGetRequestApiKey(t *testing.T) {
	server := CreateTestApiServer(t, "/api/v1/watch", "")
	defer server.Close()

	api := NewApiClient(server.URL, "foo-bar-key")
	request, err := api.getRequest("GET", "/watch", nil)
	ok(t, err)
	equals(t, "foo-bar-key", request.Header.Get("x-api-key"))
}

func TestGetWatches(t *testing.T) {
	// Start a local HTTP server
	server := CreateTestApiServer(t, "/api/v1/watch", "./test/json/getWatches.json")
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
	server := CreateTestApiServer(t, fmt.Sprintf("/api/v1/watch/%s/history/latest", id), fmt.Sprintf("./test/json/getLatestPriceSnapshot_%s.json", id))
	defer server.Close()

	api := NewApiClient(server.URL, "foo-bar-key")
	priceData, err := api.getLatestPriceSnapshot("6a4b7d5c-fee4-4616-9f43-4ac97046b595")
	ok(t, err)
	equals(t, int32(100), priceData.Price)
	equals(t, "USD", priceData.Currency)
}
