---
icon: lucide/settings
---

# Environment variables

Apart from [Block rules](block-rules.md), fltr is configured entirely through environment variables.

!!! note "Setting environment variables"
    Set variables temporarily with `#!sh export FLTR_ALLOW_UPSTREAM=http://allow.example.com`, permanently by adding the same command to your `~/.bashrc` (or your shell's equivalent), or prepend them to every `fltr` command.

## Upstreams

At startup, fltr sends an HTTP `HEAD` request to each configured upstream before it starts listening. Any HTTP response counts as reachable, including error status codes. Network and TLS failures abort the startup.

| Variable | Required | Default |
| --- | --- | --- |
| `FLTR_ALLOW_UPSTREAM` | Yes | |
| `FLTR_BLOCK_UPSTREAM` | No | |

### `FLTR_ALLOW_UPSTREAM`

URL for requests that do not match any Block Rules. This is the only required variable; fltr refuses to start without it.

!!! danger "Important"
    The Allow Upstream must be reachable at startup. fltr aborts if it cannot connect.

### `FLTR_BLOCK_UPSTREAM`

URL for requests that match any Block Rules. Without one, Blocked requests are Discarded, and the proxy returns `200 Request blocked, discarded`. See [How it works](../guide/how-it-works.md) for the full request flow.

## Block rules

| Variable | Required | Default |
| --- | --- | --- |
| `FLTR_BLOCK_RULES_FILE` | No | `block_rules.json` |
| `FLTR_CASE_SENSITIVE` | No | `false` |

### `FLTR_BLOCK_RULES_FILE`

Path to the JSON Block Rules file. Defaults to `block_rules.json` in the current directory. See [Block rules](block-rules.md) for the file format.

### `FLTR_CASE_SENSITIVE`

Enable case-sensitive Block Rule matching. By default both the Searchable Content and the Terms are lowercased before matching; set this to `true` to match case-sensitively. Invalid values fall back to case-insensitive matching.

## Request limits

| Variable | Required | Default |
| --- | --- | --- |
| `FLTR_MAX_BODY_SIZE` | No | `10 MB` |

### `FLTR_MAX_BODY_SIZE`

Maximum request body size, using decimal units (e.g., `5 MB` = 5,000,000 bytes, `1 GB`, `512 KB`) or a bare byte count. A request whose body exceeds this is rejected with `413 request too large`; see [HTTP behavior](http-behavior.md).

## Server

| Variable | Required | Default |
| --- | --- | --- |
| `FLTR_ADDR` | No | `:8080` |

### `FLTR_ADDR`

Address where the proxy listens.

## Logging

| Variable | Required | Default |
| --- | --- | --- |
| `LOG_LEVEL` | No | `info` |
| `LOG_FORMAT` | No | `text` |

### `LOG_LEVEL`

Log level: `trace`, `debug`, `info`, `warning` (`warn` also works), `error`, `fatal`, or `panic`. Invalid values fall back to `info`.

### `LOG_FORMAT`

Log format: `text` (stderr) or `json` (stdout). Invalid values fall back to `text`.
