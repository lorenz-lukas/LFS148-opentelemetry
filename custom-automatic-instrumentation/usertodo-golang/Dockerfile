## syntax=docker/dockerfile:1.7
FROM golang:1.22.12-alpine3.20 AS builder

WORKDIR /app

COPY go.mod ./
COPY cmd ./cmd
COPY internal ./internal

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux \
    go build -mod=mod -trimpath -ldflags="-s -w" -o /bin/usertodo-api ./cmd/api

FROM scratch

COPY --from=builder /bin/usertodo-api /usr/local/bin/usertodo-api
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

WORKDIR /app

EXPOSE 8081

ENTRYPOINT ["/usr/local/bin/usertodo-api"]
