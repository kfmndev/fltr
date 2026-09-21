---
icon: lucide/network
---

# HTTP behavior

What fltr sends on the wire for every request, and the responses it produces itself.

## Responses

| Situation | Status | Body |
| --- | --- | --- |
| Allowed | the Allow Upstream's response | the upstream's response |
| Blocked, with `FLTR_BLOCK_UPSTREAM` set | the Block Upstream's response | the upstream's response |
| Blocked, without `FLTR_BLOCK_UPSTREAM` | `200` | `Request blocked, discarded` |
| Body larger than `FLTR_MAX_BODY_SIZE` | `413` | `request too large` |
| `GET` or `HEAD` to `/healthz` | `200` | `ok` |
| Any other method to `/healthz` | `405` | `method not allowed` |

The [**Liveness Endpoint**](glossary.md#liveness-endpoint) (`/healthz`) is a [**Reserved Path**](glossary.md#reserved-path): it is answered by fltr itself and never inspected, matched against Block Rules, or forwarded. Its responses carry `Cache-Control: no-store`, and the `405` sets `Allow: GET, HEAD`.

## Forwarding

- The HTTP method, ordinary end-to-end headers, and body are preserved; hop-by-hop headers are removed by the reverse proxy.
- fltr never sets or forwards `X-Forwarded-For`, `X-Forwarded-Host`, or `X-Forwarded-Proto`, and strips any such headers a proxy in front of it supplied.
- The incoming URL path is **dropped**; the request is sent to the upstream URL exactly as configured.
- Query parameters already present in the upstream URL are combined with the incoming query parameters; duplicates are preserved.
- Nothing in the URL path is inspected for matching — only the body, the `Title`/`Message` headers, and the `Title`/`title` and `Message`/`message` query parameters.
