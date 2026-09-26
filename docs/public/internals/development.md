---
icon: lucide/code
---

# Development

## Project layout

- `main.go`: Startup; reads the environment, validates upstreams, loads the rules, and coordinates application components.
- `internal/config`: logging setup and [**Block Rule**](../reference/block-rules.md) file loading.
- `internal/filter`: the request handler for body reading, [**Searchable Content**](../reference/block-rules.md#matching-semantics) assembly, matching, and verdict routing.
- `internal/health`: the reserved [**Liveness Endpoint**](../reference/glossary.md#liveness-endpoint) handler that wraps the content filter.
- `internal/proxy`: upstream proxies and the startup reachability check.
- `docker/`: the Dockerfile and the s6-overlay service definitions for the published images.

## Requirements

- [Go 1.25](https://go.dev) or newer.
- [prek](https://github.com/j178/prek) (hooks for formatting, static analysis, compile, and tests).

## Getting started

```sh
go build ./...
prek install --hook-type pre-commit --hook-type commit-msg
```

## Testing

```sh
go test ./...
```

Coverage can be generated with:

```sh
go test ./... -coverprofile=coverage.out
```

## Commits

Commits must follow [Conventional Commits](https://www.conventionalcommits.org) (`feat:`, `fix:`, `docs:`, ...). The `commit-msg` hook enforces this, and the changelog is grouped by prefix.

For example:
`feat: add Docker build support (#1)`

## Docs site

The documentation site is built with [Zensical](https://zensical.org) from the public docs sources:

```sh
pip install "zensical==0.0.62"
zensical build --clean
```

The built `site/` is what gets published to GitHub Pages; do not edit it by hand.

## Terminology

The docs follow a two-tier convention for domain terms.

Multi-word terms keep their capitals in prose: `Block Rule`, `Allow Upstream`, `Block Upstream`, `Searchable Content`, `Liveness Endpoint`, `Reserved Path`, `Dropped Path`. Single-word state terms are lowercase except at the start of a sentence: `term`, `match`/`matching`, `allowed`, `blocked`, `discarded`.

The first mention of a term on a page is bold and links to its most useful destination: the term's dedicated page or section when it has one (Block Rules, upstreams, Searchable Content), otherwise its anchor in the [Glossary](../reference/glossary.md). On the page that documents a term, the mention is bold but unlinked. Later mentions are plain.

Short labels, as found in table cells, diagram nodes and step headings, keep their capitals. `CONTEXT.md` uses full capitals throughout. It is the canonical glossary, not prose.

## Releases

Releases are created using [goreleaser](https://goreleaser.com): each `v`-prefixed Git tag triggers a build of Linux binaries (amd64 and arm64), a checksum file, and a changelog assembled from the conventional commits. Docker images are published separately. See [Docker](../guide/docker.md#image-tags).
