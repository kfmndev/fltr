---
icon: lucide/network
---

# HTTP behavior

What fltr sends for every request, and the responses it produces itself.

## Responses

| Situation | Status | Body |
| --- | --- | --- |
| Allowed | the Allow Upstream's response | the upstream's response |
| Blocked, with `FLTR_BLOCK_UPSTREAM` set | the Block Upstream's response | the upstream's response |
| Blocked, without `FLTR_BLOCK_UPSTREAM` | `200` | `Request blocked, discarded` |
| Body larger than `FLTR_MAX_BODY_SIZE` | `413` | `request too large` |
| `GET` or `HEAD` to `/healthz` | `200` | `ok` |
| Any other method to `/healthz` | `405` | `method not allowed` |

The **Liveness Endpoint** (`/healthz`) is a **Reserved Path**. fltr answers it itself, so it is never inspected, matched against Block Rules, or forwarded. Its responses carry `Cache-Control: no-store`, and the `405` sets `Allow: GET, HEAD`.

## Forwarding

- The HTTP method, ordinary end-to-end headers, and body are preserved. Hop-by-hop headers are removed.
- fltr never sets or forwards `X-Forwarded-For`, `X-Forwarded-Host`, or `X-Forwarded-Proto`, and strips any such headers a proxy in front of it supplied.

!!! danger "Caution"
    When fltr fronts your upstream directly, the upstream sees the hop from fltr, not the original client.

- The incoming URL path is **dropped**. The request is sent to the configured upstream URL.
- Query parameters already present in the upstream URL are combined with the incoming query parameters. Duplicates are preserved.
- Nothing in the URL path is inspected for matching. fltr matches only the body, the `Title`/`Message` headers, and the `Title`/`title` and `Message`/`message` query parameters.
