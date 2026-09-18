---
icon: lucide/home
---

# Home

<div align="center">
<img src="assets/logo.svg" alt="fltr logo" width="250px" style="max-width:50%"/>
<br>
<b>Lightweight HTTP reverse proxy for content-based request filtering and routing</b>
<br style="margin-bottom: .5rem">
<a href="https://github.com/kfmndev/fltr/blob/main/LICENSE.md"><img src="https://img.shields.io/github/license/kfmndev/fltr?style=flat-square" alt="License"></a>
<a href="https://go.dev"><img src="https://img.shields.io/badge/go-00ADD8?style=flat-square&logo=go&logoColor=white" alt="Go"></a>
<a href="https://github.com/kfmndev/fltr/tags"><img src="https://img.shields.io/github/v/tag/kfmndev/fltr?style=flat-square" alt="Tag"></a>
<a href="https://github.com/kfmndev/fltr/actions?query=branch%3Amain"><img src="https://img.shields.io/github/actions/workflow/status/kfmndev/fltr/build.yml?branch=main&style=flat-square" alt="GitHub Actions Workflow Status"></a>
</div>

fltr inspects the body of every request, along with the `Title`/`Message` headers and the `Title`/`title` and `Message`/`message` query parameters, and matches that content against the Block Rules. A request that matches any Block Rule is Blocked; everything else is forwarded to the Allow Upstream. Method, ordinary end-to-end headers, and body are preserved (hop-by-hop headers are stripped); incoming query parameters are combined with those already in the upstream URL. The incoming URL path is dropped.

## Requirements

- Go 1.25 or newer (when running from source)
- A JSON file containing the Block Rules
- An HTTP service to receive Allowed requests

## Quick start

Start `fltr` with the minimal config: one Block Rule and an Allow Upstream.

```json title="block_rules.json"
{
    "rule": ["password", "secret"]
}
```

Set the Allow Upstream and run:

```sh
FLTR_ALLOW_UPSTREAM=http://allow.example.com \
fltr
```

The examples under [Block rules](block-rules.md) show a request that is Blocked and one that is Allowed.

## Installation

- **Pre-built binary**: download `fltr_<version>_linux_amd64` (or `linux_arm64`) from a [GitHub release](https://github.com/kfmndev/fltr/releases) and run it directly. The quickest way to try it out.
- **Docker**: images are published to `ghcr.io/kfmndev/fltr`, see [Docker](docker.md).
- **From source**: use `go build ./...`.

## Where to go next

- [Configuration](configuration.md): every environment variable and its defaults
- [Block rules](block-rules.md): file format and matching semantics
- [Docker](docker.md): images, tags, compose setup
- [How it works](how-it-works.md): startup, matching, and routing
- [Troubleshooting](troubleshooting.md): common failures and fixes
- [Development](development.md): project layout, testing, building, and releases
