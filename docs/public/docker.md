---
icon: lucide/container
---

# Docker

Pre-built images are published to the GitHub Container Registry as `ghcr.io/kfmndev/fltr`. They are built for `linux/amd64` and `linux/arm64`.

## Docker Run

```docker
docker run -d \
  --name fltr \
  -p 8080:8080 \
  -v "$PWD/block_rules.json":/config/block_rules.json \
  -e FLTR_ALLOW_UPSTREAM=http://allow.example.com \
  ghcr.io/kfmndev/fltr
```

!!! tip

    The images are built from [`linuxserver/docker-baseimage-alpine`](https://github.com/linuxserver/docker-baseimage-alpine/) (`ghcr.io/linuxserver/baseimage-alpine`). Set `PUID` and `PGID` to avoid file permission issues, see [Understanding PUID and PGID](https://docs.linuxserver.io/general/understanding-puid-and-pgid/) for details.

## Docker Compose

```yml title="docker-compose.yml"
services:
  fltr:
    image: ghcr.io/kfmndev/fltr
    container_name: fltr
    ports:
      - 8080:8080
    environment:
      - FLTR_ALLOW_UPSTREAM=http://allow.example.com
    volumes:
      - ./block_rules.json:/config/block_rules.json:ro
    restart: unless-stopped
```

## Healthcheck

The image ships a `HEALTHCHECK` that probes `http://<host>:<port>/healthz` every 30 seconds. Both come from `FLTR_ADDR`: the host is everything before the last `:`, and the port is extracted with `${FLTR_ADDR##*:}`, which strips everything up to and including that colon, so `:8080` becomes port `8080` and `0.0.0.0:9090` becomes `9090`. When `FLTR_ADDR` binds all interfaces (`:8080`, `0.0.0.0:8080`, `[::]:8080`) the probe falls back to `127.0.0.1`. Change `FLTR_ADDR` and the probe follows, so the container reports healthy as long as fltr answers on the address it was told to serve.

```sh
docker inspect --format '{{.State.Health.Status}}' fltr
```

## Image tags

The image is published under several tags, depending on what triggered the build:

- Every push to `main`: `ghcr.io/kfmndev/fltr:main` and `ghcr.io/kfmndev/fltr:sha-<short-sha>`, where `<short-sha>` is the commit's short hash.
- Every **git** version tag (e.g. `v1.0.0`): `latest`, plus the version split into `1`, `1.0`, and `1.0.0`.

Use `latest` when you want to track the most recent release, or pin a specific version to make sure nothing changes.
