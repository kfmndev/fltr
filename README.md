<div align="center">
<img src="https://github.com/kfmndev/fltr/blob/main/assets/logo.svg?raw=true" alt="fltr logo" width="250px" style="max-width:25%"/>

**Lightweight HTTP reverse proxy for content-based request filtering and routing**

[![License](https://img.shields.io/github/license/kfmndev/fltr?style=flat-square)](https://github.com/kfmndev/fltr/blob/main/LICENSE.md) [![Go](https://img.shields.io/badge/go-00ADD8?style=flat-square&logo=go&logoColor=white)](https://go.dev) [![Tag](https://img.shields.io/github/v/tag/kfmndev/fltr?style=flat-square)](https://github.com/kfmndev/fltr/tags) [![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/kfmndev/fltr/build.yml?branch=main&style=flat-square)](https://github.com/kfmndev/fltr/actions?query=branch%3Amain)
</div>

## 🎯 TL;DR

- Request content matched against block rules (case-insensitive by default)
- Example block rule `rule1: password, secret`
- Blocked, when **ALL** terms from **ANY** rule match
- Allowed requests forwarded to "allow upstream"
- Blocked requests to "block upstream" (if configured) or discarded
- Method, headers, body, and query parameters are preserved; path is **DISCARDED**

## ✅ Requirements

- Go 1.25 or newer
- A JSON file containing the block rules
- An HTTP service to receive allowed requests

## 🚀 Run it

Start `fltr` with the minimal config: one block rule and an allow upstream.

### block_rules.json

```json
{
    "rule": ["password", "secret"]
}
```

### FLTR_ALLOW_UPSTREAM

Set the environment variable temporarily like this:

```sh
export FLTR_ALLOW_UPSTREAM=http://allow.example.com
```

To make it permanent, add this to the `~/.bashrc` (or your shell's equivalent).

Alternatively, the command can be prepended with the variable:

```sh
FLTR_ALLOW_UPSTREAM=http://allow.example.com \
go run .
```

## 🐳 Run it with Docker

Pre-built images are published to the GitHub Container Registry as `ghcr.io/kfmndev/fltr`. They are built for `linux/amd64` and `linux/arm64`.

```sh
docker run -d \
  --name fltr \
  -p 8080:8080 \
  -v "$PWD/block_rules.json":/config/block_rules.json \
  -e FLTR_ALLOW_UPSTREAM=http://allow.example.com \
  ghcr.io/kfmndev/fltr
```

The container follows LinuxServer.io conventions: the block rules are read from `/config/block_rules.json` by default (set `FLTR_BLOCKED_FILE` to override the path), and the service runs as the `abc` user, remappable with the `PUID` and `PGID` environment variables (default `911`). Keep the example above in mind: `FLTR_ALLOW_UPSTREAM` is required, and the proxy listens on port `8080`.

### Compose

A [`docker-compose.yml`](docker-compose.yml) is included for a declarative setup. Its defaults live in a `.env` file next to it:

```sh
cp .env.example .env
# set FLTR_ALLOW_UPSTREAM (required)
docker compose up -d
```

The compose setup adds a few Docker-specific variables; the `FLTR_*` runtime settings are still listed in the [Configuration](#-configuration) table below:

| Variable | Default | Description |
| --- | --- | --- |
| `FLTR_RULES_FILE` | `./block_rules.json` | Host file mounted as the container's block rules |
| `FLTR_HTTP_PORT` | `8080` | Host port mapped to the container's `8080` |
| `PUID` / `PGID` | `911` | User/group the service runs as |
| `TZ` | `Etc/UTC` | Timezone |

### Image tags

The image is published under several tags, depending on what triggered the build:

- Every push to `main`: `ghcr.io/kfmndev/fltr:main` and `ghcr.io/kfmndev/fltr:sha-<short-sha>`, where `<short-sha>` is the commit's short hash.
- Every **git** version tag (e.g. `v1.0.0`): `latest`, plus the version split into `1`, `1.0`, and `1.0.0`.

So `ghcr.io/kfmndev/fltr:latest` tracks the most recent release, and `ghcr.io/kfmndev/fltr:1.0.0` pins the image to a specific version.

## 🔧 Configuration

| Variable | Required | Default | Description |
| --- | --- | --- | --- |
| `FLTR_ALLOW_UPSTREAM` | Yes | | URL for requests that do not match the blocklist |
| `FLTR_BLOCK_UPSTREAM` | No | | URL for requests that match the blocklist |
| `FLTR_BLOCKED_FILE` | No | `block_rules.json` | Path to the JSON block rules |
| `FLTR_CASE_SENSITIVE` | No | `false` | Enable case-sensitive rule matching |
| `FLTR_MAX_BODY_SIZE` | No | `10 MB` | Maximum request body size (e.g., `5 MB`, `1 GB`, `512 KB`) |
| `FLTR_ADDR` | No | `:8080` | Address where the proxy listens |
| `LOG_LEVEL` | No | `info` | Log level, such as `debug`, `info`, or `warn` |

## 🚫 Blocklist format

The block rules file is a JSON object. Its keys name rules, and each value is a non-empty array of non-empty terms. A request matches when all terms in at least one rule are present. Terms are trimmed when the file is loaded. The service refuses to start if a rule has no terms or contains a blank term after trimming.

```json
{
    "credentials": ["password", "secret"],
    "identity": ["private-key"]
}
```

With the example above, a request is blocked when it contains both `password` and `secret`, or when it contains `private-key`, somewhere in its body or `Title` and `Message` headers.

## 🧪 Try a request

With the example configuration running, send a request that matches the `credentials` rule:

```sh
curl -X POST http://localhost:8080 \
    -H 'Content-Type: text/plain' \
    -H 'Title: account details' \
    -d 'password secret'
```

Send a request that matches no rule to route it to the allow upstream:

```sh
curl -X POST http://localhost:8080 \
    -H 'Content-Type: text/plain' \
    -d 'hello service'
```

## ⚙️ How it works

During startup, an HTTP `HEAD` request is sent to each configured upstream before it starts listening. Any HTTP response is considered reachable, including error status codes. Network and TLS failures abort the startup process.

By default, the block rules are loaded from the `block_rules.json` file in the current directory. The environment variable `FLTR_BLOCKED_FILE` can specify an alternative path.

When a request reaches the proxy, its body and the `Title` and `Message` headers are matched against the block rules. By default, matching is case-insensitive, which is achieved by converting both the combined request text and the block terms to lowercase before matching. This can be configured to be case-sensitive using the `FLTR_CASE_SENSITIVE` environment variable. A request is blocked when all the terms in any configured rule appear in the combined text.

Allowed requests go to `FLTR_ALLOW_UPSTREAM`. Blocked requests go to `FLTR_BLOCK_UPSTREAM` when it is configured. Without a block upstream, blocked requests are discarded, and the proxy returns `200 Request blocked, discarded`.

It forwards requests to the configured upstream URL without preserving the incoming path. The original method, headers, and body are kept. Query parameters already present in the upstream URL are combined with the incoming query parameters.

Request bodies are limited in size by `FLTR_MAX_BODY_SIZE` (default 10 MB); larger requests are rejected with the `413 Content Too Large` status code.

## ✅ Testing

Run the unit tests with:

```sh
go test ./...
```
