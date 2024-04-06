ARG GO_VERSION=1.22.2

# build exporter
FROM golang:${GO_VERSION}-alpine as go-builder

# git is required to install using go mod
RUN apk update && apk add --no-cache git

WORKDIR /src

# download go dependencies
COPY go.mod go.sum ./
RUN go mod download

COPY . ./

# build static binary to run on distroless/static
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -installsuffix 'static' -o /changedetectionio_exporter .

# build runtime image
FROM gcr.io/distroless/static

LABEL maintainer="schaermu"

USER nonroot:nonroot
COPY --from=go-builder --chown=nonroot:nonroot /changedetectionio_exporter /changedetectionio_exporter

EXPOSE 8080

ENTRYPOINT ["/changedetectionio_exporter"]