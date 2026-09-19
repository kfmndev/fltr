ARG GO_VERSION=1.25

FROM golang:${GO_VERSION}-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/fltr .

FROM ghcr.io/linuxserver/baseimage-alpine:3.24

COPY --from=build /out/fltr /app/fltr/fltr
RUN \
  echo "**** set permissions ****" && \
  chmod +x /app/fltr/fltr && \
  lsiown abc:abc /app/fltr

COPY docker/root/ /

ENV FLTR_BLOCK_RULES_FILE=/config/block_rules.json \
    FLTR_ADDR=:8080

EXPOSE 8080

# Probe the reserved Liveness Endpoint; host and port follow FLTR_ADDR, falling
# back to loopback when it binds all interfaces.
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
  CMD HOST="${FLTR_ADDR%:*}"; PORT="${FLTR_ADDR##*:}"; \
      case "$HOST" in ""|0.0.0.0|"::"|"[::]") HOST=127.0.0.1;; esac; \
      wget -q -O /dev/null "http://${HOST}:${PORT}/healthz"
