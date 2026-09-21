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

!!! info
    [`ghcr.io/linuxserver/baseimage-alpine`](https://github.com/linuxserver/docker-baseimage-alpine/pkgs/container/baseimage-alpine) is used as a baseimage. Accordingly, set `PUID` and `PGID` to avoid file permission issues. See [Understanding PUID and PGID](https://docs.linuxserver.io/general/understanding-puid-and-pgid/) for details.

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

The image ships a `HEALTHCHECK` that probes `http://${HOST}:${PORT}/healthz` every 30 seconds.

`FLTR_ADDR` provides `${HOST}` and `${PORT}` through POSIX shell parameter expansion (supported by `linuxserver/baseimage-alpine` which uses BusyBox `ash`), realising a split at the last colon. If the address binds all interfaces (`:8080`, `0.0.0.0:8080`, `[::]:8080`) the probe falls back to `127.0.0.1`.

```sh
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
  CMD HOST="${FLTR_ADDR%:*}"; PORT="${FLTR_ADDR##*:}"; \
      case "$HOST" in ""|"0.0.0.0"|"::"|"[::]") HOST="127.0.0.1";; esac; \
      wget -q -O /dev/null "http://${HOST}:${PORT}/healthz"
```

- `:8080` becomes port `8080`
- `0.0.0.0:9090` becomes `9090`

If `FLTR_ADDR` is changed, the probe follows, so the container reports healthy as long as fltr answers on the address it was told to serve.

To check the health status while the container is running, use this:

```console
$ docker inspect --format '{{.State.Health.Status}}' fltr
healthy
```

## Image tags

The image is published under several tags, depending on what triggered the build:

- Every push to `main`: `fltr:main` and `fltr:sha-<short-sha>`, where `<short-sha>` is the commit's short hash.
- Every git version tag (e.g. `v1.0.0`): `fltr:latest`, plus the version split into `1`, `1.0`, and `1.0.0`.

!!! warning
    Only use the `latest` tag when you want to track the most recent release. It is recommended to pin a specific version to make sure nothing changes unattended.
