---
icon: lucide/code
---

# Development

## Project layout

- `main.go`: startup — reads the environment, validates upstreams, loads the rules, and wires the pieces together.
- `internal/config`: logging setup and Block Rule file loading.
- `internal/filter`: the request handler — body reading, Searchable Content assembly, matching, and verdict routing.
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

Commits must follow [Conventional Commits](https://www.conventionalcommits.org) (`feat:`, `fix:`, `docs:`, ...). The `commit-msg` hook enforces this, and the changelog and release titles are grouped by prefix.

## Docs site

The documentation site is built with [Zensical](https://zensical.org) from the public docs sources:

```sh
pip install zensical
zensical build --clean
```

The built `site/` is what gets published to GitHub Pages; do not edit it by hand.

## Releases

Releases are cut with [goreleaser](https://goreleaser.com): each `v`-prefixed git tag triggers a build of Linux binaries (amd64 and arm64), a checksum file, and a changelog assembled from the conventional commits. Docker images are published separately, see [Docker](docker.md#image-tags).
