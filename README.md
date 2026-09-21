<div align="center">
<img src="https://raw.githubusercontent.com/kfmndev/fltr/refs/heads/main/assets/logo.svg" alt="fltr logo" width="250px" style="max-width:25%"/>

**Lightweight HTTP reverse proxy for content-based request filtering and routing**

[![License](https://img.shields.io/github/license/kfmndev/fltr?style=flat-square)](https://github.com/kfmndev/fltr/blob/main/LICENSE.md) [![Go](https://img.shields.io/badge/go-00ADD8?style=flat-square&logo=go&logoColor=white)](https://go.dev) [![Tag](https://img.shields.io/github/v/tag/kfmndev/fltr?style=flat-square)](https://github.com/kfmndev/fltr/tags) [![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/kfmndev/fltr/build.yml?branch=main&style=flat-square)](https://github.com/kfmndev/fltr/actions?query=branch%3Amain)
</div>

## 🎯 TL;DR

- Request content matched against **[Block Rules](https://kfmndev.github.io/fltr/reference/block-rules/)** (case-insensitive by default)
- Example Block Rule `rule: password, secret`
- Allowed requests forwarded to the **[Allow Upstream](https://kfmndev.github.io/fltr/reference/environment-variables/#upstreams)**
- Blocked requests to the **[Block Upstream](https://kfmndev.github.io/fltr/reference/environment-variables/#upstreams)** (if configured), or discarded

> [!IMPORTANT]
> A request is blocked only when **ALL** terms from **ANY** Block Rule match. Matching covers the [**Searchable Content**](https://kfmndev.github.io/fltr/reference/block-rules/#matching-semantics): the request body, the `Title`/`Message` headers, and the `Title`/`title` and `Message`/`message` query parameters.

> [!NOTE]
> Method, ordinary end-to-end headers, and body are preserved (hop-by-hop headers are stripped); incoming query parameters are combined with those already in the upstream URL. The incoming URL path is **dropped**.

## 🚀 Quick start

Requirements: [Go 1.25](https://go.dev) or newer (to run from source), a JSON file containing the Block Rules, and an HTTP service to receive allowed requests.

Create a `block_rules.json`:

```json
{
    "rule": ["password", "secret"]
}
```

Point fltr at your Allow Upstream:

```sh
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

- [Quickstart](https://kfmndev.github.io/fltr/getting-started/): one Block Rule, one upstream, first request
- [Installation](https://kfmndev.github.io/fltr/getting-started/installation/): pre-built binary, Docker, or source
- [How it works](https://kfmndev.github.io/fltr/guide/how-it-works/): startup, matching, and routing
- [Docker](https://kfmndev.github.io/fltr/guide/docker/): images, tags, compose setup, PUID/PGID
- [Troubleshooting](https://kfmndev.github.io/fltr/guide/troubleshooting/): common failures and fixes
- [FAQ](https://kfmndev.github.io/fltr/guide/faq/): recurring questions about matching and forwarding
- [Environment variables](https://kfmndev.github.io/fltr/reference/environment-variables/): every environment variable and its defaults
- [Block Rules](https://kfmndev.github.io/fltr/reference/block-rules/): file format and matching semantics
- [HTTP behavior](https://kfmndev.github.io/fltr/reference/http-behavior/): status codes, headers, and path handling
- [Glossary](https://kfmndev.github.io/fltr/reference/glossary/): the terms fltr uses
- [Development](https://kfmndev.github.io/fltr/development/): project layout, testing, building, and releases

## Contributing

- [Contributing](CONTRIBUTING.md): setup, commit conventions, and PR flow
- [Security](SECURITY.md): how to report a vulnerability privately
- [Code of Conduct](CODE_OF_CONDUCT.md)
