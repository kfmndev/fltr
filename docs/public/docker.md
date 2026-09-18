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

    The images are based on [`linuxserver/docker-baseimage-alpine`](https://github.com/linuxserver/docker-baseimage-alpine/). Set `PUID` and `PGID` to avoid file permission issues, see [Understanding PUID and PGID](https://docs.linuxserver.io/general/understanding-puid-and-pgid/) for details.

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

## Image tags

The image is published under several tags, depending on what triggered the build:

- Every push to `main`: `ghcr.io/kfmndev/fltr:main` and `ghcr.io/kfmndev/fltr:sha-<short-sha>`, where `<short-sha>` is the commit's short hash.
- Every **git** version tag (e.g. `v1.0.0`): `latest`, plus the version split into `1`, `1.0`, and `1.0.0`.

Use `latest` when you want to track the most recent release, or pin a specific version to make sure nothing changes.
