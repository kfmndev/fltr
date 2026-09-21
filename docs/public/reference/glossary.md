---
icon: lucide/spell-check
---

# Glossary

The terms fltr uses, in the sense fltr uses them.

## Rules and matching

### Block Rule

A named set of terms. A rule *matches* a request when every one of its terms is found in the request's Searchable Content.

### Term

A single string a rule searches for within the Searchable Content.

### Match

A rule-level outcome: every term in the rule appears in the Searchable Content. It is binary — a request either matches a rule or it does not, and rules have no rank or precedence over each other.

### Searchable Content

The combined text fltr inspects: the request body, the `Title` and `Message` header values, and the `title` and `message` query parameter values.

## Verdicts and routing

### Allowed

The verdict for a request that matches no Block Rule; such a request is forwarded to the Allow Upstream.

### Blocked

The verdict for a request that matches at least one Block Rule; such a request is forwarded to the Block Upstream when configured, otherwise it is discarded.

### Allow Upstream

The HTTP service that receives allowed requests.

### Block Upstream

The optional HTTP service that receives blocked requests. When not configured, blocked requests are discarded instead.

### Discarded

The fate of a blocked request when no Block Upstream is configured: it is not forwarded and never reaches any upstream.

## Health and liveness

### Liveness Endpoint

The Reserved Path `/healthz`, answered by fltr itself with `200` whenever it can serve HTTP. Liveness is not readiness: the Liveness Endpoint never contacts an upstream and says nothing about whether requests can be forwarded.

### Reserved Path

A path fltr answers itself before content filtering, so it is never inspected, matched against Block Rules, or forwarded to any upstream.

## Forwarding

### Dropped Path

The fate of the incoming URL path: fltr does not forward it, so requests reach the upstream URL exactly as configured.
