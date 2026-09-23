---
icon: lucide/settings
---

# Environment variables

Apart from [**Block Rules**](block-rules.md), fltr is configured entirely through environment variables.

## Setting environment variables

To set variables temporarily, use `#!sh export <VAR_NAME>=<VALUE>`. Add the same command to your ~/.bashrc or the equivalent configuration file for your shell, to set them permanently.

Alternatively, prepend the environment variables to each fltr command like this `#!sh <VAR_NAME>=<VALUE> fltr`.

## Upstreams

| Variable | Required | Default |
| --- | --- | --- |
| `FLTR_ALLOW_UPSTREAM` | Yes | |
| `FLTR_BLOCK_UPSTREAM` | No | |

### `FLTR_ALLOW_UPSTREAM`

URL for requests that do not match any Block Rules. This is the only required variable.

!!! danger "Important"
    The **Allow Upstream** must be set and reachable at startup. fltr aborts the startup if it cannot connect.

### `FLTR_BLOCK_UPSTREAM`

URL for requests that match any Block Rules. Without one, blocked requests are discarded. See [How it works](how-it-works.md) for the full request flow.

## Block Rules

| Variable | Required | Default |
| --- | --- | --- |
| `FLTR_BLOCK_RULES_FILE` | No | `block_rules.json` |
| `FLTR_CASE_SENSITIVE` | No | `false` |

### `FLTR_BLOCK_RULES_FILE`

Path to the JSON Block Rules file. Defaults to `block_rules.json` in the current directory. See [Block Rules](block-rules.md) for the file format.

!!! danger "Important"
    The file must be readable and contain valid rules. Otherwise, fltr aborts the startup.

### `FLTR_CASE_SENSITIVE`

Enable case-sensitive Block Rule matching. By default both the [**Searchable Content**](block-rules.md#matching-semantics) and the terms are lowercased before matching; set this to `true` to match case-sensitively. Invalid values fall back to case-insensitive matching.

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

Log level: a valid [logrus Level](https://pkg.go.dev/github.com/sirupsen/logrus#Level), for example `debug`, `info`, or `warning`. Invalid values fall back to `info`.

### `LOG_FORMAT`

Log format: `text` (written to stderr) or `json` (written to stdout). Invalid values fall back to `text`.
