# syntax=docker/dockerfile:1.7

ARG GO_VERSION=1.26.3

FROM golang:${GO_VERSION}-alpine AS builder
WORKDIR /src

RUN apk add --no-cache ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=$(go env GOARCH) \
  go build -trimpath -ldflags="-s -w" -o /out/clipbridge-server ./cmd/server

FROM alpine:3.22
WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /out/clipbridge-server /usr/local/bin/clipbridge-server
COPY configs/config.example.yaml /app/config.yaml

RUN mkdir -p /app/data

EXPOSE 8787
VOLUME ["/app/data"]

ENTRYPOINT ["/usr/local/bin/clipbridge-server"]
CMD ["-config", "/app/config.yaml"]
