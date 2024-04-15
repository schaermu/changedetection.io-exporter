BINARY_NAME=changedetectionio_exporter
GOCOVER=go tool cover
GOTESTSUM=go run gotest.tools/gotestsum@latest

.DEFAULT_GOAL := all
.PHONY: clean test watch cover run start

all: clean test build

docker:
	docker buildx build --platform linux/amd64,linux/arm64,linux/arm/v7 -t ghcr.io/schaermu/changedetection.io-exporter:latest .

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
	$(GOTESTSUM) -f standard-verbose -- -tags=test -coverprofile=coverage.txt -race -covermode=atomic ./pkg/...

watch:
	$(GOTESTSUM) --watch -f testname -- -tags=test -coverprofile=coverage.txt -race -covermode=atomic ./...

cover:
	$(GOTESTSUM) -f pkgname -- -tags=test -coverprofile=coverage.out -race -covermode=atomic ./pkg/...
	$(GOCOVER) -func=coverage.out
	$(GOCOVER) -html=coverage.out -o coverage.html