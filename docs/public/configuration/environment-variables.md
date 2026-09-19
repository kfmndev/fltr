---
icon: lucide/settings
---

# Environment variables

fltr is configured entirely through environment variables.

| Variable | Required | Default | Description |
| --- | --- | --- | --- |
| `FLTR_ALLOW_UPSTREAM` | Yes | | URL for requests that do not match any Block Rules |
| `FLTR_BLOCK_UPSTREAM` | No | | URL for requests that match any Block Rules |
| `FLTR_BLOCK_RULES_FILE` | No | `block_rules.json` | Path to the JSON Block Rules file |
| `FLTR_CASE_SENSITIVE` | No | `false` | Enable case-sensitive Block Rule matching. Invalid values fall back to case-insensitive matching |
| `FLTR_MAX_BODY_SIZE` | No | `10 MB` | Maximum request body size, using decimal units (e.g., `5 MB` = 5,000,000 bytes, `1 GB`, `512 KB`) |
| `FLTR_ADDR` | No | `:8080` | Address where the proxy listens |
| `LOG_LEVEL` | No | `info` | Log level: `trace`, `debug`, `info`, `warning` (`warn` also works), `error`, `fatal`, or `panic`. Invalid values fall back to `info` |
| `LOG_FORMAT` | No | `text` | Log format: `text` (stderr) or `json` (stdout). Invalid values fall back to `text` |

The Allow Upstream is the only required variable; fltr refuses to start without it.

!!! note "Setting environment variables"

    Set variables temporarily with `#!sh export FLTR_ALLOW_UPSTREAM=http://allow.example.com`, permanently by adding the same command to your `~/.bashrc` (or your shell's equivalent), or prepend them to every `fltr` command.

## Upstreams

At startup, fltr sends an HTTP `HEAD` request to each configured upstream before it starts listening. Any HTTP response counts as reachable, including error status codes. Network and TLS failures abort the startup.

!!! abstract "Important"

    The Allow Upstream must be reachable at startup. fltr aborts if it cannot connect.

Blocked requests go to the Block Upstream when `FLTR_BLOCK_UPSTREAM` is configured. Without one, Blocked requests are Discarded, and the proxy returns `200 Request blocked, discarded`. See [How it works](../how-it-works.md) for the full request flow.
