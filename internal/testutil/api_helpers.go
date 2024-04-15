package testutil

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/schaermu/changedetection.io-exporter/pkg/data"
)

type ApiTestServerOptions struct {
	fmt.Stringer
	PricesAsArray bool
	SystemInfo    *data.SystemInfo
}
type ApiTestServerOption func(*ApiTestServerOptions)

func (o ApiTestServerOptions) String() string {
	return fmt.Sprintf("ApiTestServerOptions{PricesAsArray: %t, SystemInfo: %v}", o.PricesAsArray, o.SystemInfo)
}

func WithPricesAsArray() ApiTestServerOption {
	return func(o *ApiTestServerOptions) {
		o.PricesAsArray = true
	}
}

func WithSystemInfo(info *data.SystemInfo) ApiTestServerOption {
	return func(o *ApiTestServerOptions) {
		o.SystemInfo = info
	}
}

type ApiTestServer struct {
	Server  *httptest.Server
	Options ApiTestServerOptions
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
		uuid, watch := NewTestItem(fmt.Sprintf("testwatch-%s", strconv.Itoa(i+1)), rand.Float64(), "USD", 20, 15, 10)
		ret[uuid] = watch
	}
	return ret
}

func NewTestItem(title string, price float64, currency string, checkCount int, fetchTime float64, alertCount int) (string, *data.WatchItem) {
	return uuid.New().String(), &data.WatchItem{
		Title:                  title,
		Url:                    fmt.Sprintf("https://www.%s.org/", strings.ReplaceAll(strings.ToLower(title), " ", "-")),
		CheckCount:             checkCount,
		FetchTime:              fetchTime,
		NotificationAlertCount: alertCount,
		LastCheckStatus:        200,
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
		if _, err := rw.Write(res); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		rw.WriteHeader(http.StatusInternalServerError)
	}
}

func CreateTestApiServer(t *testing.T, watches map[string]*data.WatchItem, options ...ApiTestServerOption) *ApiTestServer {
	var watchDetailPattern = regexp.MustCompile(`^\/api\/v1\/watch\/(?P<UUID>[0-9a-fA-F]{8}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{4}\b-[0-9a-fA-F]{12})\/?(?P<ACTION>.+)?`)

	// pull in options
	opts := ApiTestServerOptions{
		PricesAsArray: false,
		SystemInfo:    &data.SystemInfo{Version: "1.0.0", Uptime: 100, WatchCount: len(watches), OverdueWatches: []string{}, QueueSize: 0},
	}
	for _, o := range options {
		o(&opts)
	}
	t.Log("pulled in options", opts)

	return &ApiTestServer{
		watches: watches,
		Server: httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			if req.URL.Path == "/api/v1/watch" {
				writeJson(rw, watches)
			} else if req.URL.Path == "/api/v1/systeminfo" {
				writeJson(rw, opts.SystemInfo)
			} else if watchDetailPattern.MatchString(req.URL.Path) {
				// get UUID from path
				matches := watchDetailPattern.FindStringSubmatch(req.URL.Path)
				uuidIdx := watchDetailPattern.SubexpIndex("UUID")
				uuid := matches[uuidIdx]

				// find watch
				watch, ok := watches[uuid]
				if !ok {
					t.Logf("could not find watch with id %s", uuid)
					rw.WriteHeader(http.StatusNotFound)
				} else {
					actionIndex := watchDetailPattern.SubexpIndex("ACTION")
					if actionIndex > -1 {
						switch matches[actionIndex] {
						case "history/latest":
							// return price data
							if opts.PricesAsArray {
								writeJson(rw, []data.PriceData{*watch.PriceData})
							} else {
								writeJson(rw, watch.PriceData)
							}
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
				t.Logf("could not map path %s", req.URL.Path)
				rw.WriteHeader((http.StatusNotFound))
			}
		})),
	}
}
