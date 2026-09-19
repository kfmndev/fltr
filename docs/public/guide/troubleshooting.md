---
icon: lucide/wrench
---

# Troubleshooting

Symptoms first, causes and fixes below each one. See [How it works](how-it-works.md) if you need the underlying model first.

## fltr refuses to start

### `FLTR_ALLOW_UPSTREAM is required`

Set the Allow Upstream: `export FLTR_ALLOW_UPSTREAM=https://allow.example.com`. It is the only required variable.

### `FLTR_ALLOW_UPSTREAM is unreachable` (or the Block Upstream equivalent)

fltr sent a `HEAD` request and got no HTTP response at all. Check the URL, DNS, and connectivity — self-signed TLS certificates fail here too. Any HTTP response, even an error status, proves the upstream is reachable.

### `could not load block rules`

The Block Rules file is missing, unreadable, or invalid JSON. Check the path (default `block_rules.json` in the current directory, or `FLTR_BLOCK_RULES_FILE`), and validate the JSON.

### `rule "..." has no terms` / `rule "..." has an empty term at index N`

A Block Rule key maps to an empty array, or a Term is blank after trimming. Every Block Rule needs at least one non-empty Term.

### `invalid FLTR_MAX_BODY_SIZE`

The value must be a size such as 10 MB, 512 KB, or 1 GB (decimal units: 1 KB = 1000 bytes), or a bare byte count.

## Requests behave unexpectedly

### I get `200 Request blocked, discarded` for a request I expected to be allowed

A Block Rule matched: every Term of that rule was found in the Searchable Content. Remember matching also covers the `Title`/`Message` headers and the `Title`/`title` and `Message`/`message` query parameters; matching is case-insensitive unless you set `FLTR_CASE_SENSITIVE=true`; and `strings.Contains` means Terms match inside larger words — a Block Rule with the Term `password` also matches `PasswordManager`. `LOG_LEVEL=debug` logs which Block Rule and Terms matched.

### My request reached the upstream but at the wrong path

The incoming URL path is dropped by design; requests arrive at the upstream URL exactly as configured. See [Forwarding behavior](how-it-works.md#forwarding-behavior).

### The upstream saw different query parameters than I sent

Incoming query parameters are appended to whatever query string is already part of the upstream URL; duplicates are preserved.

### `413 request too large`

The request body exceeded `FLTR_MAX_BODY_SIZE` (default 10 MB). Raise it if the contents are legitimate.

### The upstream doesn't know the client's IP or original URL

fltr does not set or forward `X-Forwarded-For`, `X-Forwarded-Host`, or `X-Forwarded-Proto` to the upstream, so adding them at a proxy in front of fltr will not make them reach the upstream.

## Docker

### Permission errors reading the rules file

The images are built from [`linuxserver/docker-baseimage-alpine`](https://github.com/linuxserver/docker-baseimage-alpine/) (`ghcr.io/linuxserver/baseimage-alpine`) and run as `abc:abc`. Set `PUID`/`PGID` to your user so the bind-mounted `block_rules.json` is readable; see [Understanding PUID and PGID](https://docs.linuxserver.io/general/understanding-puid-and-pgid/), and [Docker](docker.md) for the setup.
