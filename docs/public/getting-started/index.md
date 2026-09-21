---
icon: lucide/rocket
---

# Quickstart

This gets fltr running locally with one [**Block Rule**](../reference/glossary.md#block-rule) and one [**Allow Upstream**](../reference/glossary.md#allow-upstream). For other ways to run it, see [Installation](installation.md).

## Requirements

- Go 1.25 or newer (when running from source)
- A JSON file containing the Block Rules
- An HTTP service to receive allowed requests

## Start fltr

Create a `block_rules.json` with a single Block Rule:

```json title="block_rules.json"
{
    "rule": ["password", "secret"]
}
```

Set the Allow Upstream and run fltr:

```sh
FLTR_ALLOW_UPSTREAM=http://allow.example.com \
fltr
```

By default fltr listens on `:8080` and reads `block_rules.json` from the current directory. Both are configurable; see [Environment variables](../reference/environment-variables.md).

## Send a request

A request whose [**Searchable Content**](../reference/glossary.md#searchable-content) contains every term of a Block Rule is blocked:

```sh
curl -X POST http://localhost:8080 \
    -H 'Content-Type: text/plain' \
    -H 'Title: account details' \
    -d 'password secret'
```

With no [**Block Upstream**](../reference/glossary.md#block-upstream) configured, fltr answers `200 Request blocked, discarded`. Requests that match no Block Rule are forwarded to the Allow Upstream. [Block Rules](../reference/block-rules.md) shows both outcomes in full, and [HTTP behavior](../reference/http-behavior.md) documents every response fltr produces itself.
