BINARY_NAME=changedetectionio_exporter
GOCOVER=go tool cover
GOTESTSUM=go run gotest.tools/gotestsum@latest

.DEFAULT_GOAL := all
.PHONY: clean test watch cover run start

all: clean test build

docker:
	docker build -t ghcr.io/schaermu/changedetection.io-exporter:latest .

build:
	GOARCH=amd64 GOOS=linux go build -o ./build/${BINARY_NAME} .

run:
	./build/${BINARY_NAME}

start: clean build run

clean:
	go clean
	go clean -testcache
	rm -rf ./build

test:
	$(GOTESTSUM) -f testname -- -tags=test -coverprofile=coverage.txt -race -covermode=atomic ./...

watch:
	$(GOTESTSUM) --watch -f testname -- -tags=test -coverprofile=coverage.txt -race -covermode=atomic ./...

cover:
	$(GOTESTSUM) -f testname -- -tags=test ./... -coverprofile=coverage.out
	$(GOCOVER) -func=coverage.out
	$(GOCOVER) -html=coverage.out -o coverage.html