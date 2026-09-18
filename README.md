<div align="center">
<img src="https://github.com/kfmndev/fltr/blob/main/assets/logo.svg?raw=true" alt="fltr logo" width="250px" style="max-width:25%"/>

**Lightweight HTTP reverse proxy for content-based request filtering and routing**

[![License](https://img.shields.io/github/license/kfmndev/fltr?style=flat-square)](https://github.com/kfmndev/fltr/blob/main/LICENSE.md) [![Go](https://img.shields.io/badge/go-00ADD8?style=flat-square&logo=go&logoColor=white)](https://go.dev) [![Tag](https://img.shields.io/github/v/tag/kfmndev/fltr?style=flat-square)](https://github.com/kfmndev/fltr/tags) [![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/kfmndev/fltr/build.yml?branch=main&style=flat-square)](https://github.com/kfmndev/fltr/actions?query=branch%3Amain)
</div>

## 🎯 TL;DR

- Request content matched against block rules (case-insensitive by default)
- Example block rule `rule1: password, secret`
- Allowed requests forwarded to "allow upstream"
- Blocked requests to "block upstream" (if configured) or discarded

> [!IMPORTANT]
> A request is blocked only when **ALL** terms from **ANY** rule match. Matching covers the request body, the `Title`/`Message` headers, and the `title`/`message` query parameters.

> [!NOTE]
> Method, headers, body, and query parameters are preserved; the path is **DISCARDED**.

## 🚀 Quick start

Requirements: [Go 1.25](https://go.dev) or newer (to run from source), a JSON file containing the block rules, and an HTTP service to receive allowed requests.

Create a `block_rules.json`:

```json
{
    "rule": ["password", "secret"]
}
```

Point fltr at your allow upstream:

```properties
export FLTR_ALLOW_UPSTREAM=http://allow.example.com
```

Then run it, whichever way suits you:

```sh
# From source
go run .

# From a pre-built binary (download from a GitHub release)
fltr

# With Docker
docker run -d \
  --name fltr \
  -p 8080:8080 \
  -v "$PWD/block_rules.json":/config/block_rules.json \
  -e FLTR_ALLOW_UPSTREAM=http://allow.example.com \
  ghcr.io/kfmndev/fltr
```

## 📚 Documentation

The full documentation lives at [kfmndev.github.io/fltr](https://kfmndev.github.io/fltr/):

- [Configuration](https://kfmndev.github.io/fltr/configuration/): every environment variable and its defaults
- [Block rules](https://kfmndev.github.io/fltr/block-rules/): file format and matching semantics
- [Docker](https://kfmndev.github.io/fltr/docker/): images, tags, compose setup, PUID/PGID
- [Development](https://kfmndev.github.io/fltr/development/): how it works, testing, and building
