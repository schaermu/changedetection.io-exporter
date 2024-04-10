package testutil

import (
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
