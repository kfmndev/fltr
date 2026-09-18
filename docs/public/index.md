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

fltr inspects the body of every request, along with the `Title`/`Message` headers and the `title`/`message` query parameters, and matches that content against the block rules. A request that matches any rule is blocked; everything else is forwarded to the allow upstream. Method, headers, body, and query parameters are preserved. The incoming path is discarded.

## Requirements

- Go 1.25 or newer (when running from source)
- A JSON file containing the block rules
- An HTTP service to receive allowed requests

## Quick start

Start `fltr` with the minimal config: one block rule and an allow upstream.

```json title="block_rules.json"
{
    "rule": ["password", "secret"]
}
```

Set the allow upstream and run:

```sh
FLTR_ALLOW_UPSTREAM=http://allow.example.com \
fltr
```

The examples under [Block rules](block-rules.md) send requests that hit both the blocked and allowed paths.

## Installation

- **Pre-built binary**: download `fltr_<version>_linux_amd64` (or `linux_arm64`) from a [GitHub release](https://github.com/kfmndev/fltr/releases) and run it directly. The quickest way to try it out.
- **Docker**: images are published to `ghcr.io/kfmndev/fltr`, see [Docker](docker.md).
- **From source**: use `go build ./...`.

## Where to go next

- [Configuration](configuration.md): every environment variable and its defaults
- [Block rules](block-rules.md): file format and matching semantics
- [Docker](docker.md): images, tags, compose setup
- [Development](development.md): how it works, testing, and building
