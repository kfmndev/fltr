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

A rule-level outcome: every term in the rule appears in the Searchable Content. Matching is binary. A request either matches a rule or it doesn't, and rules have no rank or precedence over each other.

### Searchable Content

The combined text fltr inspects: the request body, the `Title`/`title` and `Message`/`message` header values, and the `Title`/`title` and `Message`/`message` query parameter values.

## Verdicts and routing

### Allowed

The verdict for a request that matches no Block Rule. Such a request is forwarded to the Allow Upstream.

### Blocked

The verdict for a request that matches at least one Block Rule. Such a request is forwarded to the Block Upstream when configured. Otherwise it is discarded.

### Allow Upstream

The HTTP service that receives allowed requests.

### Block Upstream

The optional HTTP service that receives blocked requests. When not configured, blocked requests are discarded instead.

### Discarded

A blocked request with no Block Upstream configured is not forwarded and never reaches any upstream.

## Health and liveness

### Liveness Endpoint

The Reserved Path `/healthz`, answered by fltr itself with `200` whenever it can serve HTTP. Liveness is not readiness. The Liveness Endpoint never contacts an upstream and says nothing about whether requests can be forwarded.

### Reserved Path

A path fltr answers itself before content filtering, so it is never inspected, matched against Block Rules, or forwarded to any upstream.

## Forwarding

### Dropped Path

The incoming URL path is not forwarded, so requests reach the upstream URL exactly as configured.

## Configuration

### Configuration Argument

A setting supplied as a command-line argument. It takes precedence over the same setting's environment variable, which in turn takes precedence over the default.
