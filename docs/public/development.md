---
icon: lucide/code
---

# Development

## How it works

fltr is an HTTP reverse proxy that filters requests based on their content:

1. **Startup**: an HTTP `HEAD` request is sent to each configured upstream before fltr starts listening. Any HTTP response counts as reachable; network and TLS failures abort startup. Block rules are loaded from `block_rules.json` by default, or `FLTR_BLOCK_RULES_FILE` if set.
2. **Matching**: fltr matches each incoming request's searchable content (body, `Title`/`Message` headers, `title`/`message` query parameters) against the block rules. [Block rules](block-rules.md) has the details.
3. **Routing**: allowed requests are forwarded to the allow upstream; blocked requests to the block upstream (if configured) or discarded with `200 Request blocked, discarded`.

Forwarding details:

- Method, headers, and body are preserved; the incoming path is discarded.
- Query parameters already present in the upstream URL are combined with the incoming query parameters.
- A request body larger than `FLTR_MAX_BODY_SIZE` (default 10 MB) is rejected with `413 Content Too Large`.

## Testing

Run the unit tests with:

```sh
go test ./...
```

Coverage can be generated with:

```sh
go test ./... -coverprofile=coverage.out
```
