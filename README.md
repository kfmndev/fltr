# fltr

`fltr` is a small HTTP reverse proxy that filters and routes requests based on their content. It sends matching requests to a block upstream and everything else to an allow upstream.

## How it works

The proxy checks the request body and the `Title` and `Message` headers. Matching is case-insensitive by default, but can be made case-sensitive with the `FLTR_CASE_SENSITIVE` environment variable. A request is blocked when every term in any one configured rule appears in that combined text.

Allowed requests go to `FLTR_ALLOW_UPSTREAM`. Blocked requests go to `FLTR_BLOCK_UPSTREAM` when it is configured. Without a block upstream, blocked requests are discarded and the proxy returns `200 Request blocked, discarded`. The proxy keeps the original method, headers, body, and query parameters when forwarding. It sends requests to the configured upstream URL without preserving the incoming path. Query parameters already present on the upstream URL are preserved and combined with the incoming query parameters. Request bodies are limited by `FLTR_MAX_BODY_SIZE` (default 10 MB); larger requests are rejected with a 413 status code.

## Requirements

- Go 1.23 or newer
- An HTTP service to receive allowed requests
- A JSON file containing the terms to check

## Run it

Start `fltr` with the allow and block services from your environment:

```sh
FLTR_ALLOW_UPSTREAM=http://localhost:9000 \
FLTR_BLOCK_UPSTREAM=http://localhost:9001 \
FLTR_BLOCKED_FILE=block_rules.json \
FLTR_ADDR=:8080 \
go run .
```

The proxy checks each configured upstream with an HTTP `HEAD` request before it starts listening. Any HTTP response counts as reachable, including error status codes. Network and TLS failures stop startup. If `FLTR_BLOCKED_FILE` is unset, it loads `block_rules.json` from the current directory.

### Configuration

| Variable | Required | Default | Description |
| --- | --- | --- | --- |
| `FLTR_ALLOW_UPSTREAM` | Yes | | URL for requests that do not match the blocklist |
| `FLTR_BLOCKED_FILE` | No | `block_rules.json` | Path to the JSON block rules |
| `FLTR_BLOCK_UPSTREAM` | No | | URL for requests that match the blocklist |
| `FLTR_CASE_SENSITIVE` | No | `false` | Enable case-sensitive rule matching |
| `FLTR_MAX_BODY_SIZE` | No | `10 MB` | Maximum request body size (e.g., `5 MB`, `1 GB`, `512 KB`) |
| `FLTR_ADDR` | No | `:8080` | Address where the proxy listens |
| `LOG_LEVEL` | No | `info` | Log level, such as `debug`, `info`, or `warn` |

## Blocklist format

The block rules file is a JSON object. Its keys name rules, and each value is a non-empty array of non-empty terms. A request matches when all terms in at least one rule are present. Terms are trimmed when the file is loaded; if case-insensitive matching is enabled (default), they are also lowercased. The service refuses to start if a rule has no terms or contains a blank term, including a term made blank by surrounding whitespace.

```json
{
    "credentials": ["password", "secret"],
    "identity": ["ssn"]
}
```

With the example above, a request is blocked when it contains both `password` and `secret`, or when it contains `ssn`, somewhere in its body or `Title` and `Message` headers.

## Try a request

With the example configuration running, send a request that matches the `credentials` rule:

```sh
curl -X POST http://localhost:8080 \
    -H 'Content-Type: text/plain' \
    -H 'Title: account details' \
    -d 'password secret'
```

Send a request that matches no rule to route it to the allow upstream:

```sh
curl -X POST http://localhost:8080 \
    -H 'Content-Type: text/plain' \
    -d 'hello service'
```

## Test

Run the unit tests with:

```sh
go test ./...
```

## Project layout

- `main.go` loads configuration, builds the reverse proxies, and starts the HTTP server.
- `filter.go` reads request bodies and chooses the allow or block upstream.
- `block_rules.json` contains the default block rules.
