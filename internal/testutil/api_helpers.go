package testutil

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strconv"
	"testing"

	log "github.com/sirupsen/logrus"

	"github.com/google/uuid"
	"github.com/schaermu/changedetection.io-exporter/pkg/data"
)

type ApiTestServer struct {
	Server  *httptest.Server
	watches map[string]*data.WatchItem
}

func (s *ApiTestServer) URL() string {
	return s.Server.URL
}

func (s *ApiTestServer) Close() {
	s.Server.Close()
}

func NewWatchDb(numItems int) map[string]*data.WatchItem {
	ret := make(map[string]*data.WatchItem)
	for i := 0; i < numItems; i++ {
		uuid, watch := NewTestItem(fmt.Sprintf("testwatch #%s", strconv.Itoa(i+1)), rand.Float64(), "USD")
		ret[uuid] = watch
	}
	return ret
}

func NewTestItem(title string, price float64, currency string) (string, *data.WatchItem) {
	return uuid.New().String(), &data.WatchItem{
		Title: title,
		PriceData: &data.PriceData{
			Price:        price,
			Currency:     currency,
			Availability: "InStock",
		},
	}
}

func writeJson(rw http.ResponseWriter, v any) {
	if res, err := json.Marshal(v); err == nil {
		rw.Header().Set("Content-Type", "application/json")
		rw.Write(res)
	} else {
		rw.WriteHeader(http.StatusInternalServerError)
	}
}

func CreateTestApiServer(t *testing.T, watches map[string]*data.WatchItem) *ApiTestServer {
	var watchDetailPattern = regexp.MustCompile(`^\/api\/v1\/watch\/(?P<UUID>[0-9a-fA-F]{8}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{12})\/?(?P<ACTION>.+)?`)
	return &ApiTestServer{
		watches: watches,
		Server: httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			if req.URL.Path == "/api/v1/watch" {
				writeJson(rw, watches)
			} else if watchDetailPattern.MatchString(req.URL.Path) {
				// get UUID from path
				matches := watchDetailPattern.FindStringSubmatch(req.URL.Path)
				uuidIdx := watchDetailPattern.SubexpIndex("UUID")
				uuid := matches[uuidIdx]

				// find watch
				watch, ok := watches[uuid]
				if !ok {
					log.Infof("could not find watch with id %s", uuid)
					rw.WriteHeader(http.StatusNotFound)
				} else {
					actionIndex := watchDetailPattern.SubexpIndex("ACTION")
					if actionIndex > -1 {
						switch matches[actionIndex] {
						case "history/latest":
							// return price data
							writeJson(rw, watch.PriceData)
						default:
							// return details
							writeJson(rw, watch)
						}
					} else {
						// return details
						writeJson(rw, watch)
					}
				}
			} else {
				log.Infof("could not map path %s", req.URL.Path)
				rw.WriteHeader((http.StatusNotFound))
			}
		})),
	}
}
