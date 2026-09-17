# Configuration

fltr is configured entirely through environment variables.

| Variable | Required | Default | Description |
| --- | --- | --- | --- |
| `FLTR_ALLOW_UPSTREAM` | Yes | | URL for requests that do not match any block rules |
| `FLTR_BLOCK_UPSTREAM` | No | | URL for requests that match any block rules |
| `FLTR_BLOCK_RULES_FILE` | No | `block_rules.json` | Path to the JSON block rules file |
| `FLTR_CASE_SENSITIVE` | No | `false` | Enable case-sensitive rule matching |
| `FLTR_MAX_BODY_SIZE` | No | `10 MB` | Maximum request body size (e.g., `5 MB`, `1 GB`, `512 KB`) |
| `FLTR_ADDR` | No | `:8080` | Address where the proxy listens |
| `LOG_LEVEL` | No | `info` | Log level, such as `debug`, `info`, or `warn` |
| `LOG_FORMAT` | No | `text` | Log format: `text` or `json` (`json` only writes to stdout) |

The allow upstream is the only required variable; fltr refuses to start without it.

!!! note "Persistent environment variables"

    Set variables temporarily with `export FLTR_ALLOW_UPSTREAM=http://allow.example.com`, permanently via `~/.bashrc` (or your shell's equivalent), or by prepending them to the start command.

## Upstreams

At startup, fltr sends an HTTP `HEAD` request to each configured upstream before it starts listening. Any HTTP response counts as reachable, including error status codes. Network and TLS failures abort the startup.

!!! abstract "Important"

    The allow upstream must be reachable at startup. fltr aborts if it cannot connect.

Blocked requests go to the block upstream when `FLTR_BLOCK_UPSTREAM` is configured. Without one, blocked requests are discarded, and the proxy returns `200 Request blocked, discarded`. See [Development](development.md) for the full request flow.
