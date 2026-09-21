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

fltr inspects the [**Searchable Content**](reference/block-rules.md#matching-semantics) of every request, matching it against the [**Block Rules**](reference/block-rules.md). A request that matches any Block Rule is blocked. Everything else is forwarded to the [**Allow Upstream**](reference/environment-variables.md#upstreams). fltr preserves the method, ordinary end-to-end headers, and body, while stripping hop-by-hop headers. It combines incoming query parameters with those already in the upstream URL, and drops the incoming URL path.

## Where to go next

**Getting started**

- [Quickstart](getting-started/index.md): the shortest path to a running proxy
- [Installation](getting-started/installation.md): pre-built binary, Docker, or source

**Guide**

- [How it works](guide/how-it-works.md): startup, matching, and routing
- [Docker](guide/docker.md): images, tags, compose setup
- [Troubleshooting](guide/troubleshooting.md): common failures and fixes
- [FAQ](guide/faq.md): recurring questions about matching and forwarding

**Reference**

- [Environment variables](reference/environment-variables.md): every environment variable and its defaults
- [Block Rules](reference/block-rules.md): file format and matching semantics
- [HTTP behavior](reference/http-behavior.md): status codes, headers, and path handling
- [Glossary](reference/glossary.md): the terms fltr uses

**Project**

- [Development](development.md): project layout, testing, building, and releases
