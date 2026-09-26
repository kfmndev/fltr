<div align="center">
<picture>
  <source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/kfmndev/fltr/refs/heads/main/assets/logo.svg">
  <img src="https://raw.githubusercontent.com/kfmndev/fltr/refs/heads/main/assets/logo-light.svg" alt="fltr logo" width="250px" style="max-width:25%"/>
</picture>

**Lightweight HTTP reverse proxy for content-based request filtering and routing**

[![License](https://img.shields.io/github/license/kfmndev/fltr?style=flat-square)](https://github.com/kfmndev/fltr/blob/main/LICENSE.md) [![Go](https://img.shields.io/badge/go-00ADD8?style=flat-square&logo=go&logoColor=white)](https://go.dev) [![Tag](https://img.shields.io/github/v/tag/kfmndev/fltr?style=flat-square)](https://github.com/kfmndev/fltr/tags) [![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/kfmndev/fltr/build.yml?branch=main&style=flat-square)](https://github.com/kfmndev/fltr/actions?query=branch%3Amain)
</div>

## 🎯 TL;DR

- Request content matched against **[Block Rules](https://kfmndev.github.io/fltr/reference/block-rules/)** (case-insensitive by default)
- Example Block Rule `credentials: ["password", "secret"]`
- Allowed requests forwarded to the **[Allow Upstream](https://kfmndev.github.io/fltr/reference/environment-variables/#upstreams)**
- Blocked requests to the **[Block Upstream](https://kfmndev.github.io/fltr/reference/environment-variables/#upstreams)** (if configured), or discarded

> [!IMPORTANT]
> A request is blocked only when **ALL** terms from **ANY** Block Rule match. Matching covers the [**Searchable Content**](https://kfmndev.github.io/fltr/reference/block-rules/#matching-semantics).

> [!NOTE]
> Method, ordinary end-to-end headers, and body are preserved (hop-by-hop headers are stripped); incoming query parameters are combined with those already in the upstream URL. The incoming URL path is **dropped**.

## 🚀 Quick start

Requirements: a JSON file containing the Block Rules and an HTTP service to receive allowed requests.

Create a `block_rules.json` with a single Block Rule:

```json
{
    "credentials": ["password", "secret"]
}
```

Set the Allow Upstream and run fltr:

```sh
FLTR_ALLOW_UPSTREAM=http://allow.example.com \
fltr
```

For details on building from source or running fltr in a Docker container, see [Getting started](https://kfmndev.github.io/fltr/guide/getting-started/).

## 📚 Documentation

The full documentation lives at [kfmndev.github.io/fltr](https://kfmndev.github.io/fltr/):

- **Guide**: [Getting started](https://kfmndev.github.io/fltr/guide/getting-started/), [Docker](https://kfmndev.github.io/fltr/guide/docker/), [Troubleshooting](https://kfmndev.github.io/fltr/guide/troubleshooting/), [FAQ](https://kfmndev.github.io/fltr/guide/faq/)
- **Reference**: [How it works](https://kfmndev.github.io/fltr/reference/how-it-works/), [Environment variables](https://kfmndev.github.io/fltr/reference/environment-variables/), [Block Rules](https://kfmndev.github.io/fltr/reference/block-rules/), [HTTP behavior](https://kfmndev.github.io/fltr/reference/http-behavior/), [Glossary](https://kfmndev.github.io/fltr/reference/glossary/)
- **Internals**: [Development](https://kfmndev.github.io/fltr/internals/development/)

## 📖 Project guidelines

- [Contributing](CONTRIBUTING.md): setup, commit conventions, and PR flow
- [Security](SECURITY.md): how to report a vulnerability privately
- [Code of Conduct](CODE_OF_CONDUCT.md)
