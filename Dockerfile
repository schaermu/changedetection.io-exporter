ARG GO_VERSION=1.22.2

# build exporter
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine as build
RUN apk update && apk add --no-cache git

# install dependencies
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

# build statically linked binary
COPY . .
ARG TARGETOS TARGETARCH
RUN --mount=target=. \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -ldflags="-w -s" -installsuffix 'static' -o /changedetectionio_exporter .

# build runtime image
FROM --platform=$TARGETPLATFORM gcr.io/distroless/static-debian12
USER nonroot:nonroot
COPY --from=build --chown=nonroot:nonroot /changedetectionio_exporter /changedetectionio_exporter
EXPOSE 9123
ENTRYPOINT ["/changedetectionio_exporter"]