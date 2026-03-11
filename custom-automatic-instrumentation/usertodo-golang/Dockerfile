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

FROM alpine:3.20.6

RUN apk --no-cache add ca-certificates \
    && addgroup -S appgroup \
    && adduser -S -D -H -h /app -s /sbin/nologin -G appgroup appuser

COPY --from=builder /bin/usertodo-api /usr/local/bin/usertodo-api

WORKDIR /app

USER appuser:appgroup

EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/usertodo-api"]
