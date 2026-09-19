---
icon: lucide/download
---

# Installation

fltr ships as a pre-built binary, a Docker image, and source. All paths need the same two things: a JSON file containing the Block Rules and an HTTP service to receive Allowed requests.

## Pre-built binary

Download `fltr_<version>_linux_amd64` (or `linux_arm64`) from a [GitHub release](https://github.com/kfmndev/fltr/releases) and run it directly. Each release includes a `checksums.txt` for verification.

## Docker

Images are published to `ghcr.io/kfmndev/fltr` for `linux/amd64` and `linux/arm64`:

```sh
docker run -d \
  --name fltr \
  -p 8080:8080 \
  -v "$PWD/block_rules.json":/config/block_rules.json \
  -e FLTR_ALLOW_UPSTREAM=http://allow.example.com \
  ghcr.io/kfmndev/fltr
```

See [Docker](../guide/docker.md) for compose setup, image tags, and the bundled healthcheck.

## From source

Requires [Go 1.25](https://go.dev) or newer:

```sh
go build ./...
./fltr
```

Continue with the [Quickstart](index.md) to configure a Block Rule and send your first request.
