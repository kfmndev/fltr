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

ENV FLTR_BLOCKED_FILE=/config/block_rules.json

EXPOSE 8080

LABEL org.opencontainers.image.source="https://github.com/kfmndev/fltr"
