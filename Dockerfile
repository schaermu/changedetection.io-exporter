ARG GO_VERSION=1.22.2

# build exporter
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine as build
RUN apk update && apk add --no-cache git

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# build static binary to run on distroless/static
ARG TARGETOS TARGETARCH
RUN --mount=target=. \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -ldflags="-w -s" -installsuffix 'static' -o /changedetectionio_exporter .

# build runtime image
FROM --platform=$TARGETPLATFORM gcr.io/distroless/static-debian12

LABEL maintainer="schaermu"

USER nonroot:nonroot
COPY --from=build --chown=nonroot:nonroot /changedetectionio_exporter /changedetectionio_exporter

EXPOSE 8080

ENTRYPOINT ["/changedetectionio_exporter"]